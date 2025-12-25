package service

import (
	"errors"
	"fmt"
	"time"

	"golang-rest-user/dto"
	"golang-rest-user/models"
	"golang-rest-user/repository"
	"golang-rest-user/security"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	Create(dto.CreateUserRequest) (*models.User, error)
	GetByUUID(string) (*models.User, error)
	List(page, pageSize int, search string) ([]models.User, int64, error)
	Update(uuid string, req dto.UpdateUserRequest) (*models.User, error)
	Delete(string) error
	DeleteMany([]string) (int64, error)
}

type userService struct {
	repo repository.UserRepo
}

func NewUserService(r repository.UserRepo) UserService {
	return &userService{repo: r}
}

func (s *userService) Create(req dto.CreateUserRequest) (*models.User, error) {
	// check username existing
	if _, err := s.repo.GetByUsername(req.Username); err == nil {
		return nil, fmt.Errorf("username already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// if other error (like DB err), still return it
		return nil, err
		// continue if record not found
	}

	passEncrypted, err := security.Encrypt(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Uuid:      uuid.NewString(),
		Username:  req.Username,
		Password:  passEncrypted,
		FullName:  req.FullName,
		Phone:     req.Phone,
		Position:  req.Position,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) GetByUUID(uuid string) (*models.User, error) {
	user, err := s.repo.GetByUUID(uuid)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) List(page, pageSize int, search string) ([]models.User, int64, error) {
	return s.repo.GetList(page, pageSize, search)
}

func (s *userService) Update(uuid string, req dto.UpdateUserRequest) (*models.User, error) {
	user, err := s.repo.GetByUUID(uuid)
	if err != nil {
		return nil, err
	}
	user.FullName = req.FullName
	user.Phone = req.Phone
	user.Position = req.Position
	user.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) Delete(uuid string) error {
	user, err := s.repo.GetByUUID(uuid)
	if err != nil {
		return err
	}
	return s.repo.DeleteByID(user.ID)
}

func (s *userService) DeleteMany(uuids []string) (int64, error) {
	ids := []uint{}
	for _, uuid := range uuids {
		if uuid == "" {
			continue
		}
		user, err := s.repo.GetByUUID(uuid)
		if err != nil {
			return 0, err
		}
		ids = append(ids, user.ID)
	}
	return s.repo.DeleteByIDs(ids)
}
