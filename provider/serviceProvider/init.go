package serviceProvider

import (
	"golang-rest-user/provider/mySqlProvider"
	"golang-rest-user/repository"
	"golang-rest-user/security"
	"golang-rest-user/service"
	"golang-rest-user/service/tenantSvc"
)

type AppService struct {
	TenantService        service.TenantService
	TenantConnectService tenantSvc.TenantConnect
	JWTManager           *security.Manager
}

var instance *AppService

func Init() {
	instance = &AppService{}
	masterDB := mySqlProvider.GetInstance()

	tenantRepo := repository.NewTenantRepo(masterDB)
	instance.TenantService = service.NewTenantService(tenantRepo)
	instance.TenantConnectService = tenantSvc.NewTenantConnect(tenantRepo)

	jwtConfig := security.LoadJWTConfig()
	instance.JWTManager = security.NewManager(jwtConfig)
}

func GetInstance() *AppService {
	return instance
}
