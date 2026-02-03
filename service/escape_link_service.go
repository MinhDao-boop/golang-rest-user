package service

import (
	"encoding/json"
	"golang-rest-user/dto"
	"golang-rest-user/enums"
	"golang-rest-user/models"
	"golang-rest-user/repository"
	"golang-rest-user/response"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ZoneEscapeLinkService interface {
	Upsert(userId uint, zoneUuid string, req *dto.ZoneEscapeLinksRequest) *response.Response
	GetWebView(userId uint, zoneUuid string) *response.Response
}

type ZoneEscapeLinkSvcImpl struct {
	repo    repository.ZoneEscapeLinkRepo
	zoneSvc ZoneService
}

func (s *ZoneEscapeLinkSvcImpl) validateWebViewTypeOnly(raw datatypes.JSON) error {
	var obj map[string]interface{}

	if err := json.Unmarshal(raw, &obj); err != nil {
		return response.ErrInvalidJsonObj
	}

	contents, ok := obj["contents"].([]interface{})
	if !ok {
		return response.ErrInvalidJsonArray
	}

	for _, item := range contents {
		_, ok := item.(map[string]interface{})
		if !ok {
			return response.ErrInvalidJsonObj
		}
	}
	return nil
}

func (s *ZoneEscapeLinkSvcImpl) Upsert(userId uint, zoneUuid string, req *dto.ZoneEscapeLinksRequest) *response.Response {
	newResponse := response.NewResponse()
	if err := s.validateWebViewTypeOnly(req.WebView); err != nil {
		newResponse.Err = err
		newResponse.MessageCode = response.VAD0000
		return newResponse
	}
	zone, err := s.zoneSvc.GetByUUID(zoneUuid)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	if !s.zoneSvc.CheckPermission(zone.ID, userId, enums.UserOwner) {
		newResponse.Err = response.ErrForbidden
		newResponse.MessageCode = response.FBD0000
		return newResponse
	}
	link, err := s.repo.GetByZoneID(zone.ID)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	if link == nil {
		newLink := &models.ZoneEscapeLink{
			ZoneID:  zone.ID,
			WebView: req.WebView,
		}
		newLink.UUID = uuid.New().String()
		newLink.CreatedAt = time.Now()
		if err = s.repo.Create(newLink); err != nil {
			newResponse.Err = err
			return newResponse
		}
		newResponse.MessageCode = response.SUS0000
		return newResponse
	}
	if err = s.repo.Update(zone.ID, req.WebView); err != nil {
		newResponse.Err = err
		return newResponse
	}
	newResponse.MessageCode = response.SUS0000
	return newResponse
}

func (s *ZoneEscapeLinkSvcImpl) GetWebView(userId uint, zoneUuid string) *response.Response {
	newResponse := response.NewResponse()
	zone, err := s.zoneSvc.GetByUUID(zoneUuid)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	if !s.zoneSvc.CheckPermission(zone.ID, userId, enums.UserOwner) {
		newResponse.Err = response.ErrForbidden
		newResponse.MessageCode = response.FBD0000
		return newResponse
	}
	webview, err := s.repo.FindByZoneID(zone.ID)
	if err != nil {
		newResponse.Err = err
		return newResponse
	}
	newResponse.Data = s.convertToResponse(webview)
	newResponse.MessageCode = response.SUS0000
	return newResponse
}

func (s *ZoneEscapeLinkSvcImpl) convertToResponse(link *models.ZoneEscapeLink) *dto.ZoneEscapeLinksResponse {
	return &dto.ZoneEscapeLinksResponse{
		Id:        link.ID,
		Uuid:      link.UUID,
		WebView:   link.WebView,
		CreatedAt: link.CreatedAt,
		UpdatedAt: link.UpdatedAt,
	}
}

func NewZoneEscapeLinkService(repo repository.ZoneEscapeLinkRepo, zoneSvc ZoneService) ZoneEscapeLinkService {
	return &ZoneEscapeLinkSvcImpl{repo: repo, zoneSvc: zoneSvc}
}
