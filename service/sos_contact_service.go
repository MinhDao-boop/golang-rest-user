package service

import (
	"encoding/json"
	"fmt"
	"golang-rest-user/dto"
	"golang-rest-user/enums"
	"golang-rest-user/models"
	"golang-rest-user/provider/redisProvider"
	"golang-rest-user/repository"
	"golang-rest-user/response"
	"golang-rest-user/utils"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SOSContactService interface {
	Create(req dto.CreateSOSContactRequest, userId uint, zoneUuid string) *response.Response
	ListByZone(req dto.ListSOSContactRequest, userId uint) *response.Response
	Update(req dto.UpdateSOSContactRequest, contactUuid, zoneUuid string, userId uint) *response.Response
	ToggleStatus(req dto.ToggleSOSContactRequest, contactUuid, zoneUuid string, userId uint) *response.Response
	DeleteMany(req dto.DeleteSOSContactRequest, userId uint, zoneUuid string) *response.Response
}

type SOSContactServiceImpl struct {
	sosContactRepo repository.SOSContactRepo
	zoneSOSRepo    repository.ZoneSOSRepo
	zoneSvc        ZoneService
	tenantCode     string
}

func sosContact(tenant, zoneId string, userId uint) string {
	key := fmt.Sprintf(
		"{%s}:user:%d:zone:%s:sos_contact",
		tenant, userId, zoneId)
	return key
}

func (s *SOSContactServiceImpl) Create(req dto.CreateSOSContactRequest, userId uint, zoneUuid string) *response.Response {
	resp := response.NewResponse()
	key := sosContact(s.tenantCode, zoneUuid, userId)
	normalizedPhone, err := utils.NormalizePhone(req.Phone)
	if err != nil {
		resp.Err = err
		return resp
	}
	//if !enums.IsValidRoleKey(req.Role) {
	//	return nil, errors.New("invalid role")
	//}
	zone, err := s.zoneSvc.GetByUUID(zoneUuid)
	if err != nil {
		resp.Err = err
		return resp
	}
	if !s.zoneSvc.CheckPermission(zone.ID, userId, enums.UserEditor) {
		resp.Err = response.ErrForbidden
		resp.MessageCode = response.FBD0000
		return resp
	}
	sosContact := models.SOSContact{
		Name:     req.Name,
		Role:     req.Role,
		Phone:    normalizedPhone,
		Note:     req.Note,
		IsActive: true,
	}
	sosContact.UUID = uuid.New().String()
	sosContact.CreatedAt = time.Now()
	if err = s.sosContactRepo.Create(&sosContact); err != nil {
		resp.Err = err
		return resp
	}
	mapping := &models.ZoneSOS{
		ZoneID:       zone.ID,
		SOSContactID: sosContact.ID,
	}
	mapping.UUID = uuid.New().String()
	mapping.CreatedAt = time.Now()
	if err = s.zoneSOSRepo.Create(mapping); err != nil {
		resp.Err = err
		return resp
	}
	err = redisProvider.DeleteCachePattern(key)
	if err != nil {
		log.Println(err)
	}
	resp.Data = s.convertToDTO(&sosContact)
	resp.MessageCode = response.SUS0000
	return resp
}

func (s *SOSContactServiceImpl) ListByZone(req dto.ListSOSContactRequest, userId uint) *response.Response {
	resp := response.NewResponse()
	zone, err := s.zoneSvc.GetByUUID(req.ZoneUuid)
	if err != nil {
		resp.Err = err
		return resp
	}
	if !s.zoneSvc.CheckPermission(zone.ID, userId, enums.UserViewer) {
		resp.Err = response.ErrForbidden
		resp.MessageCode = response.FBD0000
		return resp
	}
	key := sosContact(s.tenantCode, zone.UUID, userId)
	key = fmt.Sprintf("%s:page=%d&page_size=%d:", key, req.Page, req.PageSize)
	if req.Search != "" {
		key = fmt.Sprintf("%s&search=%s", key, req.Search)
	}
	if req.IsActive != nil {
		status := "inactive"
		if *req.IsActive {
			status = "active"
		}
		key = fmt.Sprintf("%s&status=%s", key, status)
	}
	var result []dto.SOSContactResponse
	value, err := redisProvider.GetCache(key)
	if err != nil {
		log.Println(err)
	}
	err = json.Unmarshal(value, &result)
	if value == nil || err != nil {
		req.Search = strings.TrimSpace(req.Search)
		contacts, _, err := s.sosContactRepo.ListByZone(zone.ID, req.Page, req.PageSize, req.Search, false, req.IsActive)
		if err != nil {
			resp.Err = err
			return resp
		}
		for _, contact := range contacts {
			result = append(result, *s.convertToDTO(&contact))
		}
		if err = redisProvider.SetCache(key, result, 300, true); err != nil {
			log.Println(err)
		}
	}
	total := len(result)
	resp.Data = &response.ListResponse{
		Page:     req.Page,
		PageSize: req.PageSize,
		Total:    int64(total),
		Contents: result,
	}
	resp.MessageCode = response.SUS0000
	return resp
}

func (s *SOSContactServiceImpl) Update(req dto.UpdateSOSContactRequest, contactUuid, zoneUuid string, userId uint) *response.Response {
	resp := response.NewResponse()
	key := sosContact(s.tenantCode, zoneUuid, userId)
	zone, err := s.zoneSvc.GetByUUID(zoneUuid)
	if err != nil {
		resp.Err = err
		return resp
	}
	contact, err := s.sosContactRepo.GetByContactAndZone(contactUuid, zone.ID)
	if err != nil {
		resp.Err = err
		return resp
	}
	if !s.zoneSvc.CheckPermission(zone.ID, userId, enums.UserEditor) {
		resp.Err = response.ErrForbidden
		resp.MessageCode = response.FBD0000
		return resp
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
			resp.Err = err
			return resp
		}
		contact.Phone = normalizedPhone
	}
	contact.Note = req.Note
	if err = s.sosContactRepo.Updates(contactUuid, contact); err != nil {
		resp.Err = err
		return resp
	}
	err = redisProvider.DeleteCachePattern(key)
	if err != nil {
		log.Println(err)
	}
	resp.MessageCode = response.SUS0000
	return resp
}

func (s *SOSContactServiceImpl) ToggleStatus(req dto.ToggleSOSContactRequest, contactUuid, zoneUuid string, userId uint) *response.Response {
	resp := response.NewResponse()
	key := sosContact(s.tenantCode, zoneUuid, userId)
	zone, err := s.zoneSvc.GetByUUID(zoneUuid)
	if err != nil {
		resp.Err = err
		return resp
	}
	if !s.zoneSvc.CheckPermission(zone.ID, userId, enums.UserEditor) {
		resp.Err = response.ErrForbidden
		resp.MessageCode = response.FBD0000
		return resp
	}
	contact, err := s.sosContactRepo.GetByContactAndZone(contactUuid, zone.ID)
	if err != nil {
		resp.Err = err
		return resp
	}
	if req.IsActive != nil {
		contact.IsActive = *req.IsActive
	}
	if err = s.sosContactRepo.ToggleStatus(contactUuid, contact); err != nil {
		resp.Err = err
		return resp
	}
	err = redisProvider.DeleteCachePattern(key)
	if err != nil {
		log.Println(err)
	}
	resp.MessageCode = response.SUS0000
	return resp
}

func (s *SOSContactServiceImpl) DeleteMany(req dto.DeleteSOSContactRequest, userId uint, zoneUuid string) *response.Response {
	var deleted int64
	key := sosContact(s.tenantCode, zoneUuid, userId)
	resp := response.NewResponse()
	zone, err := s.zoneSvc.GetByUUID(zoneUuid)
	if err != nil {
		resp.Err = err
		resp.Data = &response.DeleteResponse{
			Deleted: deleted,
		}
		return resp
	}
	if !s.zoneSvc.CheckPermission(zone.ID, userId, enums.UserEditor) {
		resp.Err = response.ErrForbidden
		resp.MessageCode = response.FBD0000
		resp.Data = &response.DeleteResponse{
			Deleted: deleted,
		}
		return resp
	}
	deleted, err = s.sosContactRepo.DeleteMany(req.Ids, zone.ID)
	if err != nil {
		resp.Err = err
		resp.Data = &response.DeleteResponse{
			Deleted: deleted,
		}
		return resp
	}
	err = redisProvider.DeleteCachePattern(key)
	if err != nil {
		log.Println(err)
	}
	resp.MessageCode = response.SUS0000
	resp.Data = &response.DeleteResponse{
		Deleted: deleted,
	}
	return resp
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
func NewSOSService(sosRepo repository.SOSContactRepo, zoneSOSRepo repository.ZoneSOSRepo, zoneSvc ZoneService, tenantCode string) SOSContactService {
	return &SOSContactServiceImpl{sosContactRepo: sosRepo, zoneSOSRepo: zoneSOSRepo, zoneSvc: zoneSvc, tenantCode: tenantCode}
}
