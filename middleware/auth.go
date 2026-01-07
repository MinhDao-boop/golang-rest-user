package middleware

import (
	"golang-rest-user/enums"
	"golang-rest-user/provider/redisProvider"
	"golang-rest-user/response"
	"golang-rest-user/security"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtManager *security.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")

		if !strings.HasPrefix(auth, "Bearer ") {
			response.Error(c, response.CodeBadRequest, "Unauthorized", nil, http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		claims, err := jwtManager.ParseToken(tokenStr)
		if err != nil {
			response.Error(c, response.CodeBadRequest, "Unauthorized", nil, http.StatusUnauthorized)
			return
		}

		if err != nil || claims.Type != enums.TokenTypeAccess {
			response.Error(c, response.CodeBadRequest, "Invalid access token", nil, http.StatusUnauthorized)
			return
		}
		tokenVer := claims.Version
		currentVer := redisProvider.GetTokenVer(claims.UserID, claims.TenantCode)
		if tokenVer != currentVer {
			response.Error(c, response.CodeBadRequest, "Unauthorized", nil, http.StatusUnauthorized)
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("tenant_code", claims.TenantCode)

		c.Next()
	}
}
