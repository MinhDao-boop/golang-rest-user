package routesProvider

import (
	"golang-rest-user/middleware"
	"golang-rest-user/provider/serviceProvider"
	"golang-rest-user/routes"

	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	service := serviceProvider.GetInstance()
	jwtManager := service.JWTManager
	router.Use(gin.Recovery())

	router.Use(middleware.RequestID())

	v1 := router.Group("api/v1")

	tenants := v1.Group("/tenants")
	routes.TenantRoutes(tenants)

	auth := v1.Group("/auth")
	auth.Use(middleware.TenantDBMiddleware())
	routes.AuthRoutes(auth)

	users := v1.Group("/users")
	users.Use(middleware.AuthMiddleware(jwtManager))
	routes.UserRoutes(users)

	zones := v1.Group("/zones")
	zones.Use(middleware.AuthMiddleware(jwtManager))
	routes.ZonesRoutes(zones)
}
