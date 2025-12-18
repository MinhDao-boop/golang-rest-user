package service

import (
	"errors"
	//"fmt"
	"golang-rest-user/database"
	"golang-rest-user/dto"
	"golang-rest-user/models"
	"golang-rest-user/repository"
	"golang-rest-user/security"
	"time"

	"gorm.io/gorm"
)

type TenantService interface {
	Create(dto.CreateTenantRequest) (*models.Tenant, error)
	GetByID(uint) (*models.Tenant, error)
	List(page, pageSize int, search string) ([]models.Tenant, int64, error)
	Update(id uint, req dto.UpdateTenantRequest) (*models.Tenant, error)
	Delete(id uint) error
}

type tenantService struct {
	repo repository.TenantRepo
}

func NewTenantService(r repository.TenantRepo) TenantService {
	return &tenantService{repo: r}
}

func (s *tenantService) Create(req dto.CreateTenantRequest) (*models.Tenant, error) {
	// check db name existing
	if _, err := s.repo.GetByTenantCode(req.Code); err == nil {
		return nil, errors.New("tenant code already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// if other error (like DB err), still return it
		return nil, err
		// continue if record not found
	}
	//Encrypt db user
	encryptedUser, err := security.Encrypt(req.DBUser)
	if err != nil {
		return nil, err
	}
	//Encrypt db password
	encryptedPass, err := security.Encrypt(req.DBPass)
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
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	dbCreated := false
	// create tenant database
	if err := database.CreateTenantDatabase(tenant.DBName); err != nil {
		return nil, err
	}
	// connect to tenant database
	tenantDB, err := database.ConnectTenantDB(*tenant)
	if err != nil {
		return nil, err
	}
	// migrate tenant database
	if err := database.Migrate(tenantDB); err != nil {
		return nil, err
	}
	// ping tenant database
	if err := pingDB(tenantDB); err != nil {
		return nil, err
	}
	// add tenant db to map
	database.SetTenantDB(tenant.Code, tenantDB)
	dbCreated = true
	// save tenant record
	if err := s.repo.Create(tenant); err != nil {
		if dbCreated {
			// cleanup tenant database if tenant record creation failed
			database.RemoveTenantDB(tenant.Code)
			database.DropTenantDatabase(tenant.DBName)
		}
		return nil, err
	}
	return tenant, nil
}

func (s *tenantService) GetByID(id uint) (*models.Tenant, error) {
	return s.repo.GetByID(id)
}

func (s *tenantService) List(page, pageSize int, search string) ([]models.Tenant, int64, error) {
	return s.repo.GetList(page, pageSize, search)
}

func (s *tenantService) Update(id uint, req dto.UpdateTenantRequest) (*models.Tenant, error) {
	tenant, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	//Decrypt old db user
	oldDBUser, err := security.Decrypt(tenant.DBUser)
	if err != nil {
		return nil, err
	}
	//Decrypt old db password
	oldDBPass, err := security.Decrypt(tenant.DBPass)
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
		return tenant, nil
	}
	//Encrypt db user
	encryptedUser, err := security.Encrypt(req.DBUser)
	if err != nil {
		return nil, err
	}
	//Encrypt db password
	encryptedPass, err := security.Encrypt(req.DBPass)
	if err != nil {
		return nil, err
	}
	tenant.Name = req.Name
	tenant.DBUser = encryptedUser
	tenant.DBPass = encryptedPass
	tenant.DBHost = req.DBHost
	tenant.DBPort = req.DBPort
	tenant.UpdatedAt = time.Now().UTC()
	dbConnected := false
	// connect to tenant database
	newDB, err := database.ConnectTenantDB(*tenant)
	if err != nil {
		return nil, err
	}
	//ping new db connection
	if err := pingDB(newDB); err != nil {
		return nil, err
	}
	dbConnected = true
	// save tenant record
	if err := s.repo.Update(tenant); err != nil {
		if dbConnected {
			// cleanup new db connection if tenant record update failed
			database.CloseTenantDB(newDB)
		}
		return nil, err
	}
	// swap tenant db in map
	oldDB := database.SwapTenantDB(tenant.Code, newDB)
	// close old db connection
	if err := database.CloseTenantDB(oldDB); err != nil {
		return nil, err
	}

	return tenant, nil
}

func needReconnect(oldTenant *models.Tenant, req dto.UpdateTenantRequest) bool {
	return oldTenant.DBUser != req.DBUser ||
		oldTenant.DBPass != req.DBPass ||
		oldTenant.DBHost != req.DBHost ||
		oldTenant.DBPort != req.DBPort
}

func pingDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
	return sqlDB.Ping()
}

func (s *tenantService) Delete(id uint) error {
	tenant, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	// remove tenant db from map
	oldDB := database.RemoveTenantDB(tenant.Code)
	// close old db connection
	if err := database.CloseTenantDB(oldDB); err != nil {
		return err
	}
	return s.repo.DeleteByID(id)
}
