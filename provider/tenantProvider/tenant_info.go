package tenantProvider

import (
	"golang-rest-user/database"
	"golang-rest-user/models"
	"golang-rest-user/repository"
	"golang-rest-user/security"
	"golang-rest-user/service/tenantSvc"
	"golang-rest-user/utils"
	"log"
	"time"

	"gorm.io/gorm"
)

type TenantInfo struct {
	Info        *models.Tenant
	db          *gorm.DB
	UserService tenantSvc.UserService
	AuthService tenantSvc.AuthService
}

func (t *TenantInfo) Init() error {
	t.Info.DBUser, _ = utils.AESGCMDecrypt(t.Info.DBUser)
	t.Info.DBPass, _ = utils.AESGCMDecrypt(t.Info.DBPass)

	err := database.CreateTenantDB(t.Info.DBUser, t.Info.DBPass, t.Info.DBHost, t.Info.DBPort, t.Info.DBName)
	if err != nil {
		log.Println(err)
		return err
	}

	t.db, _ = database.InitTenantDB(t.Info.DBUser, t.Info.DBPass, t.Info.DBHost, t.Info.DBPort, t.Info.DBName)

	sqlDB, err := t.db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	t.InitService()
	t.Migrate()
	return nil
}

func (t *TenantInfo) InitService() {

	userRepo := repository.NewUserRepo(t.db)
	t.UserService = tenantSvc.NewUserService(t.Info.Code, userRepo)

	jwtConfig := security.LoadJWTConfig()
	jwtManager := security.NewManager(jwtConfig)
	refreshTokenRedis := repository.NewRefreshTokenRedisRepo()
	t.AuthService = tenantSvc.NewAuthService(userRepo, refreshTokenRedis, jwtManager)
}

func (t *TenantInfo) Migrate() {

	err := t.db.AutoMigrate(&models.User{})
	if err != nil {
		log.Println(err)
	}
}
