package service

import (
	"golang-rest-user/response"
	"golang-rest-user/utils"
	"strings"
	"time"

	"golang-rest-user/dto"
	"golang-rest-user/models"
	"golang-rest-user/repository"

	"github.com/google/uuid"
)

type UserService interface {
	Create(dto.CreateUserRequest) *response.Response
	GetByUUID(string) (*dto.UserResponse, error)
	List(request dto.ListUsersRequest) *response.Response
	Update(uuid string, req dto.UpdateUserRequest) *response.Response
	DeleteMany(dto.DeleteUserRequest) *response.Response
	GetByID(uint) (*dto.UserResponse, error)
	GetByUsername(string) (*models.User, error)
}

type userServiceImpl struct {
	tenantCode string
	repo       repository.UserRepo
}

func (s *userServiceImpl) GetByUsername(username string) (*models.User, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userServiceImpl) GetByID(userID uint) (*dto.UserResponse, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	return s.convertToUserResponse(user), nil
}

func NewUserService(tenantCode string, r repository.UserRepo) UserService {
	return &userServiceImpl{repo: r, tenantCode: tenantCode}
}

func (s *userServiceImpl) convertToUserResponse(user *models.User) *dto.UserResponse {
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

func (s *userServiceImpl) Create(req dto.CreateUserRequest) *response.Response {
	newResponse := response.NewResponse()
	if _, err := s.repo.GetByUsername(req.Username); err == nil {
		newResponse.Err = response.ErrExistingUsername
		newResponse.MessageCode = response.REG0007
		return newResponse
	}

	passEncrypted, err := utils.AESGCMEncrypt(req.Password)
	if err != nil {
		newResponse.Err = err
		return newResponse
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

	if err = s.repo.Create(user); err != nil {
		newResponse.Err = err
		return newResponse
	}
	newResponse.Data = s.convertToUserResponse(user)
	newResponse.MessageCode = response.SUS0000
	return newResponse
}

func (s *userServiceImpl) GetByUUID(uuid string) (*dto.UserResponse, error) {
	user, err := s.repo.GetByUUID(uuid)
	if err != nil {
		return nil, err
	}
	return s.convertToUserResponse(user), nil
}

func (s *userServiceImpl) List(req dto.ListUsersRequest) *response.Response {
	newResponse := response.NewResponse()
	req.Search = strings.TrimSpace(req.Search)
	users, total, err := s.repo.GetList(req.Page, req.PageSize, req.Search)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	var result []dto.UserResponse
	for _, u := range users {
		result = append(result, *s.convertToUserResponse(&u))
	}
	newResponse.Data = &response.ListResponse{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
		Contents: result,
	}
	newResponse.MessageCode = response.SUS0000
	return newResponse
}

func (s *userServiceImpl) Update(uuid string, req dto.UpdateUserRequest) *response.Response {
	newResponse := response.NewResponse()
	user, err := s.repo.GetByUUID(uuid)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	user.FullName = req.FullName
	user.Phone = req.Phone
	user.Position = req.Position
	user.UpdatedAt = time.Now().UTC()

	if err = s.repo.Update(user); err != nil {
		newResponse.Err = err
		return newResponse
	}
	newResponse.MessageCode = response.SUS0000
	return newResponse
}

func (s *userServiceImpl) DeleteMany(req dto.DeleteUserRequest) *response.Response {
	newResponse := response.NewResponse()
	req.UUIDs = strings.TrimSpace(req.UUIDs)
	uuids := strings.Split(req.UUIDs, ",")
	deleted, err := s.repo.DeleteByUUIDs(uuids)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	newResponse.MessageCode = response.SUS0000
	newResponse.Data = &response.DeleteResponse{
		Deleted: deleted,
	}
	return newResponse
}
