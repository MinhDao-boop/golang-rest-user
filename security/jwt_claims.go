package security

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID     uint   `json:"user_id"`
	Username   string `json:"username"`
	TenantCode string `json:"tenant_code"`
	Type       string `json:"type"`
	jwt.RegisteredClaims
}
