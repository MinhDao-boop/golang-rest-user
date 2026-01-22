package service

import (
	"errors"
	"fmt"
	"golang-rest-user/dto"
	"golang-rest-user/enums"
	"golang-rest-user/models"
	"golang-rest-user/repository"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ZoneService interface {
	CreateZone(request *dto.ZoneDTORequest, userID uint) (*dto.ZoneDTOResponse, error)
	UpdateZone(request *dto.ZoneDTORequest, uuid string, userID uint) (*dto.ZoneDTOResponse, error)
	GetZonesByUser(userID uint, isOwner, isShared bool, zoneUUID string) ([]dto.ZoneDTOResponse, error)
	DeleteZones(uuid string, userID uint) (int64, error)
	GetByUUID(uuid string) (*dto.ZoneDTOResponse, error)
	GetByUUIDs(uuids []string) ([]models.Zone, error)
	CheckOwnership(path string, userID uint) bool
}

type zoneServiceImpl struct {
	zoneRepo     repository.ZoneRepo
	userZoneRepo repository.UserZoneRepo
}

func (s *zoneServiceImpl) CheckOwnership(path string, userID uint) bool {
	curPermission, err := s.userZoneRepo.GetPermission(userID, path)
	if err != nil || strings.Compare(curPermission, string(enums.UserOwner)) != 0 {
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

func (s *zoneServiceImpl) DeleteZones(uuid string, userID uint) (int64, error) {
	zone, err := s.zoneRepo.GetByUUID(uuid)
	if err != nil {
		return 0, err
	}
	if !s.CheckOwnership(zone.Path, userID) {
		return 0, errors.New("permission denied")
	}
	return s.zoneRepo.DeleteByPath(zone.Path)
}

func (s *zoneServiceImpl) CreateZone(request *dto.ZoneDTORequest, userID uint) (*dto.ZoneDTOResponse, error) {
	var parentPath string
	var parentLevel int
	if request.ParentID != nil {

		parentZone, err := s.zoneRepo.GetByID(*request.ParentID)
		if err != nil {
			return nil, err
		}
		if !s.CheckOwnership(parentZone.Path, userID) {
			return nil, errors.New("permission denied")
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
	return s.convertToZoneDTOResponse(&newZone), nil
}
func (s *zoneServiceImpl) UpdateZone(request *dto.ZoneDTORequest, uuid string, userID uint) (*dto.ZoneDTOResponse, error) {
	zone, err := s.zoneRepo.GetByUUID(uuid)
	if err != nil {
		return nil, err
	}
	if !s.CheckOwnership(zone.Path, userID) {
		return nil, errors.New("permission denied")
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
	return s.convertToZoneDTOResponse(zone), nil
}

func (s *zoneServiceImpl) GetZonesByUser(userID uint, isOwner, isShared bool, zoneUUID string) ([]dto.ZoneDTOResponse, error) {
	var rootPath string
	if zoneUUID != "" {
		rootZone, err := s.zoneRepo.GetByUUID(zoneUUID)
		if err != nil {
			return nil, err
		}
		rootPath = rootZone.Path
	}

	userZones, err := s.userZoneRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		for _, z := range subZones {

			// nếu có zoneUUID → filter subtree
			if rootPath != "" && !strings.HasPrefix(z.Path, rootPath) {
				continue
			}

			zoneMap[z.ID] = z
		}
	}

	// 3. convert map → slice
	zones := make([]models.Zone, 0, len(zoneMap))
	for _, z := range zoneMap {
		zones = append(zones, z)
	}

	// 4. sort theo level
	sort.SliceStable(zones, func(i, j int) bool {
		return zones[i].Level < zones[j].Level
	})

	// 5. convert DTO
	res := make([]dto.ZoneDTOResponse, 0, len(zones))
	for _, z := range zones {
		res = append(res, *s.convertToZoneDTOResponse(&z))
	}
	if len(res) == 0 {
		return nil, nil
	}
	return res, nil
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
