package repository

import (
	"errors"
	"golang-rest-user/models"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ZoneEscapeLinkRepo interface {
	GetByZoneID(zoneID uint) (*models.ZoneEscapeLink, error)
	Create(*models.ZoneEscapeLink) error
	Update(zoneID uint, webview datatypes.JSON) error
	FindByZoneID(zoneID uint) (*models.ZoneEscapeLink, error)
}

type ZoneEscapeLinkRepoImpl struct {
	db *gorm.DB
}

func (r *ZoneEscapeLinkRepoImpl) GetByZoneID(zoneID uint) (*models.ZoneEscapeLink, error) {
	var record models.ZoneEscapeLink
	if err := r.db.Where("zone_id = ?", zoneID).First(&record).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &record, nil
}

func (r *ZoneEscapeLinkRepoImpl) Create(newLink *models.ZoneEscapeLink) error {
	return r.db.Create(&newLink).Error
}

func (r *ZoneEscapeLinkRepoImpl) Update(zoneID uint, webview datatypes.JSON) error {
	return r.db.Model(&models.ZoneEscapeLink{}).Where("zone_id = ?", zoneID).
		Updates(models.ZoneEscapeLink{WebView: webview}).Error
}

func (r *ZoneEscapeLinkRepoImpl) FindByZoneID(zoneID uint) (*models.ZoneEscapeLink, error) {
	var record models.ZoneEscapeLink
	if err := r.db.Where("zone_id = ?", zoneID).First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func NewZoneEscapeLinkRepository(db *gorm.DB) ZoneEscapeLinkRepo {
	return &ZoneEscapeLinkRepoImpl{db: db}
}
