package tenant

import (
	"golang-rest-user/dto"
	"golang-rest-user/provider/tenantProvider"
	"golang-rest-user/response"
	"net/http"

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
	userResponse, err := service.ShareService.GetSharedUser(zoneUUID, userID)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusInternalServerError)
		return
	}
	response.Success(c, userResponse)

}

// POST /zones/share
func ShareZone(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	shareUserID := c.GetUint("user_id")
	var req = dto.ShareDTORequest{}
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	service := tenantProvider.GetTenantInfo(tenantCode)
	shareResponse, err := service.ShareService.ShareZone(shareUserID, req)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	response.Success(c, shareResponse)
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
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	if err := service.ShareService.UpdatePermission(zoneUUID, userUUID, userID, req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	response.Success(c, nil)
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
	total, err := service.ShareService.RevokeUser(zoneUUID, userUUID, userID)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	response.Success(c, gin.H{"deleted": total})
}
