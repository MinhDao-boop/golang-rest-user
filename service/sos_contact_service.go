package service

import (
	"golang-rest-user/dto"
	"golang-rest-user/enums"
	"golang-rest-user/models"
	"golang-rest-user/provider/redisProvider"
	"golang-rest-user/repository"
	"golang-rest-user/response"
	"golang-rest-user/utils"
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

//	func (s *SOSContactServiceImpl) GetByUuid(contactUuid string) (*dto.SOSContactResponse, error) {
//		contact, err := s.sosContactRepo.GetByContactAndZone(contactUuid)
//		if err != nil {
//			return nil, err
//		}
//		return s.convertToDTO(contact), nil
//	}
func (s *SOSContactServiceImpl) filterContacts(contacts []dto.SOSContactResponse, req dto.ListSOSContactRequest) []dto.SOSContactResponse {
	var result []dto.SOSContactResponse
	search := strings.ToLower(strings.TrimSpace(req.Search))

	for _, c := range contacts {

		if req.IsActive != nil && c.IsActive != *req.IsActive {
			continue
		}

		if search != "" {
			if !strings.Contains(strings.ToLower(c.Name), search) &&
				!strings.Contains(strings.ToLower(c.Phone), search) {
				continue
			}
		}

		result = append(result, c)
	}

	return result
}

func paginate(data []dto.SOSContactResponse, page, pageSize int) []dto.SOSContactResponse {
	if page <= 0 || pageSize <= 0 {
		return data
	}

	start := (page - 1) * pageSize
	if start >= len(data) {
		return []dto.SOSContactResponse{}
	}

	end := start + pageSize
	if end > len(data) {
		end = len(data)
	}

	return data[start:end]
}

func (s *SOSContactServiceImpl) Create(req dto.CreateSOSContactRequest, userId uint, zoneUuid string) *response.Response {
	newResponse := response.NewResponse()
	normalizedPhone, err := utils.NormalizePhone(req.Phone)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	//if !enums.IsValidRoleKey(req.Role) {
	//	return nil, errors.New("invalid role")
	//}
	zone, err := s.zoneSvc.GetByUUID(zoneUuid)
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
		IsActive: true,
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
	newResponse.Data = s.convertToDTO(&sosContact)
	newResponse.MessageCode = response.SUS0000
	_ = redisProvider.RevokeAllContacts(s.tenantCode, zone.UUID)
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
	//contacts, total, err := s.sosContactRepo.ListByZone(zone.ID, req.Page, req.PageSize, req.Search, req.IsAll, req.IsActive)
	//if err != nil {
	//	newResponse.Err = err
	//	return newResponse
	//}
	cachedContacts, err := redisProvider.GetFullContactKeys(s.tenantCode, req.ZoneUuid)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}

	var contacts []dto.SOSContactResponse

	if cachedContacts == nil {

		dbContacts, _, err := s.sosContactRepo.ListByZone(zone.ID, 1, 50, "", true, nil)
		if err != nil {
			newResponse.Err = err
			return newResponse
		}

		for _, c := range dbContacts {
			contacts = append(contacts, *s.convertToDTO(&c))
		}

		_ = redisProvider.SetContactKeys(s.tenantCode, req.ZoneUuid, contacts, 5*time.Minute)

	} else {
		contacts = cachedContacts
	}

	filtered := s.filterContacts(contacts, req)
	total := int64(len(filtered))

	var result []dto.SOSContactResponse
	if !req.IsAll {
		result = paginate(filtered, req.Page, req.PageSize)
	} else {
		result = filtered
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

func (s *SOSContactServiceImpl) Update(req dto.UpdateSOSContactRequest, contactUuid, zoneUuid string, userId uint) *response.Response {
	newResponse := response.NewResponse()
	zone, err := s.zoneSvc.GetByUUID(zoneUuid)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	contact, err := s.sosContactRepo.GetByContactAndZone(contactUuid, zone.ID)
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
	if err = s.sosContactRepo.Updates(contactUuid, contact); err != nil {
		newResponse.Err = err
		return newResponse
	}
	_ = redisProvider.RevokeAllContacts(s.tenantCode, zone.UUID)
	newResponse.MessageCode = response.SUS0000
	return newResponse
}

func (s *SOSContactServiceImpl) ToggleStatus(req dto.ToggleSOSContactRequest, contactUuid, zoneUuid string, userId uint) *response.Response {
	newResponse := response.NewResponse()
	zone, err := s.zoneSvc.GetByUUID(zoneUuid)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	if !s.zoneSvc.CheckPermission(zone.ID, userId, enums.UserEditor) {
		newResponse.Err = response.ErrForbidden
		newResponse.MessageCode = response.FBD0000
		return newResponse
	}
	contact, err := s.sosContactRepo.GetByContactAndZone(contactUuid, zone.ID)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	if req.IsActive != nil {
		contact.IsActive = *req.IsActive
	}
	if err = s.sosContactRepo.ToggleStatus(contactUuid, contact); err != nil {
		newResponse.Err = err
		return newResponse
	}
	_ = redisProvider.RevokeAllContacts(s.tenantCode, zone.UUID)
	newResponse.MessageCode = response.SUS0000
	return newResponse
}

func (s *SOSContactServiceImpl) DeleteMany(req dto.DeleteSOSContactRequest, userId uint, zoneUuid string) *response.Response {
	var deleted int64
	newResponse := response.NewResponse()
	zone, err := s.zoneSvc.GetByUUID(zoneUuid)
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
	deleted, err = s.sosContactRepo.DeleteMany(req.Ids, zone.ID)
	if err != nil {
		newResponse.Err = err
		newResponse.Data = &response.DeleteResponse{
			Deleted: deleted,
		}
		return newResponse
	}
	_ = redisProvider.RevokeAllContacts(s.tenantCode, zone.UUID)
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
func NewSOSService(sosRepo repository.SOSContactRepo, zoneSOSRepo repository.ZoneSOSRepo, zoneSvc ZoneService, tenantCode string) SOSContactService {
	return &SOSContactServiceImpl{sosContactRepo: sosRepo, zoneSOSRepo: zoneSOSRepo, zoneSvc: zoneSvc, tenantCode: tenantCode}
}
