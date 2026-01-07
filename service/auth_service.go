package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"golang-rest-user/enums"
	"golang-rest-user/provider/redisProvider"
	"golang-rest-user/utils"
	"time"

	"golang-rest-user/dto"
	"golang-rest-user/models"
	"golang-rest-user/repository"
	"golang-rest-user/security"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(req dto.CreateUserRequest) (*dto.UserResponse, error)
	Login(tenantCode string, req dto.LoginRequest) (map[string]interface{}, error)
	Refresh(tenantCode, refreshToken string) (map[string]interface{}, error)
	Logout(tenantCode, refreshToken string) error
}

type authService struct {
	userRepo   repository.UserRepo
	jwtManager *security.Manager
}

func NewAuthService(userRepo repository.UserRepo, jwtManager *security.Manager) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (s *authService) Register(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	if _, err := s.userRepo.GetByUsername(req.Username); err == nil {
		return nil, errors.New("username already exists")
	}

	encryptedPass, err := utils.AESGCMEncrypt(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		UUID:     uuid.NewString(),
		Username: req.Username,
		Password: encryptedPass,
		FullName: req.FullName,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return ConvertToUserResponse(user), nil
}

func (s *authService) Login(tenantCode string, req dto.LoginRequest) (map[string]interface{}, error) {

	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	decryptedPass, _ := utils.AESGCMDecrypt(user.Password)
	if decryptedPass != req.Password {
		return nil, errors.New("invalid credentials")
	}

	ver := redisProvider.GetTokenVer(user.ID, tenantCode)

	aToken, err := s.jwtManager.GenerateToken(user.ID, user.Username, tenantCode, enums.TokenTypeAccess, 900, ver)
	if err != nil {
		return nil, err
	}

	rToken, err := s.jwtManager.GenerateToken(user.ID, user.Username, tenantCode, enums.TokenTypeRefresh, 604800, ver)
	if err != nil {
		return nil, err
	}

	hash := hashToken(rToken.Token)
	ttl := time.Duration(rToken.ExpiresIn) * time.Second

	err = redisProvider.Create(hash, user.ID, tenantCode, ttl)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"access_token":       aToken.Token,
		"access_expires_in":  aToken.ExpiresIn,
		"refresh_token":      rToken.Token,
		"refresh_expires_in": rToken.ExpiresIn,
	}, nil
}

func hashToken(rToken string) string {
	h := sha256.Sum256([]byte(rToken))
	return hex.EncodeToString(h[:])
}

func (s *authService) Refresh(tenantCode, rToken string) (map[string]interface{}, error) {
	claims, err := s.jwtManager.ParseToken(rToken)

	if claims == nil {
		return nil, errors.New("invalid token")
	}

	if err != nil || claims.Type != enums.TokenTypeRefresh {
		return nil, errors.New("invalid refresh token")
	}

	if tenantCode != claims.TenantCode {
		return nil, errors.New("invalid refresh token")
	}

	if err := redisProvider.FindValidByHash(hashToken(rToken), claims.TenantCode, claims.UserID); err != nil {
		return nil, errors.New("refresh token revoked")
	}

	ver := redisProvider.GetTokenVer(claims.UserID, claims.TenantCode)

	//revoke old refresh token
	if err = redisProvider.Revoke(hashToken(rToken), claims.TenantCode, claims.UserID); err != nil {
		return nil, err
	}

	newAToken, _ := s.jwtManager.GenerateToken(claims.UserID, claims.Username, claims.TenantCode, enums.TokenTypeAccess, 900, ver)
	newRToken, _ := s.jwtManager.GenerateToken(claims.UserID, claims.Username, claims.TenantCode, enums.TokenTypeRefresh, 604800, ver)

	hash := hashToken(newRToken.Token)
	ttl := time.Duration(newRToken.ExpiresIn) * time.Second

	if err = redisProvider.Create(hash, claims.UserID, claims.TenantCode, ttl); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"access_token":       newAToken.Token,
		"access_expires_in":  newAToken.ExpiresIn,
		"refresh_token":      newRToken.Token,
		"refresh_expires_in": newRToken.ExpiresIn,
	}, nil
}

func (s *authService) Logout(tenantCode, rToken string) error {
	claims, err := s.jwtManager.ParseToken(rToken)
	if claims == nil {
		return errors.New("invalid token")
	}
	if err != nil || claims.Type != enums.TokenTypeRefresh {
		return errors.New("invalid refresh token")
	}
	if tenantCode != claims.TenantCode {
		return errors.New("invalid refresh token")
	}

	if err := redisProvider.IncreaseTokenVer(claims.UserID, claims.TenantCode); err != nil {
		return err
	}

	return redisProvider.RevokeAllByUser(claims.TenantCode, claims.UserID)
}
