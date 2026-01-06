package mySqlProvider

import (
	"golang-rest-user/database"
	"golang-rest-user/models"
	"log"

	"gorm.io/gorm"
)

var instance *gorm.DB

func Init() {
	instance = &gorm.DB{}
	instance = database.ConnectMasterDB()
	if err := instance.AutoMigrate(&models.Tenant{}); err != nil {
		log.Fatalf("failed to auto migrate tenant: %v", err)
	}
}

func GetInstance() *gorm.DB {
	return instance
}
