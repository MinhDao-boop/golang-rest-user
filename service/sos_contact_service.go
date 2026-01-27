package service

import (
	"errors"
	"golang-rest-user/dto"
	"golang-rest-user/enums"
	"golang-rest-user/models"
	"golang-rest-user/repository"
	"golang-rest-user/utils"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SOSContactService interface {
	Create(req dto.CreateSOSContactRequest, zoneUuid string) (*dto.SOSContactResponse, error)
	ListByZone(zoneUuid string, page, pageSize int, search string, isAll bool, isActive *bool) ([]dto.SOSContactResponse, int64, error)
	Update(req dto.UpdateSOSContactRequest, uuid string) (*dto.SOSContactResponse, error)
	ToggleStatus(req dto.ToggleSOSContactRequest, uuid string) (*dto.SOSContactResponse, error)
	DeleteMany(req dto.DeleteSOSContactRequest) (int64, error)
	convertToDTO(contact *models.SOSContact) *dto.SOSContactResponse
}

type SOSContactServiceImpl struct {
	sosContactRepo repository.SOSContactRepo
	zoneSOSRepo    repository.ZoneSOSRepo
	zoneSvc        ZoneService
}

func (s *SOSContactServiceImpl) Create(req dto.CreateSOSContactRequest, zoneUuid string) (*dto.SOSContactResponse, error) {
	normalizedPhone, err := utils.NormalizePhone(req.Phone)
	if err != nil {
		return nil, err
	}
	if !enums.IsValidRoleKey(req.Role) {
		return nil, errors.New("invalid role")
	}
	zone, err := s.zoneSvc.GetByUUID(zoneUuid)
	if err != nil {
		return nil, err
	}
	sosContact := models.SOSContact{
		Name:     req.Name,
		Role:     req.Role,
		Phone:    normalizedPhone,
		Note:     req.Note,
		IsActive: enums.ContactActive,
	}
	sosContact.UUID = uuid.New().String()
	sosContact.CreatedAt = time.Now()
	if err := s.sosContactRepo.Create(&sosContact); err != nil {
		return nil, err
	}
	mapping := &models.ZoneSOS{
		ZoneID:       zone.ID,
		SOSContactID: sosContact.ID,
	}
	mapping.UUID = uuid.New().String()
	mapping.CreatedAt = time.Now()
	if err := s.zoneSOSRepo.Create(mapping); err != nil {
		return nil, err
	}
	return s.convertToDTO(&sosContact), nil
}

func (s *SOSContactServiceImpl) ListByZone(zoneUuid string, page, pageSize int, search string, isAll bool, isActive *bool) ([]dto.SOSContactResponse, int64, error) {
	zone, err := s.zoneSvc.GetByUUID(zoneUuid)
	if err != nil {
		return nil, 0, err
	}
	search = strings.TrimSpace(search)
	contacts, total, err := s.sosContactRepo.ListByZone(zone.ID, page, pageSize, search, isAll, isActive)
	if err != nil {
		return nil, 0, err
	}
	var result []dto.SOSContactResponse
	for _, contact := range contacts {
		result = append(result, *s.convertToDTO(&contact))
	}
	return result, total, nil
}

func (s *SOSContactServiceImpl) Update(req dto.UpdateSOSContactRequest, uuid string) (*dto.SOSContactResponse, error) {
	contact, err := s.sosContactRepo.GetByUUID(uuid)
	if err != nil {
		return nil, err
	}
	updates := utils.BuildPatchMap(req)
	if len(updates) == 0 {
		return s.convertToDTO(contact), nil
	}
	//if req.Name != nil {
	//	contact.Name = *req.Name
	//}
	if value, ok := updates["name"]; ok {
		contact.Name = value.(string)
	}

	if value, ok := updates["role"]; ok {
		role := value.(enums.SosRoleKey)
		if !enums.IsValidRoleKey(role) {
			return nil, errors.New("invalid role")
		}
		contact.Role = role
	}
	if value, ok := updates["phone"]; ok {
		phoneRaw := value.(string)
		normalizedPhone, err := utils.NormalizePhone(phoneRaw)
		if err != nil {
			return nil, err
		}
		contact.Phone = normalizedPhone
	}
	if err = s.sosContactRepo.Update(uuid, updates); err != nil {
		return nil, err
	}
	return s.convertToDTO(contact), nil
}

func (s *SOSContactServiceImpl) ToggleStatus(req dto.ToggleSOSContactRequest, uuid string) (*dto.SOSContactResponse, error) {
	contact, err := s.sosContactRepo.GetByUUID(uuid)
	if err != nil {
		return nil, err
	}
	updates := utils.BuildPatchMap(req)
	if len(updates) == 0 {
		return s.convertToDTO(contact), nil
	}
	if value, ok := updates["is_active"]; ok {
		isActive := value.(enums.SOSContactStatus)
		if !enums.IsValidStatus(isActive) {
			return nil, errors.New("invalid status")
		}
		contact.IsActive = isActive
	}
	if err = s.sosContactRepo.Update(uuid, updates); err != nil {
		return nil, err
	}
	return s.convertToDTO(contact), nil
}

func (s *SOSContactServiceImpl) DeleteMany(req dto.DeleteSOSContactRequest) (int64, error) {
	return s.sosContactRepo.DeleteByIds(req.Ids)
}

func (s *SOSContactServiceImpl) convertToDTO(contact *models.SOSContact) *dto.SOSContactResponse {
	return &dto.SOSContactResponse{
		ID:        contact.ID,
		UUID:      contact.UUID,
		Name:      contact.Name,
		Role:      contact.Role,
		Phone:     contact.Phone,
		Note:      contact.Note,
		IsActive:  contact.IsActive,
		CreatedAt: contact.CreatedAt,
		UpdatedAt: contact.UpdatedAt,
	}
}
func NewSOSService(sosRepo repository.SOSContactRepo, zoneSOSRepo repository.ZoneSOSRepo, zoneSvc ZoneService) SOSContactService {
	return &SOSContactServiceImpl{sosContactRepo: sosRepo, zoneSOSRepo: zoneSOSRepo, zoneSvc: zoneSvc}
}
