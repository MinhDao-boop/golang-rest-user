package main

import (
	"golang-rest-user/config"
	"golang-rest-user/database"
	"golang-rest-user/handler"
	"golang-rest-user/models"

	"golang-rest-user/repository"
	"golang-rest-user/routes"
	"golang-rest-user/service"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	config.InitRedis()
	masterDB := database.ConnectMasterDB()
	masterDB.AutoMigrate(&models.Tenant{})
	if err := database.InitTenantDBs(masterDB); err != nil {
		log.Fatal(err)
	}
	r := gin.Default()

	tntRepo := repository.NewTenantRepo(masterDB)
	tntSvc := service.NewTenantService(tntRepo)
	tntHandler := handler.NewTenantHandler(tntSvc)

	userHandler := handler.NewUserHandler()

	authHandler := handler.NewAuthHandler()

	routes.RegisterRoutes(r, userHandler, tntHandler, authHandler)

	r.Run(":8080")

}
