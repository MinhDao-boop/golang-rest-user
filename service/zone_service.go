package service

import (
	"fmt"
	"golang-rest-user/dto"
	"golang-rest-user/enums"
	"golang-rest-user/models"
	"golang-rest-user/repository"
	"time"

	"github.com/google/uuid"
)

type ZoneService interface {
	CreateZone(request *dto.ZoneDTORequest, userID uint) (*dto.ZoneDTOResponse, error)
	UpdateZone(request *dto.ZoneDTORequest, uuid string) (*dto.ZoneDTOResponse, error)
	GetUserZones(userID uint) ([]dto.ZoneDTOResponse, error)
	DeleteZones(uuid string) (int64, error)
	GetSharedZone(userID uint) ([]dto.ZoneDTOResponse, error)
}

type zoneServiceImpl struct {
	zoneRepo     repository.ZoneRepo
	userZoneRepo repository.UserZoneRepo
}

func (s *zoneServiceImpl) DeleteZones(uuid string) (int64, error) {
	zone, err := s.zoneRepo.GetByUUID(uuid)
	if err != nil {
		return 0, err
	}
	return s.zoneRepo.DeleteByPath(zone.Path)
}

func (s *zoneServiceImpl) CreateZone(request *dto.ZoneDTORequest, userID uint) (*dto.ZoneDTOResponse, error) {
	var parentPath string
	var parentLevel int
	//if _, err := s.zoneRepo.GetByName(request.Name); err == nil {
	//	return nil, fmt.Errorf("zone with name %s already exists", request.Name)
	//}
	if request.ParentID != nil {

		parentZone, err := s.zoneRepo.GetByID(*request.ParentID)
		if err != nil {
			return nil, err
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
		return nil, err
	}

	if request.ParentID == nil {
		newZone.Path = fmt.Sprintf("%d/", newZone.ID)
	} else {
		newZone.Path = fmt.Sprintf("%s%d/", parentPath, newZone.ID)
	}
	if err := s.zoneRepo.UpdateZonePath(newZone.ID, newZone.Path); err != nil {
		return nil, err
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
			return nil, err
		}
	}
	return convertToZoneDTOResponse(&newZone), nil
}
func (s *zoneServiceImpl) UpdateZone(request *dto.ZoneDTORequest, uuid string) (*dto.ZoneDTOResponse, error) {
	zone, err := s.zoneRepo.GetByUUID(uuid)
	if err != nil {
		return nil, err
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
	if err := s.zoneRepo.Update(zone); err != nil {
		return nil, err
	}
	return convertToZoneDTOResponse(zone), nil
}
func (s *zoneServiceImpl) GetUserZones(userID uint) ([]dto.ZoneDTOResponse, error) {
	zoneID, err := s.userZoneRepo.GetZoneID(userID)
	zone, _ := s.zoneRepo.GetByID(zoneID)
	subZones, err := s.zoneRepo.GetSubtreeByPath(zone.Path)
	if err != nil {
		return nil, err
	}
	zoneResponses := make([]dto.ZoneDTOResponse, 0)
	for _, z := range subZones {
		zoneResponses = append(zoneResponses, *convertToZoneDTOResponse(&z))
	}
	return zoneResponses, nil
}

func (s *zoneServiceImpl) GetSharedZone(userID uint) ([]dto.ZoneDTOResponse, error) {
	var zoneResponses []dto.ZoneDTOResponse
	userZones, err := s.userZoneRepo.GetSharedZone(userID)
	if err != nil {
		return nil, err
	}
	for _, uz := range userZones {
		zone, _ := s.zoneRepo.GetByID(uz.ZoneID)
		zoneResponses = append(zoneResponses, *convertToZoneDTOResponse(zone))
	}
	return zoneResponses, nil
}

func NewZoneService(zoneRepo repository.ZoneRepo, userZoneRepo repository.UserZoneRepo) ZoneService {
	return &zoneServiceImpl{zoneRepo: zoneRepo, userZoneRepo: userZoneRepo}
}

func convertToZoneDTOResponse(zone *models.Zone) *dto.ZoneDTOResponse {
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
