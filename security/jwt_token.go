package security

import (
	"golang-rest-user/enums"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Manager struct {
	jwtConfig *JWTConfig
}

func NewManager(jwtConfig *JWTConfig) *Manager {
	return &Manager{jwtConfig: jwtConfig}
}

func (m *Manager) GenerateToken(userID uint, username, tenantCode string, tokenType enums.TokenType, ttl, ver int) (*TokenResult, error) {
	jti, _ := uuid.NewUUID()
	claims := &Claims{
		Username:   username,
		UserID:     userID,
		TenantCode: tenantCode,
		Type:       tokenType,
		Version:    ver,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.jwtConfig.Issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(ttl) * time.Second)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        jti.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.jwtConfig.SecretKey)
	if err != nil {
		return nil, err
	}
	return &TokenResult{
		Token:     signed,
		ExpiresIn: ttl,
	}, nil
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
