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
	ShareZone(shareUserID uint, req dto.ShareDTORequest) error
	RevokeUser(zoneUUID, userUUID string, userID uint) (int64, error)
	UpdatePermission(zoneUUID, userUUID string, userID uint, req dto.ShareDTORequest) error
	GetSharedUser(zoneUUID string, userID uint) ([]dto.UserResponse, error)
}

type shareServiceImpl struct {
	userZoneRepo repository.UserZoneRepo
	zoneSvc      ZoneService
	userSvc      UserService
}

func (s *shareServiceImpl) GetSharedUser(zoneUUID string, shareUserID uint) ([]dto.UserResponse, error) {
	zone, err := s.zoneSvc.GetByUUID(zoneUUID)
	if err != nil {
		return nil, err
	}
	if !s.zoneSvc.CheckOwnership(zone.Path, shareUserID) {
		return nil, errors.New("permission denied")
	}
	var userResponse []dto.UserResponse
	userZones, err := s.userZoneRepo.GetSharedUser(zone.ID)
	if err != nil {
		return nil, err
	}
	for _, uz := range userZones {
		user, _ := s.userSvc.GetByID(uz.UserID)
		userResponse = append(userResponse, *user)
	}
	return userResponse, nil
}

func (s *shareServiceImpl) UpdatePermission(zoneUUID, userUUID string, shareUserID uint, req dto.ShareDTORequest) error {
	zone, err := s.zoneSvc.GetByUUID(zoneUUID)
	if err != nil {
		return err
	}
	if !s.zoneSvc.CheckOwnership(zone.Path, shareUserID) {
		return errors.New("permission denied")
	}
	user, _ := s.userSvc.GetByUUID(userUUID)
	if !enums.IsValidUserPermission(req.Permission) {
		return errors.New("invalid permission")
	}
	return s.userZoneRepo.UpdatePermission(user.ID, zone.ID, req.Permission)
}

func (s *shareServiceImpl) ShareZone(shareUserID uint, req dto.ShareDTORequest) error {
	if !enums.IsValidUserPermission(req.Permission) {
		return errors.New("invalid permission")
	}
	sharedUser, err := s.userSvc.GetByUUID(req.UserUUID)
	if err != nil {
		return err
	}
	if sharedUser.ID == shareUserID {
		return errors.New("invalid sharing")
	}

	zones, err := s.zoneSvc.GetByUUIDs(req.ZoneUUIDs)
	if err != nil {
		return err
	}
	if len(zones) != len(req.ZoneUUIDs) {
		return errors.New("some zones not found")
	}

	for _, z := range zones {
		curPermission, err := s.userZoneRepo.GetPermission(shareUserID, z.Path)
		if err != nil || strings.Compare(curPermission, string(enums.UserOwner)) != 0 {
			return errors.New("permission denied")
		}
	}

	// 4. batch insert
	userZones := make([]models.UserZone, 0, len(zones))
	now := time.Now()

	for _, z := range zones {
		if _, err := s.userZoneRepo.GetByUserIdAndZoneId(sharedUser.ID, z.ID); err == nil {
			continue
		}
		userZone := models.UserZone{
			UserID:     sharedUser.ID,
			ZoneID:     z.ID,
			Permission: req.Permission,
		}
		userZone.UUID = uuid.New().String()
		userZone.CreatedAt = now
		userZones = append(userZones, userZone)
	}
	if len(userZones) == 0 {
		return nil
	}
	return s.userZoneRepo.BatchInsert(userZones)
}

func (s *shareServiceImpl) RevokeUser(zoneUUID, userUUID string, shareUserID uint) (int64, error) {
	zone, err := s.zoneSvc.GetByUUID(zoneUUID)
	if err != nil {
		return 0, err
	}
	if !s.zoneSvc.CheckOwnership(zone.Path, shareUserID) {
		return 0, errors.New("permission denied")
	}
	user, _ := s.userSvc.GetByUUID(userUUID)
	return s.userZoneRepo.Delete(user.ID, zone.ID)
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

func NewShareService(userZoneRepo repository.UserZoneRepo, userSvc UserService, zoneSvc ZoneService) ShareService {
	return &shareServiceImpl{
		userZoneRepo: userZoneRepo,
		zoneSvc:      zoneSvc,
		userSvc:      userSvc,
	}
}
