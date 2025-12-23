package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	jwtConfig *JWTConfig
}

func NewManager(jwtConfig *JWTConfig) *Manager {
	return &Manager{jwtConfig: jwtConfig}
}

func (m *Manager) GenerateAccessToken(userID uint, username, tenantCode string) (string, error) {
	claims := &Claims{
		UserID:     userID,
		Username:   username,
		TenantCode: tenantCode,
		Type:       "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.jwtConfig.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.jwtConfig.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.jwtConfig.SecretKey)
}

func (m *Manager) GenerateRefreshToken(userID uint, tenantCode string) (string, error) {
	rClaims := &Claims{
		UserID:     userID,
		TenantCode: tenantCode,
		Type:       "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.jwtConfig.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.jwtConfig.RefreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	rToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rClaims)
	return rToken.SignedString(m.jwtConfig.SecretKey)
}

func (m *Manager) ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return m.jwtConfig.SecretKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
