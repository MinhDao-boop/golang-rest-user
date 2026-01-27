package repository

import (
	"golang-rest-user/models"

	"gorm.io/gorm"
)

type ZoneSOSRepo interface {
	Create(*models.ZoneSOS) error
	Delete([]uint, uint) (int64, error)
}

type ZoneSOSRepoImpl struct {
	db *gorm.DB
}

func (r *ZoneSOSRepoImpl) Create(zoneContact *models.ZoneSOS) error {
	return r.db.Create(zoneContact).Error
}

func (r *ZoneSOSRepoImpl) Delete(contactIds []uint, zoneId uint) (int64, error) {
	res := r.db.Unscoped().Where("sos_contact_id IN (?) AND zone_id = ?", contactIds, zoneId).Delete(&models.ZoneSOS{})
	return res.RowsAffected, res.Error
}

func NewZoneSOSRepo(db *gorm.DB) ZoneSOSRepo {
	return &ZoneSOSRepoImpl{db: db}
}
