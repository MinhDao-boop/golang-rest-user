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
}

type shareServiceImpl struct {
	userZoneRepo repository.UserZoneRepo
	zoneRepo     repository.ZoneRepo
	userRepo     repository.UserRepo
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

func (s *shareServiceImpl) ShareZone(userID uint, zoneUUID string, req dto.ShareDTORequest) (*dto.ShareDTOResponse, error) {
	zone, _ := s.zoneRepo.GetByUUID(zoneUUID)
	curPermission, err := s.userZoneRepo.GetPermission(userID, zone.Path)
	if err != nil || strings.Compare(curPermission, string(enums.UserOwner)) != 0 {
		return nil, errors.New("permission denied")
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
	zone, _ := s.zoneRepo.GetByUUID(zoneUUID)
	user, _ := s.userRepo.GetByUUID(userUUID)
	curPermission, err := s.userZoneRepo.GetPermission(userID, zone.Path)
	if err != nil || strings.Compare(curPermission, string(enums.UserOwner)) != 0 {
		return 0, errors.New("permission denied")
	}
	return s.userZoneRepo.Delete(user.ID, zone.ID)
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
