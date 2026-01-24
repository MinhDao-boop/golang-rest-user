package repository

import (
	"golang-rest-user/enums"
	"golang-rest-user/models"

	"gorm.io/gorm"
)

type SOSContactRepo interface {
	Create(*models.SOSContact) error
	Update(string, map[string]interface{}) error
	ListPaging(zoneId uint, page, pageSize int, search string) ([]models.SOSContact, int64, error)
	DeleteByUUIDs([]string, uint) (int64, error)
	GetByUUID(string) (*models.SOSContact, error)
	GetByPhone(interface{}) (*models.SOSContact, error)
}

type SOSContactRepoImpl struct {
	db *gorm.DB
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

func (r *SOSContactRepoImpl) ListPaging(zoneId uint, page, pageSize int, search string) (SosEmployees []models.SOSContact, total int64, err error) {
	offset := (page - 1) * pageSize
	query := r.db.Model(&models.SOSContact{})
	query = query.Where("is_active = ?", enums.ContactActive)
	query = query.Where("zone_id = ?", zoneId)
	query = query.Where("name LIKE ? OR phone LIKE ?", "%"+search+"%", "%"+search+"%")
	if err = query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err = query.Order("id ASC").Offset(offset).Limit(pageSize).Find(&SosEmployees).Error; err != nil {
		return nil, 0, err
	}
	return SosEmployees, total, nil
}

func (r *SOSContactRepoImpl) DeleteByUUIDs(uuid []string, zoneId uint) (int64, error) {
	res := r.db.Unscoped().Where("uuid IN (?) AND zone_id = ?", uuid, zoneId).Delete(&models.SOSContact{})
	return res.RowsAffected, res.Error
}

func NewSOSContactRepo(db *gorm.DB) SOSContactRepo {
	return &SOSContactRepoImpl{db: db}
}
