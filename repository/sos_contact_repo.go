package repository

import (
	"golang-rest-user/enums"
	"golang-rest-user/models"

	"gorm.io/gorm"
)

type SOSContactRepo interface {
	Create(*models.SOSContact) error
	Update(string, map[string]interface{}) error
	ListByZone(zoneId uint, page, pageSize int, search string, isAll bool, isActive *bool) ([]models.SOSContact, int64, error)
	DeleteByIds([]uint) (int64, error)
	GetByUUID(string) (*models.SOSContact, error)
	GetByPhone(interface{}) (*models.SOSContact, error)
}

type SOSContactRepoImpl struct {
	db *gorm.DB
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
		if err := query.Count(&total).Error; err != nil {
			return nil, 0, err
		}
		offset := (page - 1) * pageSize
		query = query.Order("sos_contacts.id ASC").Limit(pageSize).Offset(offset)
	}
	if err := query.Find(&sosContacts).Error; err != nil {
		return nil, 0, err
	}
	if isAll {
		total = int64(len(sosContacts))
	}
	return sosContacts, total, nil
}

func (r *SOSContactRepoImpl) ToggleStatus(uuid string, status enums.SOSContactStatus) error {
	return r.db.Model(models.SOSContact{}).Where("uuid = ?", uuid).Update("status", status).Error
}

func (r *SOSContactRepoImpl) GetByPhone(phone interface{}) (*models.SOSContact, error) {
	var contact models.SOSContact
	if err := r.db.Model(&models.SOSContact{}).Where("phone = ?", phone).First(&contact).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

func (r *SOSContactRepoImpl) GetByUUID(uuid string) (*models.SOSContact, error) {
	var contact models.SOSContact
	if err := r.db.Model(&models.SOSContact{}).Where("uuid = ?", uuid).First(&contact).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

func (r *SOSContactRepoImpl) Create(sosEmployee *models.SOSContact) error {
	return r.db.Create(sosEmployee).Error
}

func (r *SOSContactRepoImpl) Update(uuid string, updates map[string]interface{}) error {
	return r.db.Model(&models.SOSContact{}).Where("uuid = ?", uuid).
		Updates(updates).Error
}

func (r *SOSContactRepoImpl) DeleteByIds(ids []uint) (int64, error) {
	res := r.db.Unscoped().Delete(&models.SOSContact{}, ids)
	return res.RowsAffected, res.Error
}

func NewSOSContactRepo(db *gorm.DB) SOSContactRepo {
	return &SOSContactRepoImpl{db: db}
}
