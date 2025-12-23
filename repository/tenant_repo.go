package repository

import (
	"strings"

	"golang-rest-user/models"

	"gorm.io/gorm"
)

type TenantRepo interface {
	Create(tenant *models.Tenant) error
	GetByID(id uint) (*models.Tenant, error)
	GetList(page, pageSize int, search string) (tenants []models.Tenant, total int64, err error)
	Update(tenant *models.Tenant) error
	DeleteByID(id uint) error
	GetByTenantCode(tenantCode string) (*models.Tenant, error)
}

type tenantRepo struct {
	db *gorm.DB
}

func NewTenantRepo(db *gorm.DB) TenantRepo {
	return &tenantRepo{db: db}
}

func (r *tenantRepo) Create(tenant *models.Tenant) error {
	return r.db.Create(tenant).Error
}

func (r *tenantRepo) GetByID(id uint) (*models.Tenant, error) {
	var t models.Tenant
	if err := r.db.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *tenantRepo) GetByTenantCode(tenantCode string) (*models.Tenant, error) {
	var t models.Tenant
	if err := r.db.Where("code = ?", tenantCode).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *tenantRepo) GetList(page, pageSize int, search string) (tenants []models.Tenant, total int64, err error) {
	offset := (page - 1) * pageSize
	query := r.db.Model(&models.Tenant{})
	if strings.TrimSpace(search) != "" {
		query = query.Where("tenant_name LIKE ?", "%"+search+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset(offset).Limit(pageSize).Find(&tenants).Error; err != nil {
		return nil, 0, err
	}
	return tenants, total, nil
}

func (r *tenantRepo) Update(tenant *models.Tenant) error {
	return r.db.Save(tenant).Error
}

func (r *tenantRepo) DeleteByID(id uint) error {
	return r.db.Delete(&models.Tenant{}, id).Error
}
