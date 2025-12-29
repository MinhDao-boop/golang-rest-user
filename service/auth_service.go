package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"

	"golang-rest-user/dto"
	"golang-rest-user/models"
	"golang-rest-user/repository"
	"golang-rest-user/security"
)

type AuthService interface {
	Register(req dto.CreateUserRequest) (*models.User, error)
	Login(tenantCode string, req dto.LoginRequest) (map[string]string, error)
	Refresh(refreshToken string) (map[string]string, error)
	Logout(refreshToken string) error
}

type authService struct {
	userRepo         repository.UserRepo
	refreshTokenRepo repository.RefreshTokenRedis
	jwtManager       *security.Manager
}

func NewAuthService(userRepo repository.UserRepo, refreshTokenRepo repository.RefreshTokenRedis,
	jwtManager *security.Manager) AuthService {
	return &authService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtManager:       jwtManager,
	}
}

func (s *authService) Register(req dto.CreateUserRequest) (*models.User, error) {

	if _, err := s.userRepo.GetByUsername(req.Username); err == nil {
		return nil, errors.New("username already exists")
	}

	encryptedPass, err := security.Encrypt(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Uuid:     uuid.NewString(),
		Username: req.Username,
		Password: encryptedPass,
		FullName: req.FullName,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(tenantCode string, req dto.LoginRequest) (map[string]string, error) {

	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	decryptedPass, _ := security.Decrypt(user.Password)
	if decryptedPass != req.Password {
		return nil, errors.New("invalid credentials")
	}

	aToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Username, tenantCode)
	if err != nil {
		return nil, err
	}

	rToken, err := s.jwtManager.GenerateRefreshToken(user.ID, tenantCode)
	if err != nil {
		return nil, err
	}

	hash := hashToken(rToken.Token)
	ttl := time.Duration(rToken.ExpiresIn) * time.Second

	err = s.refreshTokenRepo.Create(hash, user.ID, tenantCode, ttl)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":       aToken.Token,
		"access_expires_in":  strconv.FormatInt(aToken.ExpiresIn, 10),
		"refresh_token":      rToken.Token,
		"refresh_expires_in": strconv.FormatInt(rToken.ExpiresIn, 10),
	}, nil
}

func hashToken(rToken string) string {
	h := sha256.Sum256([]byte(rToken))
	return hex.EncodeToString(h[:])
}

func (s *authService) Refresh(rToken string) (map[string]string, error) {
	claims, err := s.jwtManager.ParseToken(rToken)

	if err != nil || claims.Type != "refresh" {
		return nil, errors.New("invalid refresh token")
	}

	if err := s.refreshTokenRepo.FindValidByHash(hashToken(rToken)); err != nil {
		return nil, errors.New("refresh token revoked")
	}

	//revoke old refresh token
	if err = s.refreshTokenRepo.Revoke(hashToken(rToken)); err != nil {
		return nil, err
	}

	user, _ := s.userRepo.GetByID(claims.UserID)

	newAToken, _ := s.jwtManager.GenerateAccessToken(claims.UserID, user.Username, claims.TenantCode)
	newRToken, _ := s.jwtManager.GenerateRefreshToken(claims.UserID, claims.TenantCode)

	hash := hashToken(newRToken.Token)
	ttl := time.Duration(newRToken.ExpiresIn) * time.Second

	if err = s.refreshTokenRepo.Create(hash, claims.UserID, claims.TenantCode, ttl); err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":       newAToken.Token,
		"access_expires_in":  strconv.FormatInt(newAToken.ExpiresIn, 10),
		"refresh_token":      newRToken.Token,
		"refresh_expires_in": strconv.FormatInt(newRToken.ExpiresIn, 10),
	}, nil
}

func (s *authService) Logout(rToken string) error {
	claims, err := s.jwtManager.ParseToken(rToken)
	if err != nil || claims.Type != "refresh" {
		return errors.New("invalid refresh token")
	}

	return s.refreshTokenRepo.RevokeAllByUser(claims.UserID)
}
