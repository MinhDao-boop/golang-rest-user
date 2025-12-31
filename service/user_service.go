package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"golang-rest-user/dto"
	"golang-rest-user/models"
	"golang-rest-user/repository"
	"golang-rest-user/security"

	"github.com/google/uuid"
)

type UserService interface {
	Create(context.Context, dto.CreateUserRequest) (*dto.UserResponse, error)
	GetByUUID(context.Context, string) (*dto.UserResponse, error)
	List(db context.Context, page, pageSize int, search string) ([]dto.UserResponse, int64, error)
	Update(db context.Context, uuid string, req dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteMany(context.Context, []string) (int64, error)
}

type userService struct {
	repo repository.UserRepo
}

func NewUserService(r repository.UserRepo) UserService {
	return &userService{repo: r}
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

func (s *userService) Create(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	// check username existing
	if _, err := s.repo.GetByUsername(ctx, req.Username); err == nil {
		return nil, fmt.Errorf("username already exists")
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

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return ConvertToUserResponse(user), nil
}

func (s *userService) GetByUUID(ctx context.Context, uuid string) (*dto.UserResponse, error) {
	user, err := s.repo.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return ConvertToUserResponse(user), nil
}

func (s *userService) List(ctx context.Context, page, pageSize int, search string) ([]dto.UserResponse, int64, error) {
	search = strings.TrimSpace(search)
	users, total, err := s.repo.GetList(ctx, page, pageSize, search)
	if err != nil {
		return nil, 0, err
	}
	var result []dto.UserResponse
	for _, u := range users {
		result = append(result, *ConvertToUserResponse(&u))
	}
	return result, total, nil
}

func (s *userService) Update(ctx context.Context, uuid string, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.repo.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	user.FullName = req.FullName
	user.Phone = req.Phone
	user.Position = req.Position
	user.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return ConvertToUserResponse(user), nil
}

func (s *userService) DeleteMany(ctx context.Context, uuids []string) (int64, error) {
	ids := []uint{}
	for _, uu := range uuids {
		if uu == "" {
			continue
		}
		user, err := s.repo.GetByUUID(ctx, uu)
		if err != nil {
			return 0, err
		}
		ids = append(ids, user.ID)
	}
	return s.repo.DeleteByIDs(ctx, ids)
}
