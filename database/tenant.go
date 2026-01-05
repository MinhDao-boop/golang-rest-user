package database

import (
	"database/sql"
	"fmt"
	"golang-rest-user/config"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitTenantDB(dbUser, dbPass, dbHost, dbPort, dbName string) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Error connect tenant db %s", dbName)
		return nil, err
	}

	log.Printf("tenant db %s connected", dbName)

	return db, nil
}

func CreateTenantDB(dbUser, dbPass, dbHost, dbPort, dbName string) error {
	cfg := config.LoadConfig()
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		cfg.DBName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Connect to master DB failed")
		return err
	}
	db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			log.Printf("Close DB failed")
		}
	}(sqlDB)

	return nil
}
