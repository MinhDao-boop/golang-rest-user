package service

import (
	"golang-rest-user/dto"
	"golang-rest-user/enums"
	"golang-rest-user/models"
	"golang-rest-user/repository"
	"golang-rest-user/response"
	"golang-rest-user/utils"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SOSContactService interface {
	Create(req dto.CreateSOSContactRequest, userId uint) *response.Response
	ListByZone(req dto.ListSOSContactRequest, userId uint) *response.Response
	Update(req dto.UpdateSOSContactRequest, contactUuid string, userId uint) *response.Response
	ToggleStatus(req dto.ToggleSOSContactRequest, contactUuid string, userId uint) *response.Response
	DeleteMany(req dto.DeleteSOSContactRequest, userId uint) *response.Response
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

func (s *SOSContactServiceImpl) Create(req dto.CreateSOSContactRequest, userId uint) *response.Response {
	newResponse := response.NewResponse()
	normalizedPhone, err := utils.NormalizePhone(req.Phone)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	//if !enums.IsValidRoleKey(req.Role) {
	//	return nil, errors.New("invalid role")
	//}
	zone, err := s.zoneSvc.GetByUUID(req.ZoneUuid)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	if !s.zoneSvc.CheckPermission(zone.ID, userId, enums.UserEditor) {
		newResponse.Err = response.ErrForbidden
		newResponse.MessageCode = response.FBD0000
		return newResponse
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
	if err = s.sosContactRepo.Create(&sosContact); err != nil {
		newResponse.Err = err
		return newResponse
	}
	mapping := &models.ZoneSOS{
		ZoneID:       zone.ID,
		SOSContactID: sosContact.ID,
	}
	mapping.UUID = uuid.New().String()
	mapping.CreatedAt = time.Now()
	if err = s.zoneSOSRepo.Create(mapping); err != nil {
		newResponse.Err = err
		return newResponse
	}
	newResponse.Data = sosContact
	newResponse.MessageCode = response.SUS0000
	return newResponse
}

func (s *SOSContactServiceImpl) ListByZone(req dto.ListSOSContactRequest, userId uint) *response.Response {
	newResponse := response.NewResponse()
	zone, err := s.zoneSvc.GetByUUID(req.ZoneUuid)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	if !s.zoneSvc.CheckPermission(zone.ID, userId, enums.UserViewer) {
		newResponse.Err = response.ErrForbidden
		newResponse.MessageCode = response.FBD0000
		return newResponse
	}
	req.Search = strings.TrimSpace(req.Search)
	contacts, total, err := s.sosContactRepo.ListByZone(zone.ID, req.Page, req.PageSize, req.Search, req.IsAll, req.IsActive)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	var result []dto.SOSContactResponse
	for _, contact := range contacts {
		result = append(result, *s.convertToDTO(&contact))
	}
	if req.IsAll {
		newResponse.Data = &response.ListResponse{
			Page:     req.Page,
			PageSize: int(total),
			Total:    total,
			Contents: result,
		}
		newResponse.MessageCode = response.SUS0000
		return newResponse
	}
	newResponse.Data = &response.ListResponse{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    total,
		Contents: result,
	}
	newResponse.MessageCode = response.SUS0000
	return newResponse
}

func (s *SOSContactServiceImpl) Update(req dto.UpdateSOSContactRequest, contactUuid string, userId uint) *response.Response {
	newResponse := response.NewResponse()
	zone, err := s.zoneSvc.GetByUUID(req.ZoneUuid)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	contact, err := s.sosContactRepo.GetByUUID(contactUuid)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	if !s.zoneSvc.CheckPermission(zone.ID, userId, enums.UserEditor) {
		newResponse.Err = response.ErrForbidden
		newResponse.MessageCode = response.FBD0000
		return newResponse
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
			newResponse.Err = err
			return newResponse
		}
		contact.Phone = normalizedPhone
	}
	contact.Note = req.Note
	if err = s.sosContactRepo.Update(contactUuid, contact); err != nil {
		newResponse.Err = err
		return newResponse
	}
	newResponse.MessageCode = response.SUS0000
	return newResponse
}

func (s *SOSContactServiceImpl) ToggleStatus(req dto.ToggleSOSContactRequest, contactUuid string, userId uint) *response.Response {
	newResponse := response.NewResponse()
	zone, err := s.zoneSvc.GetByUUID(req.ZoneUuid)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	if !s.zoneSvc.CheckPermission(zone.ID, userId, enums.UserEditor) {
		newResponse.Err = response.ErrForbidden
		newResponse.MessageCode = response.FBD0000
		return newResponse
	}
	contact, err := s.sosContactRepo.GetByUUID(contactUuid)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	if req.IsActive != nil {
		isActive := *req.IsActive
		if !enums.IsValidStatus(isActive) {
			newResponse.Err = response.ErrInvalidStatus
			return newResponse
		}
		contact.IsActive = isActive
	}
	if err = s.sosContactRepo.Update(contactUuid, contact); err != nil {
		newResponse.Err = err
		return newResponse
	}
	newResponse.MessageCode = response.SUS0000
	return newResponse
}

func (s *SOSContactServiceImpl) DeleteMany(req dto.DeleteSOSContactRequest, userId uint) *response.Response {
	var deleted int64
	newResponse := response.NewResponse()
	zone, err := s.zoneSvc.GetByUUID(req.ZoneUuid)
	if err != nil {
		newResponse.Err = err
		newResponse.Data = &response.DeleteResponse{
			Deleted: deleted,
		}
		return newResponse
	}
	if !s.zoneSvc.CheckPermission(zone.ID, userId, enums.UserEditor) {
		newResponse.Err = response.ErrForbidden
		newResponse.MessageCode = response.FBD0000
		newResponse.Data = &response.DeleteResponse{
			Deleted: deleted,
		}
		return newResponse
	}
	deleted, err = s.sosContactRepo.DeleteByIds(req.Ids)
	if err != nil {
		newResponse.Err = err
		newResponse.Data = &response.DeleteResponse{
			Deleted: deleted,
		}
		return newResponse
	}
	newResponse.MessageCode = response.SUS0000
	newResponse.Data = &response.DeleteResponse{
		Deleted: deleted,
	}
	return newResponse
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
