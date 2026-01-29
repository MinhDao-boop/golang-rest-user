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
	"golang-rest-user/security"
)

type AuthService interface {
	Register(req dto.CreateUserRequest) (*dto.UserResponse, error)
	Login(req dto.LoginRequest) (map[string]interface{}, error)
	Refresh(refreshToken string) (map[string]interface{}, error)
	Logout(refreshToken string) error
}

type authServiceImpl struct {
	userSvc    UserService
	jwtManager *security.Manager
	tenantCode string
}

func NewAuthService(userSvc UserService, jwtManager *security.Manager, tenantCode string) AuthService {
	return &authServiceImpl{
		userSvc:    userSvc,
		jwtManager: jwtManager,
		tenantCode: tenantCode,
	}
}

func (s *authServiceImpl) Register(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	userResponse, err := s.userSvc.Create(req)
	if err != nil {
		return nil, err
	}
	return userResponse, nil
}

func (s *authServiceImpl) Login(req dto.LoginRequest) (map[string]interface{}, error) {

	user, err := s.userSvc.GetByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	decryptedPass, _ := utils.AESGCMDecrypt(user.Password)
	if decryptedPass != req.Password {
		return nil, errors.New("invalid credentials")
	}

	ver := redisProvider.GetTokenVer(user.ID, s.tenantCode)

	aToken, err := s.jwtManager.GenerateToken(user.ID, user.Username, s.tenantCode, enums.TokenTypeAccess, 3600, ver)
	if err != nil {
		return nil, err
	}

	rToken, err := s.jwtManager.GenerateToken(user.ID, user.Username, s.tenantCode, enums.TokenTypeRefresh, 604800, ver)
	if err != nil {
		return nil, err
	}

	hash := hashToken(rToken.Token)
	ttl := time.Duration(rToken.ExpiresIn) * time.Second

	err = redisProvider.Create(hash, user.ID, s.tenantCode, ttl)
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

func (s *authServiceImpl) Refresh(rToken string) (map[string]interface{}, error) {
	claims, err := s.jwtManager.ParseToken(rToken)

	if claims == nil {
		return nil, errors.New("invalid token")
	}

	if err != nil || claims.Type != enums.TokenTypeRefresh {
		return nil, errors.New("invalid refresh token")
	}

	if s.tenantCode != claims.TenantCode {
		return nil, errors.New("invalid tenant code")
	}

	if err := redisProvider.FindValidByHash(hashToken(rToken), claims.TenantCode, claims.UserID); err != nil {
		return nil, errors.New("refresh token revoked")
	}

	ver := redisProvider.GetTokenVer(claims.UserID, claims.TenantCode)

	//revoke old refresh token
	if err = redisProvider.Revoke(hashToken(rToken), claims.TenantCode, claims.UserID); err != nil {
		return nil, err
	}

	newAToken, _ := s.jwtManager.GenerateToken(claims.UserID, claims.Username, claims.TenantCode, enums.TokenTypeAccess, 3600, ver)
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

func (s *authServiceImpl) Logout(rToken string) error {
	claims, err := s.jwtManager.ParseToken(rToken)
	if claims == nil {
		return errors.New("invalid token")
	}
	if err != nil || claims.Type != enums.TokenTypeRefresh {
		return errors.New("invalid refresh token")
	}
	if s.tenantCode != claims.TenantCode {
		return errors.New("invalid tenant code")
	}

	if err := redisProvider.IncreaseTokenVer(claims.UserID, claims.TenantCode); err != nil {
		return err
	}

	return redisProvider.RevokeAllByUser(claims.TenantCode, claims.UserID)
}
