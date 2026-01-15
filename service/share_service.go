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
	ShareZone(userID uint, zoneUUID string, req dto.ShareDTORequest) (*dto.ShareDTOResponse, error)
	RevokeUser(zoneUUID, userUUID string, userID uint) (int64, error)
	UpdatePermission(zoneUUID, userUUID string, userID uint, req dto.ShareDTORequest) error
	GetSharedUser(zoneUUID string, userID uint) ([]dto.UserResponse, error)
}

type shareServiceImpl struct {
	userZoneRepo repository.UserZoneRepo
	zoneRepo     repository.ZoneRepo
	userRepo     repository.UserRepo
}

func (s *shareServiceImpl) GetSharedUser(zoneUUID string, userID uint) ([]dto.UserResponse, error) {
	var userResponse []dto.UserResponse
	zone, err := s.checkOwnerPermission(zoneUUID, userID)
	if err != nil {
		return nil, err
	}
	userZones, err := s.userZoneRepo.GetSharedUser(zone.ID)
	if err != nil {
		return nil, err
	}
	for _, uz := range userZones {
		user, _ := s.userRepo.GetByID(uz.UserID)
		userResponse = append(userResponse, *convertToUserResponse(user))
	}
	return userResponse, nil
}

func (s *shareServiceImpl) UpdatePermission(zoneUUID, userUUID string, userID uint, req dto.ShareDTORequest) error {
	zone, err := s.checkOwnerPermission(zoneUUID, userID)
	if err != nil {
		return err
	}
	user, _ := s.userRepo.GetByUUID(userUUID)
	if !enums.IsValidUserPermission(string(req.Permission)) {
		return errors.New("invalid permission")
	}
	return s.userZoneRepo.UpdatePermission(user.ID, zone.ID, req.Permission)
}

func (s *shareServiceImpl) ShareZone(userID uint, zoneUUID string, req dto.ShareDTORequest) (*dto.ShareDTOResponse, error) {
	zone, err := s.checkOwnerPermission(zoneUUID, userID)
	if err != nil {
		return nil, err
	}
	if userID == req.UserID {
		return nil, errors.New("sharing denied")
	}
	if !enums.IsValidUserPermission(string(req.Permission)) {
		return nil, errors.New("invalid permission")
	}
	userZone := models.UserZone{
		UserID:     req.UserID,
		ZoneID:     zone.ID,
		Permission: req.Permission,
	}
	userZone.UUID = uuid.New().String()
	userZone.CreatedAt = time.Now()
	if err := s.userZoneRepo.Create(&userZone); err != nil {
		return nil, err
	}
	return convertToShareDTOResponse(&userZone), nil
}

func (s *shareServiceImpl) RevokeUser(zoneUUID, userUUID string, userID uint) (int64, error) {
	user, _ := s.userRepo.GetByUUID(userUUID)
	zone, err := s.checkOwnerPermission(zoneUUID, userID)
	if err != nil {
		return 0, err
	}
	return s.userZoneRepo.Delete(user.ID, zone.ID)
}

func (s *shareServiceImpl) checkOwnerPermission(zoneUUID string, userID uint) (*models.Zone, error) {
	zone, _ := s.zoneRepo.GetByUUID(zoneUUID)
	curPermission, err := s.userZoneRepo.GetPermission(userID, zone.Path)
	if err != nil || strings.Compare(curPermission, string(enums.UserOwner)) != 0 {
		return nil, errors.New("permission denied")
	}
	return zone, nil
}

func convertToShareDTOResponse(userZone *models.UserZone) *dto.ShareDTOResponse {
	return &dto.ShareDTOResponse{
		UUID:       userZone.UUID,
		UserID:     userZone.UserID,
		ZoneID:     userZone.ZoneID,
		Permission: userZone.Permission,
		CreatedAt:  userZone.CreatedAt,
		UpdatedAt:  userZone.UpdatedAt,
	}
}

func NewShareService(
	userZoneRepo repository.UserZoneRepo,
	zoneRepo repository.ZoneRepo,
	userRepo repository.UserRepo,
) ShareService {
	return &shareServiceImpl{
		userZoneRepo: userZoneRepo,
		zoneRepo:     zoneRepo,
		userRepo:     userRepo,
	}
}
