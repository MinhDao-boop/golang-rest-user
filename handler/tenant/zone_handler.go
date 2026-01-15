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
	service := tenantProvider.GetTenantInfo(tenantCode)
	zoneResponse, err := service.ZoneService.GetUserZones(userId)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	response.Success(c, zoneResponse)
}

// GET /zones/share-with-me
func ListSharedZones(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	userID := c.GetUint("user_id")
	service := tenantProvider.GetTenantInfo(tenantCode)
	zoneResponses, err := service.ZoneService.GetSharedZone(userID)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusInternalServerError)
		return
	}
	response.Success(c, zoneResponses)
}

// PUT /zone/:uuid
func UpdateZone(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	uuid := c.Param("uuid")
	if tenantCode == "" {
		return
	}
	service := tenantProvider.GetTenantInfo(tenantCode)
	var req = dto.ZoneDTORequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	zoneResponse, err := service.ZoneService.UpdateZone(&req, uuid)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	response.Success(c, zoneResponse)
}

// DELETE /zones/:uuid
func DeleteZone(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	uuid := c.Param("uuid")
	if tenantCode == "" {
		return
	}
	service := tenantProvider.GetTenantInfo(tenantCode)
	deleted, err := service.ZoneService.DeleteZones(uuid)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	response.Success(c, gin.H{"zone deleted": deleted})
}
