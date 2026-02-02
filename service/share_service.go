package service

import (
	"golang-rest-user/dto"
	"golang-rest-user/enums"
	"golang-rest-user/models"
	"golang-rest-user/repository"
	"golang-rest-user/response"
)

type ShareService interface {
	ShareZone(shareUserID uint, zoneUUID string, req dto.ShareDTORequest) *response.Response
	RevokeUser(zoneUUID, userUUID string, userID uint) *response.Response
	UpdatePermission(zoneUUID, userUUID string, userID uint, req dto.ShareDTORequest) *response.Response
	GetSharedUser(zoneUUID string, userID uint) *response.Response
}

type shareServiceImpl struct {
	userZoneRepo repository.UserZoneRepo
	zoneSvc      ZoneService
	userSvc      UserService
}

func (s *shareServiceImpl) GetSharedUser(zoneUUID string, shareUserID uint) *response.Response {
	newResponse := response.NewResponse()
	zone, err := s.zoneSvc.GetByUUID(zoneUUID)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	if !s.zoneSvc.CheckPermission(zone.ID, shareUserID, enums.UserOwner) {
		newResponse.Err = response.ErrForbidden
		newResponse.MessageCode = response.FBD0000
		return newResponse
	}
	var userResponse []dto.UserResponse
	userZones, err := s.userZoneRepo.GetSharedUser(zone.ID)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	for _, uz := range userZones {
		user, _ := s.userSvc.GetByID(uz.UserID)
		userResponse = append(userResponse, *user)
	}
	newResponse.Data = userResponse
	newResponse.MessageCode = response.SUS0000
	return newResponse
}

func (s *shareServiceImpl) UpdatePermission(zoneUUID, userUUID string, shareUserID uint, req dto.ShareDTORequest) *response.Response {
	newResponse := response.NewResponse()
	zone, err := s.zoneSvc.GetByUUID(zoneUUID)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	if !s.zoneSvc.CheckPermission(zone.ID, shareUserID, enums.UserOwner) {
		newResponse.Err = response.ErrForbidden
		newResponse.MessageCode = response.FBD0000
		return newResponse
	}
	user, _ := s.userSvc.GetByUUID(userUUID)
	if !enums.IsValidUserPermission(req.Permission) {
		newResponse.Err = response.ErrInvalidPermission
		newResponse.MessageCode = response.VAD0000
		return newResponse
	}

	if err = s.userZoneRepo.UpdatePermission(user.ID, zone.ID, req.Permission); err != nil {
		newResponse.Err = err
		return newResponse
	}
	newResponse.MessageCode = response.SUS0000
	return newResponse
}

func (s *shareServiceImpl) ShareZone(shareUserID uint, shareZoneUUID string, req dto.ShareDTORequest) *response.Response {
	newResponse := response.NewResponse()
	if !enums.IsValidUserPermission(req.Permission) {
		newResponse.Err = response.ErrInvalidPermission
		newResponse.MessageCode = response.VAD0000
		return newResponse
	}
	sharedUser, err := s.userSvc.GetByUUID(req.UserUUID)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	if sharedUser.ID == shareUserID {
		newResponse.Err = response.ErrInvalidSharing
		newResponse.MessageCode = response.VAD0000
		return newResponse
	}

	shareZone, err := s.zoneSvc.GetByUUID(shareZoneUUID)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	if !s.zoneSvc.CheckPermission(shareZone.ID, shareUserID, enums.UserEditor) {
		newResponse.Err = response.ErrForbidden
		newResponse.MessageCode = response.FBD0000
		return newResponse
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
	//for _, zoneUUID := range req.ZoneUUIDs {
	//	if _, ok := zoneMap[zoneUUID]; ok {
	//		userZone := models.UserZone{
	//			UserID:     &sharedUser.ID,
	//			ZoneID:     zoneMap[zoneUUID].ID,
	//			Permission: req.Permission,
	//		}
	//		_ = s.userZoneRepo.Create(&userZone)
	//	}
	//}
	return nil
}

func (s *shareServiceImpl) RevokeUser(shareZoneUUID, sharedUserUUID string, shareUserID uint) *response.Response {
	newResponse := response.NewResponse()
	shareZone, err := s.zoneSvc.GetByUUID(shareZoneUUID)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	if !s.zoneSvc.CheckPermission(shareZone.ID, shareUserID, enums.UserOwner) {
		newResponse.Err = response.ErrForbidden
		newResponse.MessageCode = response.FBD0000
		return newResponse
	}
	user, _ := s.userSvc.GetByUUID(sharedUserUUID)
	var zoneIDs []uint
	zones, err := s.zoneSvc.GetSubtreeByPath(shareZone.Path)
	for _, zone := range zones {
		zoneIDs = append(zoneIDs, zone.ID)
	}
	deleted, err := s.userZoneRepo.Delete(user.ID, zoneIDs)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	newResponse.Data = &response.DeleteResponse{
		Deleted: deleted,
	}
	newResponse.MessageCode = response.SUS0000
	return newResponse
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
