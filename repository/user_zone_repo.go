package repository

import (
	"golang-rest-user/enums"
	"golang-rest-user/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserZoneRepo interface {
	Create(*models.UserZone) error
	UpdatePermission(userID, zoneID uint, permission enums.UserPermission) error
	Delete(userID uint, zoneID []uint) (int64, error)
	GetPermission(userID, zoneId uint) (*uint, error)
	GetSharedUser(uint) ([]models.UserZone, error)
	GetByUserID(uint) ([]models.UserZone, error)
	CountOwnerPermission(shareUserID uint, paths []string) (int64, error)
	BatchInsert([]models.UserZone) error
	GetUserZone(userID uint, zoneID uint) (*models.UserZone, error)
}

type userZoneRepoImpl struct {
	db *gorm.DB
}

func (r *userZoneRepoImpl) GetUserZone(userID uint, zoneID uint) (*models.UserZone, error) {
	var userZone models.UserZone
	if err := r.db.Model(&models.UserZone{}).
		Where("user_id = ? AND zone_id = ?", userID, zoneID).First(&userZone).Error; err != nil {
		return nil, err
	}
	return &userZone, nil
}

func (r *userZoneRepoImpl) BatchInsert(userZones []models.UserZone) error {
	return r.db.
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "zone_uuid"}},
			DoUpdates: clause.AssignmentColumns([]string{"permission"}),
		}).
		Create(userZones).Error
}

func (r *userZoneRepoImpl) CountOwnerPermission(shareUserID uint, paths []string) (int64, error) {
	var count int64
	err := r.db.
		Model(&models.UserZone{}).
		Table("user_zones uz").
		Joins("JOIN zones z on z.id = uz.zone_id").
		Where("uz.user_id = ?", shareUserID).
		Where("uz.permission = ?", enums.UserViewer).
		Scopes(func(db *gorm.DB) *gorm.DB {
			for _, path := range paths {
				db = db.Or("? LIKE CONCAT(z.path, '%')", path)
			}
			return db
		}).Distinct("z.id").Count(&count).Error
	return count, err
}

func (r *userZoneRepoImpl) GetByUserID(userID uint) (userZones []models.UserZone, err error) {
	if err = r.db.Where("user_id = ?", userID).
		Find(&userZones).Error; err != nil {
		return nil, err
	}
	return userZones, nil
}

func (r *userZoneRepoImpl) GetPermission(userID, zoneId uint) (*uint, error) {
	var (
		permission *uint
		path       string
	)
	if err := r.db.Model(&models.Zone{}).Select("path").Where("id = ?", zoneId).Scan(&path).Error; err != nil {
		return nil, err
	}
	err := r.db.Table("user_zones uz").
		Select("uz.permission").
		Joins("JOIN zones z on uz.zone_id = z.id").
		Where("uz.user_id = ? AND ? LIKE CONCAT(z.path, '%')", userID, path).
		Order("z.level DESC").
		Limit(1).Scan(&permission).Error
	if err != nil {
		return nil, err
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
	return r.db.Save(userZone).Error
}

func (r *userZoneRepoImpl) UpdatePermission(userID, zoneID uint, permission enums.UserPermission) error {
	return r.db.Model(&models.UserZone{}).Where("user_id = ? AND zone_id = ?", userID, zoneID).
		Update("permission", permission).Error
}

func (r *userZoneRepoImpl) Delete(userID uint, zoneID []uint) (int64, error) {
	res := r.db.Unscoped().Where("user_id = ? AND zone_id IN ?", userID, zoneID).Delete(&models.UserZone{})
	return res.RowsAffected, res.Error
}

func NewUserZoneRepo(db *gorm.DB) UserZoneRepo {
	return &userZoneRepoImpl{db: db}
}
