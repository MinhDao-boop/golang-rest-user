package tenant

import (
	"golang-rest-user/dto"
	"golang-rest-user/provider/tenantProvider"
	"golang-rest-user/response"
	"golang-rest-user/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// POST /zones/:zone_uuid/sos-contacts
func CreateSOSContact(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	zoneUuid := c.Param("zone_uuid")
	service := tenantProvider.GetTenantInfo(tenantCode)
	//raw, _ := c.GetRawData()
	//log.Println(string(raw))
	var req dto.CreateSOSContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	contactResponse, err := service.SOSService.Create(req, zoneUuid)
	if err != nil {
		if strings.Contains(err.Error(), "exists") {
			response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusConflict)
			return
		}
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	response.Success(c, contactResponse)
}

// GET /zones/:zone_uuid/sos-contacts
func ListSOSContact(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	service := tenantProvider.GetTenantInfo(tenantCode)
	zoneUuid := c.Param("zone_uuid")
	page, pageSize := utils.GetPageAndPageSize(c)
	search := c.Query("search")
	contactResponse, total, err := service.SOSService.ListPaging(zoneUuid, page, pageSize, search)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusInternalServerError)
		return
	}
	response.Success(c, gin.H{
		"data":     contactResponse,
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}

// PATCH /zones/:zone_uuid/sos-contacts/:contact_uuid
func UpdateSOSContact(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	contactUuid := c.Param("contact_uuid")
	service := tenantProvider.GetTenantInfo(tenantCode)
	var req = dto.UpdateSOSContactRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	contactResponse, err := service.SOSService.Update(req, contactUuid)
	if err != nil {
		if strings.Contains(err.Error(), "exists") {
			response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusConflict)
			return
		}
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	response.Success(c, contactResponse)
}

// PATCH /zones/:zone_uuid/sos-contacts/:contact_uuid/toggle
func ToggleSOSContactStatus(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	service := tenantProvider.GetTenantInfo(tenantCode)
	contactUuid := c.Param("contact_uuid")
	var req = dto.ToggleSOSContactRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	contactResponse, err := service.SOSService.ToggleStatus(req, contactUuid)
	if err != nil {
		if strings.Contains(err.Error(), "invalid") {
			response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusConflict)
			return
		}
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
	}
	response.Success(c, contactResponse)
}

// DELETE /zones/:zone_uuid/sos-contacts/contact_uuid
func DeleteSOSContact(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	service := tenantProvider.GetTenantInfo(tenantCode)
	zoneUuid := c.Param("zone_uuid")
	var req = dto.DeleteSOSContactRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	deleted, err := service.SOSService.DeleteMany(req, zoneUuid)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	response.Success(c, gin.H{"deleted": deleted})
}
