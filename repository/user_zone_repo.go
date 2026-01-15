package repository

import (
	"golang-rest-user/enums"
	"golang-rest-user/models"

	"gorm.io/gorm"
)

type UserZoneRepo interface {
	Create(*models.UserZone) error
	UpdatePermission(userID, zoneID uint, permission enums.UserPermission) error
	Delete(userID, zoneID uint) (int64, error)
	GetPermission(userID uint, path string) (string, error)
	GetZoneID(userID uint) (uint, error)
	GetSharedUser(uint) ([]models.UserZone, error)
	GetSharedZone(uint) ([]models.UserZone, error)
}

type userZoneRepoImpl struct {
	db *gorm.DB
}

func (r *userZoneRepoImpl) GetSharedZone(userID uint) (userZones []models.UserZone, err error) {
	if err = r.db.Where("user_id = ?", userID).
		Find(&userZones).Error; err != nil {
		return nil, err
	}
	return userZones, nil
}

func (r *userZoneRepoImpl) GetZoneID(userID uint) (uint, error) {
	var userZone models.UserZone
	err := r.db.Table("user_zones").Where("user_id = ?", userID).First(&userZone).Error

	return userZone.ZoneID, err
}

func (r *userZoneRepoImpl) GetPermission(userID uint, path string) (string, error) {
	var permission string
	err := r.db.Table("user_zones uz").
		Select("uz.permission").
		Joins("JOIN zones z on uz.zone_id = z.id").
		Where("uz.user_id = ? AND ? LIKE CONCAT(z.path, '%')", userID, path).
		Order("z.level DESC").
		Limit(1).Scan(&permission).Error
	if err != nil {
		return "", err
	}
	return permission, nil
}

func (r *userZoneRepoImpl) GetSharedUser(zoneID uint) (userZones []models.UserZone, err error) {
	if err = r.db.Where("zone_id = ?", zoneID).
		Find(&userZones).Error; err != nil {
		return nil, err
	}
	return userZones, nil
}

func (r *userZoneRepoImpl) Create(userZone *models.UserZone) error {
	return r.db.Create(userZone).Error
}

func (r *userZoneRepoImpl) UpdatePermission(userID, zoneID uint, permission enums.UserPermission) error {
	return r.db.Model(&models.UserZone{}).Where("user_id = ? AND zone_id = ?", userID, zoneID).
		Update("permission", permission).Error
}

func (r *userZoneRepoImpl) Delete(userID, zoneID uint) (int64, error) {
	res := r.db.Unscoped().Where("user_id = ? AND zone_id = ?", userID, zoneID).Delete(&models.UserZone{})
	return res.RowsAffected, res.Error
}

func NewUserZoneRepo(db *gorm.DB) UserZoneRepo {
	return &userZoneRepoImpl{db: db}
}
