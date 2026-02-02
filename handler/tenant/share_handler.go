package tenant

import (
	"golang-rest-user/dto"
	"golang-rest-user/provider/tenantProvider"
	"golang-rest-user/response"

	"github.com/gin-gonic/gin"
)

// GET /zones/share/:zone_uuid
func GetSharedUsers(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	userID := c.GetUint("user_id")
	zoneUUID := c.Param("zone_uuid")
	service := tenantProvider.GetTenantInfo(tenantCode)
	userResponse := service.ShareService.GetSharedUser(zoneUUID, userID)
	response.Data(c, userResponse)
}

// POST /zones/share/:zone_uuid
func ShareZone(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	shareUserID := c.GetUint("user_id")
	shareZoneUUID := c.Param("zone_uuid")
	var req = dto.ShareDTORequest{}
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, response.VAD0000, err)
		return
	}
	service := tenantProvider.GetTenantInfo(tenantCode)
	shareResponse := service.ShareService.ShareZone(shareUserID, shareZoneUUID, req)
	response.Data(c, shareResponse)
}

// PUT /zones/share/:zone_uuid/:user_uuid
func UpdatePermission(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	userID := c.GetUint("user_id")
	zoneUUID := c.Param("zone_uuid")
	userUUID := c.Param("user_uuid")
	service := tenantProvider.GetTenantInfo(tenantCode)
	var req = dto.ShareDTORequest{}
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, response.VAD0000, err)
		return
	}
	shareResponse := service.ShareService.UpdatePermission(zoneUUID, userUUID, userID, req)
	response.Data(c, shareResponse)
}

// DELETE /zones/share/:zone_uuid/:user_uuid
func RevokeZone(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	userID := c.GetUint("user_id")
	zoneUUID := c.Param("zone_uuid")
	userUUID := c.Param("user_uuid")
	service := tenantProvider.GetTenantInfo(tenantCode)
	shareResponse := service.ShareService.RevokeUser(zoneUUID, userUUID, userID)
	response.Data(c, shareResponse)
}
