package main

import (
	"golang-rest-user/config"
	"golang-rest-user/database"
	"golang-rest-user/handler"
	"golang-rest-user/models"
	"golang-rest-user/security"

	"golang-rest-user/repository"
	"golang-rest-user/routes"
	"golang-rest-user/service"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	config.InitRedis()
	masterDB := database.ConnectMasterDB()
	err := masterDB.AutoMigrate(&models.Tenant{})
	if err != nil {
		return
	}
	if err := database.InitTenantDBs(masterDB); err != nil {
		log.Fatal(err)
	}
	r := gin.Default()

	tntRepo := repository.NewTenantRepo(masterDB)
	tntSvc := service.NewTenantService(tntRepo)
	tntHandler := handler.NewTenantHandler(tntSvc)

	userRepo := repository.NewUserRepo()
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	jwtConfig := security.LoadJWTConfig()
	jwtManager := security.NewManager(jwtConfig)
	refreshTokenRedis := repository.NewRefreshTokenRedisRepo()
	authSvc := service.NewAuthService(userRepo, refreshTokenRedis, jwtManager)
	authHandler := handler.NewAuthHandler(authSvc)

	routes.RegisterRoutes(r, userHandler, tntHandler, authHandler)

	r.Run(":8080")

}
