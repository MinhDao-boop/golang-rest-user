package service

import (
	"fmt"
	"golang-rest-user/utils"
	"strings"
	"time"

	"golang-rest-user/dto"
	"golang-rest-user/models"
	"golang-rest-user/repository"

	"github.com/google/uuid"
)

type UserService interface {
	Create(dto.CreateUserRequest) (*dto.UserResponse, error)
	GetByUUID(string) (*dto.UserResponse, error)
	List(page, pageSize int, search string) ([]dto.UserResponse, int64, error)
	Update(uuid string, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteMany(dto.DeleteUserRequest) (int64, error)
	GetByID(uint) (*dto.UserResponse, error)
}

type userService struct {
	tenantCode string
	repo       repository.UserRepo
}

func (s *userService) GetByID(userID uint) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	return s.convertToUserResponse(user), nil
}

func NewUserService(tenantCode string, r repository.UserRepo) UserService {
	return &userService{repo: r, tenantCode: tenantCode}
}

func (s *userService) convertToUserResponse(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:        user.ID,
		UUID:      user.UUID,
		Username:  user.Username,
		FullName:  user.FullName,
		Phone:     user.Phone,
		Position:  user.Position,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *userService) Create(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// check username existing
	if _, err := s.repo.GetByUsername(req.Username); err == nil {
		return nil, fmt.Errorf("username already exists")
	}

	passEncrypted, err := utils.AESGCMEncrypt(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: req.Username,
		Password: passEncrypted,
		FullName: req.FullName,
		Phone:    req.Phone,
		Position: req.Position,
	}
	user.UUID = uuid.New().String()
	user.CreatedAt = time.Now()

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return s.convertToUserResponse(user), nil
}

func (s *userService) GetByUUID(uuid string) (*dto.UserResponse, error) {
	user, err := s.repo.GetByUUID(uuid)
	if err != nil {
		return nil, err
	}
	return s.convertToUserResponse(user), nil
}

func (s *userService) List(page, pageSize int, search string) ([]dto.UserResponse, int64, error) {
	search = strings.TrimSpace(search)
	users, total, err := s.repo.GetList(page, pageSize, search)
	if err != nil {
		return nil, 0, err
	}
	var result []dto.UserResponse
	for _, u := range users {
		result = append(result, *s.convertToUserResponse(&u))
	}
	return result, total, nil
}

func (s *userService) Update(uuid string, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
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
	return s.convertToUserResponse(user), nil
}

func (s *userService) DeleteMany(req dto.DeleteUserRequest) (int64, error) {
	return s.repo.DeleteByUUIDs(req.UUIDs)
}
