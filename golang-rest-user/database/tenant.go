package database

import (
	"fmt"
	"log"
	"time"

	"golang-rest-user/config"
	"golang-rest-user/models"
	"golang-rest-user/security"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitTenantDBs(masterDB *gorm.DB) error {
	var tenants []models.Tenant

	if err := masterDB.Where("status = ?", "active").Find(&tenants).Error; err != nil {
		return err
	}

	for _, t := range tenants {
		dbUser, err := security.Decrypt(t.DBUser)
		if err != nil {
			log.Printf("❌ tenant %s decrypt db user failed", t.Code)
			continue
		}
		dbPass, err := security.Decrypt(t.DBPass)
		if err != nil {
			log.Printf("❌ tenant %s decrypt db pass failed", t.Code)
			continue
		}
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
			dbUser,
			dbPass,
			t.DBHost,
			t.DBPort,
			t.DBName,
		)

		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Printf("❌ tenant %s connect failed", t.Code)
			continue
		}

		SetTenantDB(t.Code, db)
		log.Printf("✅ tenant %s connected", t.Code)
	}

	return nil
}

func CheckConnectMasterDB(tenant models.Tenant) (bool, error) {
	cfg := config.LoadConfig()
	//Decrypt
	dbUser, err := security.Decrypt(tenant.DBUser)
	if err != nil {
		return false, err
	}
	dbPass, err := security.Decrypt(tenant.DBPass)
	if err != nil {
		return false, err
	}
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		dbUser,
		dbPass,
		tenant.DBHost,
		tenant.DBPort,
		cfg.DBName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("❌ tenant %s connect to master DB failed", tenant.Code)
		return false, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return false, err
	}
	defer sqlDB.Close()
	return true, nil
}

func PingDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	return sqlDB.Ping()
}

func CreateTenantDatabase(dbName string) error {
	var masterDB = DB
	return masterDB.Exec("CREATE DATABASE IF NOT EXISTS " + dbName).Error
}

func ConnectTenantDB(tenant models.Tenant) (*gorm.DB, error) {
	//Decrypt db user
	dbUser, err := security.Decrypt(tenant.DBUser)
	if err != nil {
		return nil, err
	}
	//Decrypt db password
	dbPass, err := security.Decrypt(tenant.DBPass)
	if err != nil {
		return nil, err
	}
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		dbUser,
		dbPass,
		tenant.DBHost,
		tenant.DBPort,
		tenant.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
	)
}

func CloseTenantDB(oldDB *gorm.DB) error {
	if oldDB != nil {
		sqlDB, err := oldDB.DB()
		if err != nil {
			return err
		}
		sqlDB.Close()
	}
	return nil
}

func DropTenantDatabase(dbName string) error {
	var masterDB = DB
	return masterDB.Exec("DROP DATABASE IF EXISTS " + dbName).Error
}
