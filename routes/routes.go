package routes

import (
	"golang-rest-user/handler"
	"golang-rest-user/handler/tenant"

	"github.com/gin-gonic/gin"
)

func TenantRoutes(r *gin.RouterGroup) {
	r.GET("", handler.ListTenant)                   // GET /api/v1/tenants
	r.POST("", handler.CreateTenant)                // POST /api/v1/tenants
	r.GET("/:tenant_uuid", handler.GetByTenantCode) // GET /api/v1/tenants/:tenant_uuid
	r.PUT("/:tenant_uuid", handler.UpdateTenant)    // PUT /api/v1/tenants/:tenant_uuid
	r.DELETE("/:tenant_uuid", handler.DeleteTenant) // DELETE /api/v1/tenants/:tenant_uuid
}

func UserRoutes(r *gin.RouterGroup) {
	r.GET("", tenant.ListUsers)                // GET /api/v1/users
	r.POST("", tenant.CreateUser)              // POST /api/v1/users
	r.DELETE("", tenant.DeleteManyUsers)       // DELETE /api/v1/users
	r.GET("/:user_uuid", tenant.GetByUserUUID) // GET /api/v1/users/:user_uuid
	r.PUT("/:user_uuid", tenant.UpdateUser)    // PUT /api/v1/users/:user_uuid
}

func AuthRoutes(r *gin.RouterGroup) {
	r.POST("/register", tenant.Register) // POST /api/v1/auth/register
	r.POST("/login", tenant.Login)       // POST /api/v1/auth/login
	r.POST("/logout", tenant.Logout)     // POST /api/v1/auth/logout
	r.POST("/refresh", tenant.Refresh)   // POST /api/v1/auth/refresh
}

func ZonesRoutes(r *gin.RouterGroup) {
	r.GET("", tenant.ListZones)                // GET /api/v1/zones
	r.POST("", tenant.CreateZone)              // POST /api/v1/zones
	r.PUT("/:zone_uuid", tenant.UpdateZone)    // PUT /api/v1/zones/:uuid
	r.DELETE("/:zone_uuid", tenant.DeleteZone) // DELETE /api/v1/zones/:uuid
}

func ShareRoutes(r *gin.RouterGroup) {
	r.GET("/:zone_uuid", tenant.GetSharedUsers)              // GET /api/v1/zones/share/:zone_uuid
	r.POST("", tenant.ShareZone)                             // POST /api/v1/zones/share
	r.PUT("/:zone_uuid/:user_uuid", tenant.UpdatePermission) // PUT /api/v1/zones/share/:zone_uuid/:user_uuid
	r.DELETE("/:zone_uuid/:user_uuid", tenant.RevokeZone)    // DELETE /api/v1/zones/share/:zone_uuid/:user_uuid
}
