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
	Create(req dto.CreateSOSContactRequest, zoneUuid string, userId uint) (*dto.SOSContactResponse, error)
	ListByZone(page, pageSize int, search string, isAll bool, isActive *bool, zoneUuid string, userId uint) (*dto.ListResponse, error)
	Update(req dto.UpdateSOSContactRequest, contactUuid, zoneUuid string, userId uint) (*dto.SOSContactResponse, error)
	ToggleStatus(req dto.ToggleSOSContactRequest, contactUuid, zoneUuid string, userId uint) (*dto.SOSContactResponse, error)
	DeleteMany(req dto.DeleteSOSContactRequest, zoneUuid string, userId uint) (int64, error)
	GetByUuid(string) (*dto.SOSContactResponse, error)
	convertToDTO(contact *models.SOSContact) *dto.SOSContactResponse
}

type SOSContactServiceImpl struct {
	sosContactRepo repository.SOSContactRepo
	zoneSOSRepo    repository.ZoneSOSRepo
	zoneSvc        ZoneService
}

func (s *SOSContactServiceImpl) GetByUuid(contactUuid string) (*dto.SOSContactResponse, error) {
	contact, err := s.sosContactRepo.GetByUUID(contactUuid)
	if err != nil {
		return nil, err
	}
	return s.convertToDTO(contact), nil
}

func (s *SOSContactServiceImpl) Create(req dto.CreateSOSContactRequest, zoneUuid string, userId uint) (*dto.SOSContactResponse, error) {
	normalizedPhone, err := utils.NormalizePhone(req.Phone)
	if err != nil {
		return nil, err
	}
	//if !enums.IsValidRoleKey(req.Role) {
	//	return nil, errors.New("invalid role")
	//}
	zone, err := s.zoneSvc.GetByUUID(zoneUuid)
	if err != nil {
		return nil, err
	}
	if !s.zoneSvc.CheckPermission(zone.ID, userId, enums.UserEditor) {
		return nil, errors.New("permission denied")
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

func (s *SOSContactServiceImpl) ListByZone(page, pageSize int, search string, isAll bool, isActive *bool, zoneUuid string, userId uint) (*dto.ListResponse, error) {
	response := &dto.ListResponse{}
	zone, err := s.zoneSvc.GetByUUID(zoneUuid)
	if err != nil {
		return nil, err
	}
	if !s.zoneSvc.CheckPermission(zone.ID, userId, enums.UserViewer) {
		return nil, errors.New("permission denied")
	}
	search = strings.TrimSpace(search)
	contacts, total, err := s.sosContactRepo.ListByZone(zone.ID, page, pageSize, search, isAll, isActive)
	if err != nil {
		return nil, err
	}
	var result []dto.SOSContactResponse
	for _, contact := range contacts {
		result = append(result, *s.convertToDTO(&contact))
	}
	if isAll {
		response.Data = result
		response.Page = page
		response.PageSize = int(total)
		response.Total = total
		return response, nil
	}
	response = &dto.ListResponse{
		Data:     result,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}
	return response, nil
}

func (s *SOSContactServiceImpl) Update(req dto.UpdateSOSContactRequest, contactUuid, zoneUuid string, userId uint) (*dto.SOSContactResponse, error) {
	zone, err := s.zoneSvc.GetByUUID(zoneUuid)
	if err != nil {
		return nil, err
	}
	contact, err := s.sosContactRepo.GetByUUID(contactUuid)
	if err != nil {
		return nil, err
	}
	if !s.zoneSvc.CheckPermission(zone.ID, userId, enums.UserEditor) {
		return nil, errors.New("permission denied")
	}
	if req.Name != nil {
		contact.Name = *req.Name
	}
	if req.Role != nil {
		role := *req.Role
		//if !enums.IsValidRoleKey(role) {
		//	return nil, errors.New("invalid role")
		//}
		contact.Role = role
	}
	if req.Phone != nil {
		phone := *req.Phone
		normalizedPhone, err := utils.NormalizePhone(phone)
		if err != nil {
			return nil, err
		}
		contact.Phone = normalizedPhone
	}
	contact.Note = req.Note
	if err = s.sosContactRepo.Update(contactUuid, contact); err != nil {
		return nil, err
	}
	return s.convertToDTO(contact), nil
}

func (s *SOSContactServiceImpl) ToggleStatus(req dto.ToggleSOSContactRequest, contactUuid, zoneUuid string, userId uint) (*dto.SOSContactResponse, error) {
	zone, err := s.zoneSvc.GetByUUID(zoneUuid)
	if err != nil {
		return nil, err
	}
	if !s.zoneSvc.CheckPermission(zone.ID, userId, enums.UserEditor) {
		return nil, errors.New("permission denied")
	}
	contact, err := s.sosContactRepo.GetByUUID(contactUuid)
	if err != nil {
		return nil, err
	}
	if req.IsActive != nil {
		isActive := *req.IsActive
		if !enums.IsValidStatus(isActive) {
			return nil, errors.New("invalid status")
		}
		contact.IsActive = isActive
	}
	if err = s.sosContactRepo.Update(contactUuid, contact); err != nil {
		return nil, err
	}
	return s.convertToDTO(contact), nil
}

func (s *SOSContactServiceImpl) DeleteMany(req dto.DeleteSOSContactRequest, zoneUuid string, userId uint) (int64, error) {
	zone, err := s.zoneSvc.GetByUUID(zoneUuid)
	if err != nil {
		return 0, err
	}
	if !s.zoneSvc.CheckPermission(zone.ID, userId, enums.UserEditor) {
		return 0, errors.New("permission denied")
	}
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
