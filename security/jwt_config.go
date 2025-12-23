package security

import (
	"os"
	"time"
)

type JWTConfig struct {
	SecretKey       []byte
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Issuer          string
}

func LoadJWTConfig() *JWTConfig {
	return &JWTConfig{
		SecretKey:       []byte(os.Getenv("JWT_SECRET_KEY")),
		AccessTokenTTL:  time.Minute * 15,
		RefreshTokenTTL: time.Hour * 24 * 7,
		Issuer:          "golang-rest-user",
	}
}
