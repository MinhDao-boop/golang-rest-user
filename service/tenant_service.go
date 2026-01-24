package service

import (
	"errors"
	"golang-rest-user/enums"
	"golang-rest-user/models"
	"golang-rest-user/utils"
	"regexp"
	"strings"

	"golang-rest-user/dto"
	"golang-rest-user/repository"
	"time"

	"github.com/google/uuid"
)

type CallBackFunction func(mode enums.HandleTenant, tenantCode string, tenant *models.Tenant)

var dbnameRegex = regexp.MustCompile("^[a-z0-9_]{1,64}$")

type TenantService interface {
	Create(dto.CreateTenantRequest) (*dto.TenantResponse, error)
	GetByUUID(string) (*dto.TenantResponse, error)
	List(page, pageSize int, search string) ([]dto.TenantResponse, int64, error)
	ListAllTenantConnect() ([]models.Tenant, error)
	Update(tenantCode string, req dto.UpdateTenantRequest) (*dto.TenantResponse, error)
	Delete(string) error
	SetCallBackFunction(CallBackFunction)
}

type tenantServiceImpl struct {
	callBackFunction CallBackFunction
	repo             repository.TenantRepo
}

func NewTenantService(r repository.TenantRepo) TenantService {
	return &tenantServiceImpl{repo: r}
}

func (s *tenantServiceImpl) convertToTenantResponse(tenant *models.Tenant) *dto.TenantResponse {
	return &dto.TenantResponse{
		UUID:      tenant.UUID,
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

func isValidDBName(name string) bool {
	return dbnameRegex.MatchString(name)
}

func (s *tenantServiceImpl) Create(req dto.CreateTenantRequest) (*dto.TenantResponse, error) {
	// check tenant code existing
	if _, err := s.repo.GetByUUID(req.Code); err == nil {
		return nil, errors.New("tenant code already exists")
	}
	//check db name existing
	if _, err := s.repo.GetByDBName(req.DBName); err == nil {
		return nil, errors.New("db name already exists")
	}
	//Validate dbname
	if !isValidDBName(req.DBName) {
		return nil, errors.New("invalid db name")
	}
	//AESGCMEncrypt db user
	encryptedUser, err := utils.AESGCMEncrypt(req.DBUser)
	if err != nil {
		return nil, err
	}
	//AESGCMEncrypt db password
	encryptedPass, err := utils.AESGCMEncrypt(req.DBPass)
	if err != nil {
		return nil, err
	}
	tenant := &models.Tenant{
		Code:   req.Code,
		Name:   req.Name,
		DBUser: encryptedUser,
		DBPass: encryptedPass,
		DBHost: req.DBHost,
		DBPort: req.DBPort,
		DBName: req.DBName,
	}
	tenant.UUID = uuid.New().String()
	tenant.CreatedAt = time.Now()
	if s.callBackFunction != nil {
		go func() {
			s.callBackFunction(enums.AddTenantConnect, tenant.Code, tenant)
		}()
	}
	if err := s.repo.Create(tenant); err != nil {
		go func() {
			s.callBackFunction(enums.DropTenantConnect, tenant.Code, tenant)
		}()
		return nil, err
	}
	return s.convertToTenantResponse(tenant), nil
}

func (s *tenantServiceImpl) GetByUUID(uuid string) (*dto.TenantResponse, error) {
	uuid = strings.TrimSpace(strings.ToLower(uuid))
	tenant, err := s.repo.GetByUUID(uuid)
	if err != nil {
		return nil, err
	}
	return s.convertToTenantResponse(tenant), nil
}

func (s *tenantServiceImpl) List(page, pageSize int, search string) ([]dto.TenantResponse, int64, error) {
	search = strings.TrimSpace(search)
	tenants, total, err := s.repo.GetList(page, pageSize, search)
	if err != nil {
		return nil, 0, err
	}
	var result []dto.TenantResponse
	for _, t := range tenants {
		result = append(result, *s.convertToTenantResponse(&t))
	}
	return result, total, nil
}

func (s *tenantServiceImpl) ListAllTenantConnect() ([]models.Tenant, error) {
	tenants, err := s.repo.ListAll()
	if err != nil {
		return nil, err
	}
	return tenants, nil
}

func (s *tenantServiceImpl) Update(uuid string, req dto.UpdateTenantRequest) (*dto.TenantResponse, error) {
	tenant, err := s.repo.GetByUUID(uuid)
	if err != nil {
		return nil, err
	}
	//AESGCMDecrypt old db user
	oldDBUser, err := utils.AESGCMDecrypt(tenant.DBUser)
	if err != nil {
		return nil, err
	}
	//AESGCMDecrypt old db password
	oldDBPass, err := utils.AESGCMDecrypt(tenant.DBPass)
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
		return s.convertToTenantResponse(tenant), nil
	}
	//AESGCMEncrypt db user
	encryptedUser, err := utils.AESGCMEncrypt(req.DBUser)
	if err != nil {
		return nil, err
	}
	//AESGCMEncrypt db password
	encryptedPass, err := utils.AESGCMEncrypt(req.DBPass)
	if err != nil {
		return nil, err
	}
	tenant.Name = req.Name
	tenant.DBUser = encryptedUser
	tenant.DBPass = encryptedPass
	tenant.DBHost = req.DBHost
	tenant.DBPort = req.DBPort
	tenant.UpdatedAt = time.Now().UTC()

	if s.callBackFunction != nil {
		go func() {
			s.callBackFunction(enums.EditTenantConnect, tenant.Code, tenant)
		}()
	}
	if err := s.repo.Update(tenant); err != nil {
		go func() {
			s.callBackFunction(enums.DeleteTenantConnect, tenant.Code, tenant)
		}()
		return nil, err
	}

	return s.convertToTenantResponse(tenant), nil
}

func needReconnect(oldTenant *models.Tenant, req dto.UpdateTenantRequest) bool {
	return oldTenant.DBUser != req.DBUser ||
		oldTenant.DBPass != req.DBPass ||
		oldTenant.DBHost != req.DBHost ||
		oldTenant.DBPort != req.DBPort
}

func (s *tenantServiceImpl) Delete(uuid string) error {
	tenant, err := s.repo.GetByUUID(uuid)
	if err != nil {
		return err
	}
	if s.callBackFunction != nil {
		go func() {
			s.callBackFunction(enums.DeleteTenantConnect, tenant.Code, tenant)
		}()
	}
	return s.repo.DeleteByUUID(tenant.BaseModel.UUID)
}

func (s *tenantServiceImpl) SetCallBackFunction(callBackFunction CallBackFunction) {
	s.callBackFunction = callBackFunction
}
