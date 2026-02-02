package service

import (
	"crypto/sha256"
	"encoding/hex"
	"golang-rest-user/enums"
	"golang-rest-user/provider/redisProvider"
	"golang-rest-user/response"
	"golang-rest-user/utils"
	"time"

	"golang-rest-user/dto"
	"golang-rest-user/security"
)

type AuthService interface {
	Register(req dto.CreateUserRequest) *response.Response
	Login(req dto.LoginRequest) *response.Response
	Refresh(refreshToken string) *response.Response
	Logout(refreshToken string) *response.Response
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

func (s *authServiceImpl) Register(req dto.CreateUserRequest) *response.Response {
	return s.userSvc.Create(req)
}

func (s *authServiceImpl) Login(req dto.LoginRequest) *response.Response {
	newResponse := response.NewResponse()
	user, err := s.userSvc.GetByUsername(req.Username)
	if err != nil {
		newResponse.Err = response.ErrFailedAuthentication
		newResponse.MessageCode = response.AUT0003
		return newResponse
	}

	decryptedPass, _ := utils.AESGCMDecrypt(user.Password)
	if decryptedPass != req.Password {
		newResponse.Err = response.ErrFailedAuthentication
		newResponse.MessageCode = response.AUT0003
		return newResponse
	}

	ver := redisProvider.GetTokenVer(user.ID, s.tenantCode)

	accessToken, err := s.jwtManager.GenerateToken(user.ID, user.Username, s.tenantCode, enums.TokenTypeAccess, 3600, ver)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}

	refreshToken, err := s.jwtManager.GenerateToken(user.ID, user.Username, s.tenantCode, enums.TokenTypeRefresh, 604800, ver)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}

	hash := hashToken(refreshToken.Token)
	ttl := time.Duration(refreshToken.ExpiresIn) * time.Second

	err = redisProvider.Create(hash, user.ID, s.tenantCode, ttl)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	newResponse.Data = &response.TokenResponse{
		AccessToken:      accessToken.Token,
		AccessExpiresIn:  accessToken.ExpiresIn,
		RefreshToken:     refreshToken.Token,
		RefreshExpiresIn: refreshToken.ExpiresIn,
	}
	newResponse.MessageCode = response.SUS0000
	return newResponse
}

func hashToken(rToken string) string {
	h := sha256.Sum256([]byte(rToken))
	return hex.EncodeToString(h[:])
}

func (s *authServiceImpl) Refresh(rToken string) *response.Response {
	newResponse := response.NewResponse()
	claims, err := s.jwtManager.ParseToken(rToken)

	if claims == nil {
		newResponse.Err = response.ErrInvalidToken
		return newResponse
	}

	if err != nil || claims.Type != enums.TokenTypeRefresh {
		newResponse.Err = response.ErrInvalidToken
		return newResponse
	}

	if s.tenantCode != claims.TenantCode {
		newResponse.Err = response.ErrInvalidToken
		return newResponse
	}

	if err = redisProvider.FindValidByHash(hashToken(rToken), claims.TenantCode, claims.UserID); err != nil {
		newResponse.Err = err
		return newResponse
	}

	ver := redisProvider.GetTokenVer(claims.UserID, claims.TenantCode)

	//revoke old refresh token
	if err = redisProvider.Revoke(hashToken(rToken), claims.TenantCode, claims.UserID); err != nil {
		newResponse.Err = err
		return newResponse
	}

	newAccessToken, _ := s.jwtManager.GenerateToken(claims.UserID, claims.Username, claims.TenantCode, enums.TokenTypeAccess, 3600, ver)
	newRefreshToken, _ := s.jwtManager.GenerateToken(claims.UserID, claims.Username, claims.TenantCode, enums.TokenTypeRefresh, 604800, ver)

	hash := hashToken(newRefreshToken.Token)
	ttl := time.Duration(newRefreshToken.ExpiresIn) * time.Second

	if err = redisProvider.Create(hash, claims.UserID, claims.TenantCode, ttl); err != nil {
		newResponse.Err = err
		return newResponse
	}
	newResponse.Data = &response.TokenResponse{
		AccessToken:      newAccessToken.Token,
		AccessExpiresIn:  newAccessToken.ExpiresIn,
		RefreshToken:     newRefreshToken.Token,
		RefreshExpiresIn: newRefreshToken.ExpiresIn,
	}
	newResponse.MessageCode = response.SUS0000
	return newResponse
}

func (s *authServiceImpl) Logout(rToken string) *response.Response {
	newResponse := response.NewResponse()
	claims, err := s.jwtManager.ParseToken(rToken)
	if claims == nil {
		newResponse.Err = response.ErrInvalidToken
		return newResponse
	}
	if err != nil || claims.Type != enums.TokenTypeRefresh {
		newResponse.Err = response.ErrInvalidToken
		return newResponse
	}
	if s.tenantCode != claims.TenantCode {
		newResponse.Err = response.ErrInvalidToken
		return newResponse
	}

	if err = redisProvider.IncreaseTokenVer(claims.UserID, claims.TenantCode); err != nil {
		newResponse.Err = err
		return newResponse
	}

	if err = redisProvider.RevokeAllByUser(claims.TenantCode, claims.UserID); err != nil {
		newResponse.Err = err
		return newResponse
	}
	newResponse.Data = map[string]string{"message": "logged out"}
	newResponse.MessageCode = response.SUS0000
	return newResponse
}
