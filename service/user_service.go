package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"golang-rest-user/dto"
	"golang-rest-user/models"
	"golang-rest-user/repository"
	"golang-rest-user/security"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserService interface {
	Create(*gorm.DB, dto.CreateUserRequest) (*dto.UserResponse, error)
	GetByUUID(*gorm.DB, string) (*dto.UserResponse, error)
	List(db *gorm.DB, page, pageSize int, search string) ([]dto.UserResponse, int64, error)
	Update(db *gorm.DB, uuid string, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteMany(*gorm.DB, []string) (int64, error)
}

type userService struct{}

func NewUserService() UserService {
	return &userService{}
}

func ConvertToUserResponse(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		Uuid:      user.Uuid,
		Username:  user.Username,
		FullName:  user.FullName,
		Phone:     user.Phone,
		Position:  user.Position,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *userService) Create(db *gorm.DB, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	userRepo := repository.NewUserRepo(db)
	// check username existing
	if _, err := userRepo.GetByUsername(req.Username); err == nil {
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

	if err := userRepo.Create(user); err != nil {
		return nil, err
	}
	return ConvertToUserResponse(user), nil
}

func (s *userService) GetByUUID(db *gorm.DB, uuid string) (*dto.UserResponse, error) {
	userRepo := repository.NewUserRepo(db)
	user, err := userRepo.GetByUUID(uuid)
	if err != nil {
		return nil, err
	}
	return ConvertToUserResponse(user), nil
}

func (s *userService) List(db *gorm.DB, page, pageSize int, search string) ([]dto.UserResponse, int64, error) {
	userRepo := repository.NewUserRepo(db)
	search = strings.TrimSpace(search)
	users, total, err := userRepo.GetList(page, pageSize, search)
	if err != nil {
		return nil, 0, err
	}
	var result []dto.UserResponse
	for _, u := range users {
		result = append(result, *ConvertToUserResponse(&u))
	}
	return result, total, nil
}

func (s *userService) Update(db *gorm.DB, uuid string, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	userRepo := repository.NewUserRepo(db)
	user, err := userRepo.GetByUUID(uuid)
	if err != nil {
		return nil, err
	}
	user.FullName = req.FullName
	user.Phone = req.Phone
	user.Position = req.Position
	user.UpdatedAt = time.Now().UTC()

	if err := userRepo.Update(user); err != nil {
		return nil, err
	}
	return ConvertToUserResponse(user), nil
}

func (s *userService) DeleteMany(db *gorm.DB, uuids []string) (int64, error) {
	userRepo := repository.NewUserRepo(db)
	ids := []uint{}
	for _, uu := range uuids {
		if uu == "" {
			continue
		}
		user, err := userRepo.GetByUUID(uu)
		if err != nil {
			return 0, err
		}
		ids = append(ids, user.ID)
	}
	return userRepo.DeleteByIDs(ids)
}
