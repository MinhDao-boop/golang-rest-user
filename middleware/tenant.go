package middleware

import (
	//"log"
	"net/http"

	"golang-rest-user/database"
	"golang-rest-user/response"

	"github.com/gin-gonic/gin"
)

func TenantDBMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantCode := c.GetHeader("X-Tenant-Code")
		if tenantCode == "" {
			response.Error(c, response.CodeBadRequest, "X-Tenant-Code header is required", nil, http.StatusInternalServerError)
			return
		}

		db, ok := database.GetTenantDB(tenantCode)
		if !ok {
			response.Error(c, response.CodeBadRequest, "Invalid tenant code", nil, http.StatusInternalServerError)
			return
		}

		c.Set("TENANT_DB", db)
		c.Next()
	}
}
