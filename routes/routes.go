package routes

import (
	"golang-rest-user/handler"
	"golang-rest-user/handler/tenant"

	"github.com/gin-gonic/gin"
)

func TenantRoutes(r *gin.RouterGroup) {
	r.GET("", handler.ListTenant)            // GET /api/v1/tenants
	r.POST("", handler.CreateTenant)         // POST /api/v1/tenants
	r.GET("/:code", handler.GetByTenantCode) // GET /api/v1/tenants/:code
	r.PUT("/:code", handler.UpdateTenant)    // PUT /api/v1/tenants/:code
	r.DELETE("/:code", handler.DeleteTenant) // DELETE /api/v1/tenants/:code
}

func UserRoutes(r *gin.RouterGroup) {
	r.GET("", tenant.ListUsers)           // GET /api/v1/users
	r.POST("", tenant.CreateUser)         // POST /api/v1/users
	r.DELETE("", tenant.DeleteManyUsers)  // DELETE /api/v1/users?uuids=1b0f0fe4-8710-4518-b8bc-7f1e52b280e4,1c8edc4f-b1a0-4252-808b-682eb76551ad,...
	r.GET("/:uuid", tenant.GetByUserUUID) // GET /api/v1/users/:uuid
	r.PUT("/:uuid", tenant.UpdateUser)    // PUT /api/v1/users/:uuid
}

func AuthRoutes(r *gin.RouterGroup) {
	r.POST("/register", tenant.Register) // POST /api/v1/auth/register
	r.POST("/login", tenant.Login)       // POST /api/v1/auth/login
	r.POST("/logout", tenant.Logout)     // POST /api/v1/auth/logout
	r.POST("/refresh", tenant.Refresh)   // POST /api/v1/auth/refresh
}
