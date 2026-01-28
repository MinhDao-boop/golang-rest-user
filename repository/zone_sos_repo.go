package repository

import (
	"golang-rest-user/models"

	"gorm.io/gorm"
)

type ZoneSOSRepo interface {
	Create(*models.ZoneSOS) error
	Delete([]uint, uint) (int64, error)
	GetByContactId(uint) (*models.ZoneSOS, error)
}

type ZoneSOSRepoImpl struct {
	db *gorm.DB
}

func (r *ZoneSOSRepoImpl) GetByContactId(contactId uint) (*models.ZoneSOS, error) {
	var contact models.ZoneSOS
	if err := r.db.Model(&models.ZoneSOS{}).Where("id = ?", contactId).First(&contact).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

func (r *ZoneSOSRepoImpl) Create(zoneContact *models.ZoneSOS) error {
	return r.db.Create(zoneContact).Error
}

func (r *ZoneSOSRepoImpl) Delete(contactIds []uint, zoneId uint) (int64, error) {
	res := r.db.Where("sos_contact_id IN (?) AND zone_id = ?", contactIds, zoneId).Delete(&models.ZoneSOS{})
	return res.RowsAffected, res.Error
}

func NewZoneSOSRepo(db *gorm.DB) ZoneSOSRepo {
	return &ZoneSOSRepoImpl{db: db}
}
