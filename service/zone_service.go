package service

import (
	"fmt"
	"golang-rest-user/dto"
	"golang-rest-user/enums"
	"golang-rest-user/models"
	"golang-rest-user/repository"
	"golang-rest-user/response"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ZoneService interface {
	CreateZone(request *dto.ZoneDTORequest, userID uint) *response.Response
	UpdateZone(request *dto.ZoneDTORequest, uuid string, userID uint) *response.Response
	GetZonesByUser(userID uint, isOwner, isShared bool, zoneUUID string) *response.Response
	DeleteZones(uuid string, userID uint) *response.Response
	GetByUUID(uuid string) (*dto.ZoneDTOResponse, error)
	GetById(uint) (*dto.ZoneDTOResponse, error)
	GetByUUIDs(uuids []string) ([]models.Zone, error)
	CheckPermission(zoneID, userID uint, requiredPermission enums.UserPermission) bool
	GetSubtreeByPath(path string) ([]models.Zone, error)
}

type zoneServiceImpl struct {
	zoneRepo     repository.ZoneRepo
	userZoneRepo repository.UserZoneRepo
}

func (s *zoneServiceImpl) GetById(zoneId uint) (*dto.ZoneDTOResponse, error) {
	zone, err := s.zoneRepo.GetByID(zoneId)
	if err != nil {
		return nil, err
	}
	return s.convertToZoneDTOResponse(zone), nil
}

func (s *zoneServiceImpl) GetSubtreeByPath(path string) ([]models.Zone, error) {
	var zones []models.Zone
	zones, err := s.zoneRepo.GetSubtreeByPath(path)
	if err != nil {
		return nil, err
	}
	return zones, nil
}

func (s *zoneServiceImpl) CheckPermission(zoneId, userId uint, requiredPermission enums.UserPermission) bool {
	_, err := s.userZoneRepo.GetUserZone(userId, zoneId)
	if err != nil {
		return false
	}
	curPermission, err := s.userZoneRepo.GetPermission(userId, zoneId)
	if err != nil || !enums.HasPermission(enums.UserPermission(*curPermission), requiredPermission) {
		return false
	}
	return true
}

func (s *zoneServiceImpl) GetByUUIDs(uuids []string) ([]models.Zone, error) {
	zones, err := s.zoneRepo.GetByUUIDs(uuids)
	if err != nil {
		return nil, err
	}
	return zones, nil
}

func (s *zoneServiceImpl) GetByUUID(uuid string) (*dto.ZoneDTOResponse, error) {
	zone, err := s.zoneRepo.GetByUUID(uuid)
	if err != nil {
		return nil, err
	}
	return s.convertToZoneDTOResponse(zone), nil
}

func (s *zoneServiceImpl) DeleteZones(uuid string, userID uint) *response.Response {
	newResponse := response.NewResponse()
	zone, err := s.zoneRepo.GetByUUID(uuid)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	if !s.CheckPermission(zone.ID, userID, enums.UserOwner) {
		newResponse.Err = response.ErrForbidden
		newResponse.MessageCode = response.FBD0000
		return newResponse
	}
	deleted, err := s.zoneRepo.DeleteByPath(zone.Path)
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

func (s *zoneServiceImpl) CreateZone(request *dto.ZoneDTORequest, userID uint) *response.Response {
	newResponse := response.NewResponse()
	var parentPath string
	var parentLevel int
	if request.ParentID != nil {

		parentZone, err := s.zoneRepo.GetByID(*request.ParentID)
		if err != nil {
			newResponse.Err = err
			return newResponse
		}
		if !s.CheckPermission(parentZone.ID, userID, enums.UserEditor) {
			newResponse.Err = response.ErrForbidden
			newResponse.MessageCode = response.FBD0000
			return newResponse
		}
		parentPath = parentZone.Path
		parentLevel = parentZone.Level
	}
	newZone := models.Zone{
		Name:     request.Name,
		Type:     request.Type,
		Metadata: request.Metadata,
		ParentID: request.ParentID,
		Level:    parentLevel + 1,
	}
	newZone.UUID = uuid.New().String()
	newZone.CreatedAt = time.Now()
	if err := s.zoneRepo.Create(&newZone); err != nil {
		newResponse.Err = err
		return newResponse
	}

	if request.ParentID == nil {
		newZone.Path = fmt.Sprintf("%d/", newZone.ID)
	} else {
		newZone.Path = fmt.Sprintf("%s%d/", parentPath, newZone.ID)
	}
	if err := s.zoneRepo.UpdateZonePath(newZone.ID, newZone.Path); err != nil {
		newResponse.Err = err
		return newResponse
	}
	if request.ParentID == nil {
		newUserZone := &models.UserZone{
			UserID:     userID,
			ZoneID:     newZone.ID,
			Permission: enums.UserOwner,
		}
		newUserZone.UUID = uuid.New().String()
		newUserZone.CreatedAt = time.Now()
		if err := s.userZoneRepo.Create(newUserZone); err != nil {
			newResponse.Err = err
			return newResponse
		}
	}
	newResponse.Data = s.convertToZoneDTOResponse(&newZone)
	newResponse.MessageCode = response.SUS0000
	return newResponse
}
func (s *zoneServiceImpl) UpdateZone(request *dto.ZoneDTORequest, uuid string, userID uint) *response.Response {
	newResponse := response.NewResponse()
	zone, err := s.zoneRepo.GetByUUID(uuid)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	if !s.CheckPermission(zone.ID, userID, enums.UserEditor) {
		newResponse.Err = response.ErrForbidden
		newResponse.MessageCode = response.FBD0000
		return newResponse
	}
	zone.Name = request.Name
	zone.Type = request.Type
	zone.Metadata = request.Metadata
	if request.ParentID != nil {
		parentZone, _ := s.zoneRepo.GetByID(*request.ParentID)
		zone.ParentID = request.ParentID
		zone.Path = fmt.Sprintf("%s%d/", parentZone.Path, zone.ID)
		zone.Level = parentZone.Level + 1
	}
	if err = s.zoneRepo.Update(zone); err != nil {
		newResponse.Err = err
		return newResponse
	}
	newResponse.Data = s.convertToZoneDTOResponse(zone)
	newResponse.MessageCode = response.SUS0000
	return newResponse
}

func (s *zoneServiceImpl) GetZonesByUser(userID uint, isOwner, isShared bool, zoneUUID string) *response.Response {
	newResponse := response.NewResponse()
	var rootPath string
	if zoneUUID != "" {
		rootZone, err := s.zoneRepo.GetByUUID(zoneUUID)
		if err != nil {
			newResponse.Err = err
			return newResponse
		}
		rootPath = rootZone.Path
	}

	userZones, err := s.userZoneRepo.GetByUserID(userID)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}

	zoneMap := make(map[uint]models.Zone)

	for _, uz := range userZones {
		if isOwner && uz.Permission != enums.UserOwner {
			continue
		}
		if isShared && uz.Permission == enums.UserOwner {
			continue
		}

		zone, err := s.zoneRepo.GetByID(uz.ZoneID)
		if err != nil {
			continue
		}

		// subtree
		subZones, err := s.zoneRepo.GetSubtreeByPath(zone.Path)
		if err != nil {
			newResponse.Err = err
			return newResponse
		}

		for _, z := range subZones {

			// nếu có zoneUUID → filter subtree
			if rootPath != "" && !strings.HasPrefix(z.Path, rootPath) {
				continue
			}

			zoneMap[z.ID] = z
		}
	}

	zones := make([]models.Zone, 0, len(zoneMap))
	for _, z := range zoneMap {
		zones = append(zones, z)
	}

	sort.Slice(zones, func(i, j int) bool {
		return zones[i].Path < zones[j].Path
	})

	res := make([]dto.ZoneDTOResponse, 0, len(zones))
	for _, z := range zones {
		res = append(res, *s.convertToZoneDTOResponse(&z))
	}
	if len(res) == 0 {
		newResponse.MessageCode = response.SUS0000
		return newResponse
	}
	newResponse.Data = res
	newResponse.MessageCode = response.SUS0000
	return newResponse
}

func NewZoneService(zoneRepo repository.ZoneRepo, userZoneRepo repository.UserZoneRepo) ZoneService {
	return &zoneServiceImpl{zoneRepo: zoneRepo, userZoneRepo: userZoneRepo}
}

func (s *zoneServiceImpl) convertToZoneDTOResponse(zone *models.Zone) *dto.ZoneDTOResponse {
	return &dto.ZoneDTOResponse{
		ID:        zone.ID,
		UUID:      zone.UUID,
		Name:      zone.Name,
		Type:      zone.Type,
		Path:      zone.Path,
		Level:     zone.Level,
		Metadata:  zone.Metadata,
		CreatedAt: zone.CreatedAt,
		UpdatedAt: zone.UpdatedAt,
		ParentID:  zone.ParentID,
	}
}
