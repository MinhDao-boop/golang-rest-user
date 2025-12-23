package main

import (
	"golang-rest-user/database"
	"golang-rest-user/handler"
	"golang-rest-user/security"

	//"golang-rest-user/middleware"

	//"golang-rest-user/models"
	"golang-rest-user/repository"
	"golang-rest-user/routes"
	"golang-rest-user/service"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	database.ConnectMasterDB()
	DB := database.DB
	if err := database.InitTenantDBs(database.DB); err != nil {
		log.Fatal(err)
	}
	r := gin.Default()

	tntRepo := repository.NewTenantRepo(DB)
	tntSvc := service.NewTenantService(tntRepo)
	tntHandler := handler.NewTenantHandler(tntSvc)

	userHandler := handler.NewUserHandler()

	userRepoFactory := func(tenantCode string) (repository.UserRepo, error) {
		db, ok := database.GetTenantDB(tenantCode)
		if !ok {
			return nil, err
		}
		return repository.NewUserRepo(db), nil
	}

	refreshTokenRepoFactory := func(tenantCode string) (repository.RefreshTokenRepo, error) {
		db, ok := database.GetTenantDB(tenantCode)
		if !ok {
			return nil, err
		}
		return repository.NewRefreshTokenRepo(db), nil
	}

	jwtCfg := security.LoadJWTConfig()
	jwtManager := security.NewManager(jwtCfg)

	authService := service.NewAuthService(userRepoFactory, refreshTokenRepoFactory, jwtManager)
	authHandler := handler.NewAuthHandler(authService)

	routes.RegisterRoutes(r, userHandler, tntHandler, authHandler)

	r.Run(":8080")

}
