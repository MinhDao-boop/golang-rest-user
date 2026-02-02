package middleware

import (
	"golang-rest-user/response"

	"github.com/gin-gonic/gin"
)

func TenantDBMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantCode := c.GetHeader("X-Tenant-Code")
		if tenantCode == "" {
			response.Error(c, response.TEN0001, response.ErrMissingIdentifier)
			return
		}
		c.Set("TENANT_CODE", tenantCode)
		c.Next()
	}
}
