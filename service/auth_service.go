package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"golang-rest-user/dto"
	"golang-rest-user/models"
	"golang-rest-user/repository"
	"golang-rest-user/security"
)

type AuthService interface {
	Register(tenantCode string, req dto.CreateUserRequest) (*models.User, error)
	Login(tenantCode string, req dto.LoginRequest) (map[string]string, error)
	Refresh(refreshToken string) (map[string]string, error)
	Logout(refreshToken string) error
}

type authService struct {
	userRepoFactory         func(string) (repository.UserRepo, error)
	refreshTokenRepoFactory func(string) (repository.RefreshTokenRepo, error)
	jwtManager              *security.Manager
}

func NewAuthService(
	userRepoFactory func(string) (repository.UserRepo, error),
	refreshTokenRepoFactory func(string) (repository.RefreshTokenRepo, error),
	jwtManager *security.Manager,
) AuthService {
	return &authService{
		userRepoFactory:         userRepoFactory,
		refreshTokenRepoFactory: refreshTokenRepoFactory,
		jwtManager:              jwtManager,
	}
}

func (s *authService) Register(tenantCode string, req dto.CreateUserRequest) (*models.User, error) {

	repo, err := s.userRepoFactory(tenantCode)
	if err != nil {
		return nil, err
	}

	if _, err := repo.GetByUsername(req.Username); err == nil {
		return nil, errors.New("username already exists")
	}

	encryptedPass, err := security.Encrypt(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: req.Username,
		Password: encryptedPass,
		FullName: req.FullName,
	}

	if err := repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(tenantCode string, req dto.LoginRequest) (map[string]string, error) {

	repo, err := s.userRepoFactory(tenantCode)
	if err != nil {
		return nil, err
	}

	user, err := repo.GetByUsername(req.Username)
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

	hash := hashToken(rToken)

	refreshRepo, err := s.refreshTokenRepoFactory(tenantCode)
	if err != nil {
		return nil, err
	}

	err = refreshRepo.Create(&models.RefreshToken{
		UserID:    user.ID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	})
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"access_token":  aToken,
		"refresh_token": rToken,
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

	refreshRepo, err := s.refreshTokenRepoFactory(claims.TenantCode)
	if err != nil {
		return nil, err
	}

	storedRToken, err := refreshRepo.FindValidByHash(hashToken(rToken))
	if err != nil {
		return nil, errors.New("refresh token revoked")
	}

	//revoke old refresh token
	_ = refreshRepo.Revoke(storedRToken.ID)

	newAToken, _ := s.jwtManager.GenerateAccessToken(claims.UserID, claims.Username, claims.TenantCode)
	newRToken, _ := s.jwtManager.GenerateRefreshToken(claims.UserID, claims.TenantCode)

	_ = refreshRepo.Create(&models.RefreshToken{
		UserID:    claims.UserID,
		TokenHash: hashToken(newRToken),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	})

	return map[string]string{
		"access_token":  newAToken,
		"refresh_token": newRToken,
	}, nil
}

func (s *authService) Logout(rToken string) error {
	claims, err := s.jwtManager.ParseToken(rToken)
	if err != nil || claims.Type != "refresh" {
		return errors.New("invalid refresh token")
	}
	refreshRepo, err := s.refreshTokenRepoFactory(claims.TenantCode)
	if err != nil {
		return err
	}
	return refreshRepo.RevokeAllByUser(claims.UserID)
}
