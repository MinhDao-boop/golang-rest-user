package repository

import (
	"golang-rest-user/models"
	"golang-rest-user/utils"

	"gorm.io/gorm"
)

type SOSContactRepo interface {
	Create(*models.SOSContact) error
	UpdateMap(string, map[string]interface{}) error
	Updates(string, *models.SOSContact) error
	ToggleStatus(string, *models.SOSContact) error
	ListByZone(zoneId uint, page, pageSize int, search string, isAll bool, isActive *bool) ([]models.SOSContact, int64, error)
	DeleteMany([]uint, uint) (int64, error)
	GetByContactAndZone(string, uint) (*models.SOSContact, error)
	GetByPhone(interface{}) (*models.SOSContact, error)
}

type SOSContactRepoImpl struct {
	db *gorm.DB
}

func (r *SOSContactRepoImpl) ToggleStatus(uuid string, contact *models.SOSContact) error {
	return r.db.Model(&models.SOSContact{}).Where("uuid = ?", uuid).Update("is_active", contact.IsActive).Error
}

func (r *SOSContactRepoImpl) Updates(uuid string, contact *models.SOSContact) error {
	return r.db.Where("uuid = ?", uuid).Updates(contact).Error
}

func (r *SOSContactRepoImpl) ListByZone(zoneId uint, page, pageSize int, search string, isAll bool, isActive *bool) ([]models.SOSContact, int64, error) {
	var (
		sosContacts []models.SOSContact
		total       int64
	)
	query := r.db.
		Model(&models.SOSContact{}).
		Joins("JOIN zone_sos zs ON zs.sos_contact_id = sos_contacts.id").
		Where("zs.zone_id = ?", zoneId).
		Where("(name LIKE ? OR phone LIKE ?)", "%"+search+"%", "%"+search+"%")
	if isActive != nil {
		query.Where("sos_contacts.is_active = ?", *isActive)
	}
	if !isAll {
		query = query.Scopes(utils.PaginateGORM(page, pageSize))
	}
	if err := query.Find(&sosContacts).Error; err != nil {
		return nil, 0, err
	}
	if isAll {
		total = int64(len(sosContacts))
	}
	if err := r.db.Model(&models.ZoneSOS{}).Where("zone_id = ?", zoneId).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return sosContacts, total, nil
}

func (r *SOSContactRepoImpl) GetByPhone(phone interface{}) (*models.SOSContact, error) {
	var contact models.SOSContact
	if err := r.db.Model(&models.SOSContact{}).Where("phone = ?", phone).First(&contact).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

func (r *SOSContactRepoImpl) GetByContactAndZone(contactUuid string, zoneId uint) (*models.SOSContact, error) {
	var contact models.SOSContact
	if err := r.db.Model(&models.SOSContact{}).
		Joins("JOIN zone_sos zs ON zs.sos_contact_id = sos_contacts.id").
		Where("sos_contacts.uuid = ? AND zs.zone_id = ?", contactUuid, zoneId).First(&contact).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

func (r *SOSContactRepoImpl) Create(sosEmployee *models.SOSContact) error {
	return r.db.Create(sosEmployee).Error
}

func (r *SOSContactRepoImpl) UpdateMap(uuid string, updates map[string]interface{}) error {
	return r.db.Model(&models.SOSContact{}).Where("uuid = ?", uuid).
		Updates(updates).Error
}

func (r *SOSContactRepoImpl) DeleteMany(contactIds []uint, zoneId uint) (int64, error) {
	db := r.db

	subQuery := db.
		Table("zone_sos").
		Select("sos_contact_id").
		Where("zone_id = ?", zoneId)

	result := db.
		Where("id IN (?)", subQuery).
		Where("id IN (?)", contactIds).
		Delete(&models.SOSContact{})

	return result.RowsAffected, result.Error
}

func NewSOSContactRepo(db *gorm.DB) SOSContactRepo {
	return &SOSContactRepoImpl{db: db}
}
