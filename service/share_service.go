package service

import (
	"errors"
	"golang-rest-user/dto"
	"golang-rest-user/enums"
	"golang-rest-user/models"
	"golang-rest-user/repository"
)

type ShareService interface {
	ShareZone(shareUserID uint, zoneUUID string, req dto.ShareDTORequest) error
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

func (s *shareServiceImpl) ShareZone(shareUserID uint, shareZoneUUID string, req dto.ShareDTORequest) error {
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

	shareZone, err := s.zoneSvc.GetByUUID(shareZoneUUID)
	if err != nil {
		return err
	}
	if !s.zoneSvc.CheckOwnership(shareZone.Path, shareUserID) {
		return errors.New("permission denied")
	}
	//var userZones []models.UserZone
	zones, err := s.zoneSvc.GetSubtreeByPath(shareZone.Path)
	zoneMap := make(map[string]models.Zone)
	for _, z := range zones {
		zoneMap[z.UUID] = z
	}

	if len(req.ZoneUUIDs) == 0 {
		userZone := models.UserZone{
			UserID:     sharedUser.ID,
			ZoneID:     shareZone.ID,
			Permission: req.Permission,
		}
		_ = s.userZoneRepo.Create(&userZone)
	}
	for _, zoneUUID := range req.ZoneUUIDs {
		if _, ok := zoneMap[zoneUUID]; ok {
			userZone := models.UserZone{
				UserID:     sharedUser.ID,
				ZoneID:     zoneMap[zoneUUID].ID,
				Permission: req.Permission,
			}
			_ = s.userZoneRepo.Create(&userZone)
		}
	}
	return nil
}

func (s *shareServiceImpl) RevokeUser(shareZoneUUID, sharedUserUUID string, shareUserID uint) (int64, error) {
	shareZone, err := s.zoneSvc.GetByUUID(shareZoneUUID)
	if err != nil {
		return 0, err
	}
	if !s.zoneSvc.CheckOwnership(shareZone.Path, shareUserID) {
		return 0, errors.New("permission denied")
	}
	user, _ := s.userSvc.GetByUUID(sharedUserUUID)
	var zoneIDs []uint
	zones, err := s.zoneSvc.GetSubtreeByPath(shareZone.Path)
	for _, zone := range zones {
		zoneIDs = append(zoneIDs, zone.ID)
	}
	return s.userZoneRepo.Delete(user.ID, zoneIDs)
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
