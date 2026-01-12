package repository

import (
	"golang-rest-user/models"

	"gorm.io/gorm"
)

type UserZoneRepo interface {
	AddUserPermission(*models.UserZone) error
	GetUserZones(userID uint) ([]models.UserZone, error)
}

type userZoneRepoImpl struct {
	db *gorm.DB
}

func NewUserZoneRepo(db *gorm.DB) UserZoneRepo {
	return &userZoneRepoImpl{db: db}
}

func (r *userZoneRepoImpl) AddUserPermission(userZone *models.UserZone) error {
	return r.db.Create(userZone).Error
}

func (r *userZoneRepoImpl) GetUserZones(userID uint) (userZones []models.UserZone, err error) {
	if err := r.db.Where("user_id = ?", userID).
		Find(&userZones).Error; err != nil {
		return nil, err
	}
	return userZones, nil
}
