package service

import (
	"errors"
	"golang-rest-user/dto"
	"golang-rest-user/enums"
	"golang-rest-user/models"
	"golang-rest-user/repository"
	"golang-rest-user/utils"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SOSContactService interface {
	Create(req dto.CreateSOSContactRequest, zoneUuid string) (*dto.SOSContactResponse, error)
	ListPaging(zoneUuid string, page, pageSize int, search string) ([]dto.SOSContactResponse, int64, error)
	Update(req dto.UpdateSOSContactRequest, uuid string) (*dto.SOSContactResponse, error)
	ToggleStatus(req dto.ToggleSOSContactRequest, uuid string) (*dto.SOSContactResponse, error)
	DeleteMany(req dto.DeleteSOSContactRequest, zoneUuid string) (int64, error)
	convertToDTO(contact *models.SOSContact) *dto.SOSContactResponse
}

type SOSContactServiceImpl struct {
	sosContactRepo repository.SOSContactRepo
	zoneSvc        ZoneService
}

var phoneRegex = regexp.MustCompile("^(?:0|(?:\\+[1-9]))\\d{7,14}$")

func (s *SOSContactServiceImpl) Create(req dto.CreateSOSContactRequest, zoneUuid string) (*dto.SOSContactResponse, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	if req.Position == "" {
		return nil, errors.New("position is required")
	}
	if req.Phone == "" {
		return nil, errors.New("phone number is required")
	}
	if !phoneRegex.MatchString(req.Phone) {
		return nil, errors.New("invalid phone number")
	}
	if _, err := s.sosContactRepo.GetByPhone(req.Phone); err == nil {
		return nil, errors.New("phone number already exists")
	}
	zone, err := s.zoneSvc.GetByUUID(zoneUuid)
	if err != nil {
		return nil, err
	}
	sosContact := models.SOSContact{
		Name:     req.Name,
		Position: req.Position,
		Phone:    req.Phone,
		Note:     req.Note,
		IsActive: enums.ContactActive,
		ZoneID:   &zone.ID,
	}
	sosContact.UUID = uuid.New().String()
	sosContact.CreatedAt = time.Now()
	if err := s.sosContactRepo.Create(&sosContact); err != nil {
		return nil, err
	}
	return s.convertToDTO(&sosContact), nil
}

func (s *SOSContactServiceImpl) ListPaging(zoneUuid string, page, pageSize int, search string) ([]dto.SOSContactResponse, int64, error) {
	zone, err := s.zoneSvc.GetByUUID(zoneUuid)
	if err != nil {
		return nil, 0, err
	}
	search = strings.TrimSpace(search)
	contacts, total, err := s.sosContactRepo.ListPaging(zone.ID, page, pageSize, search)
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
	if value, ok := updates["name"]; ok {
		if strings.TrimSpace(value.(string)) == "" {
			return nil, errors.New("name is required")
		}
	}
	if value, ok := updates["position"]; ok {
		if strings.TrimSpace(value.(string)) == "" {
			return nil, errors.New("position is required")
		}
	}
	if value, ok := updates["phone"]; ok {
		phone := strings.TrimSpace(value.(string))
		if phone == "" {
			return nil, errors.New("phone is required")
		}
		if !phoneRegex.MatchString(phone) {
			return nil, errors.New("invalid phone number")
		}
		existed, err := s.sosContactRepo.GetByPhone(phone)
		if err == nil && existed.UUID != uuid {
			return nil, errors.New("phone number already exists")
		}
	}
	if err = s.sosContactRepo.Update(uuid, updates); err != nil {
		return nil, err
	}
	contact, _ = s.sosContactRepo.GetByUUID(uuid)
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
	}
	if err = s.sosContactRepo.Update(uuid, updates); err != nil {
		return nil, err
	}
	contact, _ = s.sosContactRepo.GetByUUID(uuid)
	return s.convertToDTO(contact), nil
}

func (s *SOSContactServiceImpl) DeleteMany(req dto.DeleteSOSContactRequest, zoneUuid string) (int64, error) {
	zone, err := s.zoneSvc.GetByUUID(zoneUuid)
	if err != nil {
		return 0, err
	}
	return s.sosContactRepo.DeleteByUUIDs(req.UUID, zone.ID)
}

func (s *SOSContactServiceImpl) convertToDTO(contact *models.SOSContact) *dto.SOSContactResponse {
	return &dto.SOSContactResponse{
		ID:        contact.ID,
		UUID:      contact.UUID,
		Name:      contact.Name,
		Position:  contact.Position,
		Phone:     contact.Phone,
		Note:      contact.Note,
		IsActive:  contact.IsActive,
		CreatedAt: contact.CreatedAt,
		UpdatedAt: contact.UpdatedAt,
	}
}
func NewSOSService(sosRepo repository.SOSContactRepo, zoneSvc ZoneService) SOSContactService {
	return &SOSContactServiceImpl{sosContactRepo: sosRepo, zoneSvc: zoneSvc}
}
