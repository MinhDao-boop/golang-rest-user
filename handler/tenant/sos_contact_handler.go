package tenant

import (
	"golang-rest-user/dto"
	"golang-rest-user/provider/tenantProvider"
	"golang-rest-user/response"

	"github.com/gin-gonic/gin"
)

// POST /sos-contacts
func CreateSOSContact(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	userId := c.GetUint("user_id")
	service := tenantProvider.GetTenantInfo(tenantCode)
	var query dto.ZoneIdQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, response.VAD0000, err)
		return
	}
	var req dto.CreateSOSContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.VAD0000, err)
		return
	}
	sosResponse := service.SOSContactService.Create(req, userId, query.ZoneUuid)
	response.Data(c, sosResponse)
}

// GET /sos-contacts
func ListSOSContact(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	userId := c.GetUint("user_id")
	service := tenantProvider.GetTenantInfo(tenantCode)
	var input dto.ListSOSContactRequest
	if err := c.ShouldBind(&input); err != nil {
		response.Error(c, response.VAD0000, err)
		return
	}
	sosResponse := service.SOSContactService.ListByZone(input, userId)
	response.Data(c, sosResponse)
}

// PATCH /sos-contacts/:contact_uuid
func UpdateSOSContact(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	userId := c.GetUint("user_id")
	contactUuid := c.Param("contact_uuid")
	service := tenantProvider.GetTenantInfo(tenantCode)
	var query dto.ZoneIdQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, response.VAD0000, err)
		return
	}
	var req dto.UpdateSOSContactRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, response.VAD0000, err)
		return
	}
	sosResponse := service.SOSContactService.Update(req, contactUuid, query.ZoneUuid, userId)
	response.Data(c, sosResponse)
}

// PATCH /sos-contacts/:contact_uuid/toggle
func ToggleSOSContactStatus(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	userId := c.GetUint("user_id")
	service := tenantProvider.GetTenantInfo(tenantCode)
	contactUuid := c.Param("contact_uuid")
	var query dto.ZoneIdQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, response.VAD0000, err)
		return
	}
	var req dto.ToggleSOSContactRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, response.VAD0000, err)
		return
	}
	sosResponse := service.SOSContactService.ToggleStatus(req, contactUuid, query.ZoneUuid, userId)
	response.Data(c, sosResponse)
}

// DELETE /sos-contacts
func DeleteSOSContact(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	userId := c.GetUint("user_id")
	service := tenantProvider.GetTenantInfo(tenantCode)
	var query dto.ZoneIdQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, response.VAD0000, err)
		return
	}
	var req dto.DeleteSOSContactRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, response.VAD0000, err)
		return
	}
	sosResponse := service.SOSContactService.DeleteMany(req, userId, query.ZoneUuid)
	response.Data(c, sosResponse)
}
