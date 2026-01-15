package tenantProvider

import (
	"golang-rest-user/models"
	"golang-rest-user/provider/mySqlProvider"
	"golang-rest-user/provider/serviceProvider"
	"golang-rest-user/repository"
	"golang-rest-user/service"
	"golang-rest-user/utils"
	"log"
	"time"

	"gorm.io/gorm"
)

type TenantInfo struct {
	Info         *models.Tenant
	db           *gorm.DB
	UserService  service.UserService
	AuthService  service.AuthService
	ZoneService  service.ZoneService
	ShareService service.ShareService
}

func (t *TenantInfo) Init() error {
	decryptedDBUser, _ := utils.AESGCMDecrypt(t.Info.DBUser)
	decryptedDBPass, _ := utils.AESGCMDecrypt(t.Info.DBPass)

	err := mySqlProvider.CreateDB(decryptedDBUser, decryptedDBPass, t.Info.DBHost, t.Info.DBPort, t.Info.DBName)
	if err != nil {
		log.Println(err)
		return err
	}

	t.db, _ = mySqlProvider.CreateInstanceDB(decryptedDBUser, decryptedDBPass, t.Info.DBHost, t.Info.DBPort, t.Info.DBName)

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
	appService := serviceProvider.GetInstance()

	userRepo := repository.NewUserRepo(t.db)
	t.UserService = service.NewUserService(t.Info.Code, userRepo)

	jwtManager := appService.JWTManager
	t.AuthService = service.NewAuthService(userRepo, jwtManager)

	zoneRepo := repository.NewZoneRepo(t.db)
	userZoneRepo := repository.NewUserZoneRepo(t.db)
	t.ZoneService = service.NewZoneService(zoneRepo, userZoneRepo)
	t.ShareService = service.NewShareService(userZoneRepo, zoneRepo, userRepo)
}

func (t *TenantInfo) Migrate() {

	err := t.db.AutoMigrate(
		&models.User{},
		&models.Zone{},
		&models.UserZone{},
	)
	if err != nil {
		log.Println(err)
	}
}

func (t *TenantInfo) Destruction() {
	db, err := t.db.DB()
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
}

func (t *TenantInfo) Drop() {
	db, err := t.db.DB()
	if err != nil {
		log.Println(err)
	}
	defer db.Close()
	db.Exec("DROP TABLE IF EXISTS %s", t.Info.DBName)
}
