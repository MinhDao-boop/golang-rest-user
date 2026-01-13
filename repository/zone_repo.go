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
	GetByName(string) (*models.Zone, error)
	UpdateZonePath(uint, string) error
	GetSubtreeByPath(path string) ([]models.Zone, error)
}

type zoneRepoImpl struct {
	db *gorm.DB
}

func (r *zoneRepoImpl) GetSubtreeByPath(path string) ([]models.Zone, error) {
	var zones []models.Zone
	err := r.db.Where("path LIKE ?", path+"%").Order("level ASC").Find(&zones).Error
	if err != nil {
		return nil, err
	}
	return zones, nil
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

func (r *zoneRepoImpl) GetByName(name string) (*models.Zone, error) {
	var zone models.Zone
	if err := r.db.Where("name = ?", name).First(&zone).Error; err != nil {
		return nil, err
	}
	return &zone, nil
}

func (r *zoneRepoImpl) UpdateZonePath(newZoneID uint, newZonePath string) error {
	return r.db.Model(&models.Zone{}).Where("id = ?", newZoneID).
		Update("path", newZonePath).Error
}

func (r *zoneRepoImpl) GetByPath(path string) ([]models.Zone, error) {
	var zones []models.Zone
	if err := r.db.Where("path LIKE ?", path+"%").Order("level ASC").Find(&zones).Error; err != nil {
		return nil, err
	}
	return zones, nil
}
