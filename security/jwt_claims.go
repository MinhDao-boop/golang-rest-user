package security

import (
	"golang-rest-user/enums"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username   string          `json:"username"`
	UserID     uint            `json:"user_id"`
	TenantCode string          `json:"tenant_code"`
	Version    int             `json:"ver"`
	Type       enums.TokenType `json:"type"`
	jwt.RegisteredClaims
}

type TokenResult struct {
	Token     string
	ExpiresIn int
}
