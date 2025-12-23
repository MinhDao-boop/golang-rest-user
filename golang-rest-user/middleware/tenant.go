package middleware

import (
	//"log"
	"net/http"

	"golang-rest-user/database"
	"golang-rest-user/response"

	"github.com/gin-gonic/gin"
)

const ContextTenantCode = "resolved_tenant_code"

func TenantContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1️⃣ Ưu tiên tenant từ JWT (đã set bởi AuthMiddleware)
		if tenantCode, ok := c.Get(ContextTenantCode); ok {
			c.Set(ContextTenantCode, tenantCode)
			c.Next()
			return
		}

		// 2️⃣ Fallback: lấy từ header (login / register)
		headerCode := c.GetHeader("X-Tenant-Code")
		if headerCode != "" {
			c.Set(ContextTenantCode, headerCode)
			//log.Println("Tenant code from header:", headerCode)
			c.Next()
			return
		}

		response.Error(c, response.CodeBadRequest, "tenant code is required", nil, http.StatusBadRequest)
	}
}

func TenantDBMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantCode := c.GetHeader("X-Tenant-Code")
		if tenantCode == "" {
			response.Error(c, response.CodeBadRequest, "X-Tenant-Code header is required", nil, http.StatusBadRequest)
			return
		}

		db, ok := database.GetTenantDB(tenantCode)
		if !ok {
			response.Error(c, response.CodeBadRequest, "Invalid tenant code", nil, http.StatusBadRequest)
			return
		}

		c.Set("TENANT_DB", db)
		c.Next()
	}
}
