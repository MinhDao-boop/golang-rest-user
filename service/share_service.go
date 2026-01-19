package service

import (
	"errors"
	"golang-rest-user/dto"
	"golang-rest-user/enums"
	"golang-rest-user/models"
	"golang-rest-user/repository"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ShareService interface {
	ShareZone(shareUserID uint, req dto.ShareDTORequest) (*dto.ShareDTOResponse, error)
	RevokeUser(zoneUUID, userUUID string, userID uint) (int64, error)
	UpdatePermission(zoneUUID, userUUID string, userID uint, req dto.ShareDTORequest) error
	GetSharedUser(zoneUUID string, userID uint) ([]dto.UserResponse, error)
}

type shareServiceImpl struct {
	userZoneRepo repository.UserZoneRepo
	zoneService  ZoneService
	userService  UserService
}

func (s *shareServiceImpl) GetSharedUser(zoneUUID string, userID uint) ([]dto.UserResponse, error) {
	zone, err := s.checkOwnerPermission(zoneUUID, userID)
	if err != nil {
		return nil, err
	}
	var userResponse []dto.UserResponse
	userZones, err := s.userZoneRepo.GetSharedUser(zone.ID)
	if err != nil {
		return nil, err
	}
	for _, uz := range userZones {
		user, _ := s.userService.GetByID(uz.UserID)
		userResponse = append(userResponse, *user)
	}
	return userResponse, nil
}

func (s *shareServiceImpl) UpdatePermission(zoneUUID, userUUID string, userID uint, req dto.ShareDTORequest) error {
	zone, err := s.checkOwnerPermission(zoneUUID, userID)
	if err != nil {
		return err
	}
	user, _ := s.userService.GetByUUID(userUUID)
	if !enums.IsValidUserPermission(req.Permission) {
		return errors.New("invalid permission")
	}
	return s.userZoneRepo.UpdatePermission(user.ID, zone.ID, req.Permission)
}

func (s *shareServiceImpl) ShareZone(shareUserID uint, req dto.ShareDTORequest) (*dto.ShareDTOResponse, error) {
	if !enums.IsValidUserPermission(req.Permission) {
		return nil, errors.New("invalid permission")
	}
	zone, err := s.checkOwnerPermission(req.ZoneUUID, shareUserID)
	if err != nil {
		return nil, err
	}
	sharedUser, err := s.userService.GetByUUID(req.UserUUID)
	if err != nil {
		return nil, err
	}
	if shareUserID == sharedUser.ID {
		return nil, errors.New("invalid sharing")
	}
	userZone := models.UserZone{
		UserID:     sharedUser.ID,
		ZoneID:     zone.ID,
		Permission: req.Permission,
	}
	userZone.UUID = uuid.New().String()
	userZone.CreatedAt = time.Now()
	if err := s.userZoneRepo.Create(&userZone); err != nil {
		return nil, err
	}
	return s.convertToShareDTOResponse(&userZone), nil
}

func (s *shareServiceImpl) RevokeUser(zoneUUID, userUUID string, shareUserID uint) (int64, error) {
	zone, err := s.checkOwnerPermission(zoneUUID, shareUserID)
	user, _ := s.userService.GetByUUID(userUUID)
	if err != nil {
		return 0, err
	}
	return s.userZoneRepo.Delete(user.ID, zone.ID)
}

func (s *shareServiceImpl) checkOwnerPermission(zoneUUID string, shareUserID uint) (*dto.ZoneDTOResponse, error) {
	zone, err := s.zoneService.GetByUUID(zoneUUID)
	if err != nil {
		return nil, err
	}
	curPermission, err := s.userZoneRepo.GetPermission(shareUserID, zone.Path)
	if err != nil || strings.Compare(curPermission, string(enums.UserOwner)) != 0 {
		return nil, errors.New("permission denied")
	}
	return zone, nil
}

func (s *shareServiceImpl) convertToShareDTOResponse(userZone *models.UserZone) *dto.ShareDTOResponse {
	return &dto.ShareDTOResponse{
		UUID:       userZone.UUID,
		UserID:     userZone.UserID,
		ZoneID:     userZone.ZoneID,
		Permission: userZone.Permission,
		CreatedAt:  userZone.CreatedAt,
		UpdatedAt:  userZone.UpdatedAt,
	}
}

func NewShareService(userZoneRepo repository.UserZoneRepo) ShareService {
	return &shareServiceImpl{userZoneRepo: userZoneRepo}
}
