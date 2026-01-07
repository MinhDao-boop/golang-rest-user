package mySqlProvider

import (
	"database/sql"
	"fmt"
	"golang-rest-user/models"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var instance *gorm.DB

func CreateInstanceDB(dbUser, dbPass, dbHost, dbPort, dbName string) (instance *gorm.DB, err error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName,
	)

	if instance, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		log.Printf("Error connect db %s", dbName)
		return nil, err
	}

	log.Printf("db %s connected", dbName)

	return instance, err
}

func CreateDB(dbUser, dbPass, dbHost, dbPort, dbName string) error {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/information_schema?charset=utf8mb4&parseTime=true",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Printf("Error connect to DB %s", dsn)
		return err
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		log.Printf("Error creating DB %s", dsn)
	}
	return nil
}

func Init() {
	var err error
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	_ = CreateDB(dbUser, dbPass, dbHost, dbPort, dbName)

	if instance, err = CreateInstanceDB(dbUser, dbPass, dbHost, dbPort, dbName); err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	if err = instance.AutoMigrate(&models.Tenant{}); err != nil {
		log.Fatalf("failed to auto migrate tenant: %v", err)
	}
}

func GetInstance() *gorm.DB {
	return instance
}
