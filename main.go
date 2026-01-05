package main

import (
	"golang-rest-user/config"
	"golang-rest-user/database"
	"golang-rest-user/handler"
	"golang-rest-user/models"
	"golang-rest-user/provider/tenantProvider"
	"golang-rest-user/service/tenantSvc"

	"golang-rest-user/repository"
	"golang-rest-user/routes"
	"golang-rest-user/service"

	"github.com/gin-gonic/gin"
)

func main() {

	config.InitRedis()
	masterDB := database.ConnectMasterDB()
	err := masterDB.AutoMigrate(&models.Tenant{})
	if err != nil {
		return
	}
	r := gin.Default()

	tntRepo := repository.NewTenantRepo(masterDB)
	tntSvc := service.NewTenantService(tntRepo)
	tntCntSvc := tenantSvc.NewTenantConnect(tntRepo)
	tntHandler := handler.NewTenantHandler(tntSvc)

	tenantProvider.Init(tntCntSvc)

	userHandler := handler.NewUserHandler()
	authHandler := handler.NewAuthHandler()

	routes.RegisterRoutes(r, userHandler, tntHandler, authHandler)

	r.Run(":8080")

}
