package service

import (
	"fmt"
	"golang-rest-user/dto"
	"golang-rest-user/enums"
	"golang-rest-user/models"
	"golang-rest-user/repository"
	"strings"
	"time"
)

type ZoneService interface {
	CreateZone(request *dto.ZoneDTORequest, userID uint) (*dto.ZoneDTOResponse, error)
	UpdateZone(request *dto.ZoneDTORequest) (*dto.ZoneDTOResponse, error)
	GetUserZones(userID uint) (*dto.ZoneDTOResponse, error)
}

type zoneServiceImpl struct {
	zoneRepo     repository.ZoneRepo
	userZoneRepo repository.UserZoneRepo
}

func NewZoneService(zoneRepo repository.ZoneRepo, userZoneRepo repository.UserZoneRepo) ZoneService {
	return &zoneServiceImpl{zoneRepo: zoneRepo, userZoneRepo: userZoneRepo}
}

func convertToZoneDTOResponse(zone *models.Zone) *dto.ZoneDTOResponse {
	return &dto.ZoneDTOResponse{
		ID:        zone.ID,
		Name:      zone.Name,
		Type:      zone.Type,
		Path:      zone.Path,
		Level:     zone.Level,
		Metadata:  zone.Metadata,
		CreatedAt: zone.CreatedAt,
		UpdatedAt: zone.UpdatedAt,
		Children:  make([]*dto.ZoneDTOResponse, 0),
	}
}

func (s *zoneServiceImpl) CreateZone(request *dto.ZoneDTORequest, userID uint) (*dto.ZoneDTOResponse, error) {
	var parentPath string
	var parentLevel int
	if _, err := s.zoneRepo.GetByName(request.Name); err == nil {
		return nil, fmt.Errorf("zone with name %s already exists", request.Name)
	}
	if request.ParentID != 0 {

		parentZone, err := s.zoneRepo.GetByID(request.ParentID)
		if err != nil {
			return nil, err
		}
		permission, _ := s.userZoneRepo.GetPermission(userID, parentZone.Path)
		if strings.Compare(permission, "owner") != 0 {
			return nil, fmt.Errorf("user is not an owner")
		}
		parentPath = parentZone.Path
		parentLevel = parentZone.Level
	}
	newZone := models.Zone{
		Name:     request.Name,
		Type:     request.Type,
		Metadata: request.Metadata,
		ParentID: &request.ParentID,
		Level:    parentLevel + 1,
	}
	newZone.CreatedAt = time.Now()
	if err := s.zoneRepo.Create(&newZone); err != nil {
		return nil, err
	}
	if request.ParentID == 0 {
		newZone.Path = fmt.Sprintf("%d/", newZone.ID)
	} else {
		newZone.Path = fmt.Sprintf("%s%d/", parentPath, newZone.ID)
	}
	if err := s.zoneRepo.UpdateZonePath(newZone.ID, newZone.Path); err != nil {
		return nil, err
	}
	if request.ParentID == 0 {
		newUserZone := &models.UserZone{
			UserID:     userID,
			ZoneID:     newZone.ID,
			Permission: enums.PermOwner,
		}
		newUserZone.CreatedAt = time.Now()
		if err := s.userZoneRepo.Create(newUserZone); err != nil {
			return nil, err
		}
	}
	return convertToZoneDTOResponse(&newZone), nil
}
func (s *zoneServiceImpl) UpdateZone(request *dto.ZoneDTORequest) (*dto.ZoneDTOResponse, error) {
	return nil, nil
}
func (s *zoneServiceImpl) GetUserZones(userID uint) (*dto.ZoneDTOResponse, error) {
	zoneID, err := s.userZoneRepo.GetZoneID(userID)
	zone, _ := s.zoneRepo.GetByID(zoneID)
	subZones, err := s.zoneRepo.GetSubtreeByPath(zone.Path)
	if err != nil {
		return nil, err
	}
	zoneResponses := buildTree(subZones, zoneID)
	return zoneResponses, nil
}

func buildTree(zones []models.Zone, zoneID uint) *dto.ZoneDTOResponse {
	nodeMap := make(map[uint]*dto.ZoneDTOResponse)
	//var roots []dto.ZoneDTOResponse

	for _, z := range zones {
		nodeMap[z.ID] = convertToZoneDTOResponse(&z)
	}

	for _, z := range zones {
		node := nodeMap[z.ID]

		if parent, ok := nodeMap[*z.ParentID]; ok {
			parent.Children = append(parent.Children, node)
			continue
		}

		//roots = append(roots, *node)
	}

	return nodeMap[zoneID]
}
