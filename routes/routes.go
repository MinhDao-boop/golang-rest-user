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
		auth := v1.Group("/auth")
		auth.Use(middleware.TenantContextMiddleware())
		{
			auth.POST("/register", a.Register) // POST /api/v1/auth/register
			auth.POST("/login", a.Login)       // POST /api/v1/auth/login
			auth.POST("/logout", a.Logout)     // POST /api/v1/auth/logout
			auth.POST("refresh", a.Refresh)    // POST /api/v1/auth/refresh
		}
		tenants := v1.Group("/tenants")
		{
			tenants.GET("", t.ListTenantResponse)   // GET /api/v1/tenants
			tenants.POST("", t.CreateTenantRequest) // POST /api/v1/tenants
			tenants.GET("/:id", t.GetByTenantID)    // GET /api/v1/tenants/:id
			tenants.PUT("/:id", t.UpdateTenant)     // PUT /api/v1/tenants/:id
			tenants.DELETE("/:id", t.DeleteTenant)  // DELETE /api/v1/tenants/:id
		}

		users := v1.Group("/users")
		users.Use(middleware.TenantDBMiddleware())
		{
			users.GET("", u.ListUsersResponse)     // GET /api/v1/users
			users.POST("", u.CreateUser)           // POST /api/v1/users
			users.DELETE("", u.DeleteManyUsers)    // DELETE /api/v1/users?ids=1,2,3
			users.GET("/:id", u.GetByUserID)       // GET /api/v1/users/:id
			users.PUT("/:id", u.UpdateUserRequest) // PUT /api/v1/users/:id
			users.DELETE("/:id", u.DeleteUser)     // DELETE /api/v1/users/:id
		}
	}
}
