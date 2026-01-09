package repository

import (
	"golang-rest-user/models"

	"gorm.io/gorm"
)

type ZoneRepo interface {
	Create(*models.Zone) error
	Update(*models.Zone) error
	DeleteByIDs([]uint) (deleted int64, err error)
	GetByID(uint) (*models.Zone, error)
	GetByPaths([]string) ([]models.Zone, error)
	UpdateZonePath(uint, string) error
	AddUserPermission(*models.UserZone) error
	GetUserZones(userID uint) ([]models.UserZone, error)
}

type zoneRepoImpl struct {
	db *gorm.DB
}

func NewZoneRepo(db *gorm.DB) ZoneRepo {
	return &zoneRepoImpl{db: db}
}

func (r *zoneRepoImpl) Create(user *models.Zone) error {
	return r.db.Create(user).Error
}

func (r *zoneRepoImpl) Update(user *models.Zone) error {
	return r.db.Save(user).Error
}

func (r *zoneRepoImpl) DeleteByIDs(ids []uint) (deleted int64, err error) {
	return 0, nil
}

func (r *zoneRepoImpl) GetByID(id uint) (*models.Zone, error) {
	var zone models.Zone
	if err := r.db.First(&zone, id).Error; err != nil {
		return nil, err
	}
	return &zone, nil
}

func (r *zoneRepoImpl) UpdateZonePath(newZoneID uint, newZonePath string) error {
	return r.db.Where("id = ?", newZoneID).
		Update("path", newZonePath).Error
}

func (r *zoneRepoImpl) AddUserPermission(userZone *models.UserZone) error {
	return r.db.Create(userZone).Error
}

func (r *zoneRepoImpl) GetUserZones(userID uint) (userZones []models.UserZone, err error) {
	if err := r.db.Where("user_id = ?", userID).
		Find(&userZones).Error; err != nil {
		return nil, err
	}
	return userZones, nil
}

func (r *zoneRepoImpl) GetByPaths(paths []string) ([]models.Zone, error) {
	var zones []models.Zone
	for _, path := range paths {
		var zone models.Zone
		if err := r.db.Where("path LIKE ?", path).First(&zone).Error; err != nil {
			return nil, err
		}
		zones = append(zones, zone)
	}
	return zones, nil
}
