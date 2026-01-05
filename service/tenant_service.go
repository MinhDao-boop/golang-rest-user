package service

import (
	"errors"
	"golang-rest-user/models"
	"golang-rest-user/provider/tenantProvider"
	"strings"

	"golang-rest-user/dto"
	"golang-rest-user/repository"
	"golang-rest-user/security"
	"time"

	"gorm.io/gorm"
)

type TenantService interface {
	Create(dto.CreateTenantRequest) (*dto.TenantResponse, error)
	GetByTenantCode(string) (*dto.TenantResponse, error)
	List(page, pageSize int, search string) ([]dto.TenantResponse, int64, error)
	Update(tenantCode string, req dto.UpdateTenantRequest) (*dto.TenantResponse, error)
	Delete(string) error
}

type tenantService struct {
	repo repository.TenantRepo
}

func NewTenantService(r repository.TenantRepo) TenantService {
	return &tenantService{repo: r}
}

func convertToTenantResponse(tenant *models.Tenant) *dto.TenantResponse {
	return &dto.TenantResponse{
		ID:        tenant.ID,
		Code:      tenant.Code,
		Name:      tenant.Name,
		DBHost:    tenant.DBHost,
		DBPort:    tenant.DBPort,
		DBName:    tenant.DBName,
		Status:    tenant.Status,
		CreatedAt: tenant.CreatedAt.Format(time.RFC3339),
		UpdatedAt: tenant.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *tenantService) Create(req dto.CreateTenantRequest) (*dto.TenantResponse, error) {
	// check db name existing
	if _, err := s.repo.GetByTenantCode(req.Code); err == nil {
		return nil, errors.New("tenant code already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// if other error (like DB err), still return it
		return nil, err
		// continue if record not found
	}
	//AESGCMEncrypt db user
	encryptedUser, err := security.AESGCMEncrypt(req.DBUser)
	if err != nil {
		return nil, err
	}
	//AESGCMEncrypt db password
	encryptedPass, err := security.AESGCMEncrypt(req.DBPass)
	if err != nil {
		return nil, err
	}
	tenant := &models.Tenant{
		Code:      req.Code,
		Name:      req.Name,
		DBUser:    encryptedUser,
		DBPass:    encryptedPass,
		DBHost:    req.DBHost,
		DBPort:    req.DBPort,
		DBName:    req.DBName,
		CreatedAt: time.Now().UTC(),
	}
	if err := s.repo.Create(tenant); err != nil {
		return nil, err
	}
	tenantProvider.AddInstance(tenant)

	return convertToTenantResponse(tenant), nil
}

func (s *tenantService) GetByTenantCode(tenantCode string) (*dto.TenantResponse, error) {
	tenantCode = strings.TrimSpace(strings.ToLower(tenantCode))
	tenant, err := s.repo.GetByTenantCode(tenantCode)
	if err != nil {
		return nil, err
	}
	return convertToTenantResponse(tenant), nil
}

func (s *tenantService) List(page, pageSize int, search string) ([]dto.TenantResponse, int64, error) {
	search = strings.TrimSpace(search)
	tenants, total, err := s.repo.GetList(page, pageSize, search)
	if err != nil {
		return nil, 0, err
	}
	var result []dto.TenantResponse
	for _, t := range tenants {
		result = append(result, *convertToTenantResponse(&t))
	}
	return result, total, nil
}

func (s *tenantService) Update(tenantCode string, req dto.UpdateTenantRequest) (*dto.TenantResponse, error) {
	tenant, err := s.repo.GetByTenantCode(tenantCode)
	if err != nil {
		return nil, err
	}
	//AESGCMDecrypt old db user
	oldDBUser, err := security.AESGCMDecrypt(tenant.DBUser)
	if err != nil {
		return nil, err
	}
	//AESGCMDecrypt old db password
	oldDBPass, err := security.AESGCMDecrypt(tenant.DBPass)
	if err != nil {
		return nil, err
	}
	oldTenant := &models.Tenant{
		DBUser: oldDBUser,
		DBPass: oldDBPass,
	}
	if !needReconnect(oldTenant, req) {
		// no need to reconnect, just update other fields
		tenant.Name = req.Name
		tenant.UpdatedAt = time.Now().UTC()
		if err := s.repo.Update(tenant); err != nil {
			return nil, err
		}
		return convertToTenantResponse(tenant), nil
	}
	//AESGCMEncrypt db user
	encryptedUser, err := security.AESGCMEncrypt(req.DBUser)
	if err != nil {
		return nil, err
	}
	//AESGCMEncrypt db password
	encryptedPass, err := security.AESGCMEncrypt(req.DBPass)
	if err != nil {
		return nil, err
	}
	tenant.Name = req.Name
	tenant.DBUser = encryptedUser
	tenant.DBPass = encryptedPass
	tenant.DBHost = req.DBHost
	tenant.DBPort = req.DBPort
	tenant.UpdatedAt = time.Now().UTC()

	tenantProvider.EditInstance(tenant)
	if err := s.repo.Update(tenant); err != nil {
		tenantProvider.DeleteInstance(tenant.Code)
		return nil, err
	}

	return convertToTenantResponse(tenant), nil
}

func needReconnect(oldTenant *models.Tenant, req dto.UpdateTenantRequest) bool {
	return oldTenant.DBUser != req.DBUser ||
		oldTenant.DBPass != req.DBPass ||
		oldTenant.DBHost != req.DBHost ||
		oldTenant.DBPort != req.DBPort
}

func (s *tenantService) Delete(tenantCode string) error {
	tenant, err := s.repo.GetByTenantCode(tenantCode)
	if err != nil {
		return err
	}
	tenantProvider.DeleteInstance(tenantCode)
	return s.repo.DeleteByID(tenant.ID)
}
