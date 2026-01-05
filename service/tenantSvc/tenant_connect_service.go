package tenantSvc

import (
	"golang-rest-user/models"
	"golang-rest-user/repository"
)

type TenantConnectInterface interface {
	ListAllTenantConnect() ([]models.Tenant, error)
}

type tenantConnectImpl struct {
	repo repository.TenantRepo
}

func NewTenantConnect(repo repository.TenantRepo) TenantConnectInterface {
	return &tenantConnectImpl{repo: repo}
}

func (t *tenantConnectImpl) ListAllTenantConnect() ([]models.Tenant, error) {
	tenants, err := t.repo.ListAll()
	if err != nil {
		return nil, err
	}
	return tenants, nil
}
