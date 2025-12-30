package routes

import (
	"golang-rest-user/handler"
	"golang-rest-user/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, u *handler.UserHandler, t *handler.TenantHandler, a *handler.AuthHandler) {
	r.Use(middleware.RequestID())
	v1 := r.Group("/api/v1")
	{
		tenants := v1.Group("/tenants")
		{
			tenants.GET("", t.ListTenant)            // GET /api/v1/tenants
			tenants.POST("", t.CreateTenant)         // POST /api/v1/tenants
			tenants.GET("/:code", t.GetByTenantCode) // GET /api/v1/tenants/:code
			tenants.PUT("/:code", t.UpdateTenant)    // PUT /api/v1/tenants/:code
			tenants.DELETE("/:code", t.DeleteTenant) // DELETE /api/v1/tenants/:code
		}

		users := v1.Group("/users")
		users.Use(middleware.TenantDBMiddleware())
		{
			users.GET("", u.ListUsers)           // GET /api/v1/users
			users.POST("", u.CreateUser)         // POST /api/v1/users
			users.DELETE("", u.DeleteManyUsers)  // DELETE /api/v1/users?uuids=1b0f0fe4-8710-4518-b8bc-7f1e52b280e4,1c8edc4f-b1a0-4252-808b-682eb76551ad,...
			users.GET("/:uuid", u.GetByUserUUID) // GET /api/v1/users/:uuid
			users.PUT("/:uuid", u.UpdateUser)    // PUT /api/v1/users/:uuid
		}
		auth := v1.Group("/auth")
		auth.Use(middleware.TenantDBMiddleware())
		{
			auth.POST("/register", a.Register) // POST /api/v1/auth/register
			auth.POST("/login", a.Login)       // POST /api/v1/auth/login
			auth.POST("/logout", a.Logout)     // POST /api/v1/auth/logout
			auth.POST("/refresh", a.Refresh)   // POST /api/v1/auth/refresh
		}
	}
}
