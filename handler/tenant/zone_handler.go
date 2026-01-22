package tenant

import (
	"golang-rest-user/dto"
	"golang-rest-user/provider/tenantProvider"
	"golang-rest-user/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// POST /zones
func CreateZone(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	userId := c.GetUint("user_id")
	if tenantCode == "" {
		return
	}
	service := tenantProvider.GetTenantInfo(tenantCode)

	var req = dto.ZoneDTORequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	zoneResponse, err := service.ZoneService.CreateZone(&req, userId)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	response.Success(c, zoneResponse)
}

// GET /zones
func ListZones(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	userId := c.GetUint("user_id")
	if tenantCode == "" {
		return
	}
	isOwner := c.Query("is_owner") == "true"
	isShared := c.Query("is_shared") == "true"
	zoneUUID := c.Query("zone_uuid")
	service := tenantProvider.GetTenantInfo(tenantCode)
	zoneResponse, err := service.ZoneService.GetZonesByUser(userId, isOwner, isShared, zoneUUID)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	response.Success(c, zoneResponse)
}

// PUT /zone/:zone_uuid
func UpdateZone(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	userId := c.GetUint("user_id")
	zoneUUID := c.Param("zone_uuid")
	if tenantCode == "" {
		return
	}
	service := tenantProvider.GetTenantInfo(tenantCode)
	var req = dto.ZoneDTORequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	zoneResponse, err := service.ZoneService.UpdateZone(&req, zoneUUID, userId)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	response.Success(c, zoneResponse)
}

// DELETE /zones/:zone_uuid
func DeleteZone(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	userId := c.GetUint("user_id")
	zoneUUID := c.Param("zone_uuid")
	if tenantCode == "" {
		return
	}
	service := tenantProvider.GetTenantInfo(tenantCode)
	deleted, err := service.ZoneService.DeleteZones(zoneUUID, userId)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	response.Success(c, gin.H{"zone deleted": deleted})
}
