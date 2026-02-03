package tenant

import (
	"golang-rest-user/dto"
	"golang-rest-user/provider/tenantProvider"
	"golang-rest-user/response"

	"github.com/gin-gonic/gin"
)

// POST /zone/:zone_uuid/escape
func UpsertZoneEscapeLink(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	userId := c.GetUint("user_id")
	zoneUuid := c.Param("zone_uuid")
	service := tenantProvider.GetTenantInfo(tenantCode)
	var req dto.ZoneEscapeLinksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.VAD0000, err)
		return
	}
	linkResponse := service.ZoneEscapeLinkService.Upsert(userId, zoneUuid, &req)
	response.Data(c, linkResponse)
}

// GET /zone/:zone_uuid/escape
func GetWebView(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	userId := c.GetUint("user_id")
	zoneUuid := c.Param("zone_uuid")
	service := tenantProvider.GetTenantInfo(tenantCode)
	linkResponse := service.ZoneEscapeLinkService.GetWebView(userId, zoneUuid)
	response.Data(c, linkResponse)
}
