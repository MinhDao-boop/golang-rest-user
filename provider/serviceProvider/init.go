package serviceProvider

import (
	"golang-rest-user/provider/mySqlProvider"
	"golang-rest-user/repository"
	"golang-rest-user/security"
	"golang-rest-user/service"
)

type AppService struct {
	TenantService service.TenantService
	JWTManager    *security.Manager
}

var instance *AppService

func Init() {
	instance = &AppService{}
	masterDB := mySqlProvider.GetInstance()

	tenantRepo := repository.NewTenantRepo(masterDB)
	instance.TenantService = service.NewTenantService(tenantRepo)

	jwtConfig := security.LoadJWTConfig()
	instance.JWTManager = security.NewManager(jwtConfig)
}

func GetInstance() *AppService {
	return instance
}
