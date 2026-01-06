package routesProvider

import (
	"golang-rest-user/middleware"
	"golang-rest-user/provider/serviceProvider"
	"golang-rest-user/routes"

	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	appService := serviceProvider.GetInstance()
	router.Use(gin.Recovery())

	router.Use(middleware.RequestID())

	v1 := router.Group("api/v1")

	tenants := v1.Group("/tenants")
	routes.TenantRoutes(tenants)

	auth := v1.Group("/auth")
	auth.Use(middleware.TenantDBMiddleware())
	routes.AuthRoutes(auth)

	users := v1.Group("/users")
	users.Use(middleware.AuthMiddleware(appService.JWTManager))
	routes.UserRoutes(users)
}
