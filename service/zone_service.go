package service

import (
	"fmt"
	"golang-rest-user/dto"
	"golang-rest-user/enums"
	"golang-rest-user/models"
	"golang-rest-user/repository"
	"time"
)

type ZoneService interface {
	CreateZone(request *dto.ZoneDTORequest, userID uint) (*dto.ZoneDTOResponse, error)
	UpdateZone(request *dto.ZoneDTORequest) (*dto.ZoneDTOResponse, error)
	GetUserZones(userID uint) ([]dto.ZoneDTOResponse, error)
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
	}
}

func (s *zoneServiceImpl) CreateZone(request *dto.ZoneDTORequest, userID uint) (*dto.ZoneDTOResponse, error) {
	var parentPath string
	var parentLevel int
	if request.ParentID != 0 {
		parent, err := s.zoneRepo.GetByID(request.ParentID)
		if err != nil {
			return nil, err
		}
		parentPath = parent.Path
		parentLevel = parent.Level
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
		if err := s.userZoneRepo.AddUserPermission(newUserZone); err != nil {
			return nil, err
		}
	}
	return convertToZoneDTOResponse(&newZone), nil
}
func (s *zoneServiceImpl) UpdateZone(request *dto.ZoneDTORequest) (*dto.ZoneDTOResponse, error) {
	return nil, nil
}
func (s *zoneServiceImpl) GetUserZones(userID uint) ([]dto.ZoneDTOResponse, error) {
	userZones, err := s.userZoneRepo.GetUserZones(userID)
	if err != nil {
		return nil, err
	}
	if len(userZones) == 0 {
		return []dto.ZoneDTOResponse{}, nil
	}
	var paths []string
	for _, u := range userZones {
		z, _ := s.zoneRepo.GetByID(u.ZoneID)
		paths = append(paths, z.Path)
	}

	zones, _ := s.zoneRepo.GetByPaths(paths)
	var zoneResponses []dto.ZoneDTOResponse
	for _, z := range zones {
		zoneRes := convertToZoneDTOResponse(&z)
		zoneResponses = append(zoneResponses, *zoneRes)
	}
	return zoneResponses, nil
}

//func buildTree(nodes []models.Zone) []models.Zone {
//	nodeMap := make(map[uint]*models.Zone)
//	var rootNodes []models.Zone
//
//	for i := range nodes {
//		nodeMap[nodes[i].ID] = &nodes[i]
//	}
//
//	for i := range nodes {
//		node := &nodes[i]
//		if node.ParentID != nil && nodeMap[*node.ParentID] != nil {
//			parent := nodeMap[*node.ParentID]
//			parent.Children = append(parent.Children, *node)
//		} else {
//			rootNodes = append(rootNodes, *node)
//		}
//	}
//	return rootNodes
//}
