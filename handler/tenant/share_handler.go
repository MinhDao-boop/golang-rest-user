package tenant

import (
	"golang-rest-user/dto"
	"golang-rest-user/provider/tenantProvider"
	"golang-rest-user/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GET /zones/:uuid/share
func GetSharedUsers(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	userID := c.GetUint("user_id")
	zoneUUID := c.Param("uuid")
	service := tenantProvider.GetTenantInfo(tenantCode)
	userResponse, err := service.ShareService.GetSharedUser(zoneUUID, userID)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusInternalServerError)
		return
	}
	response.Success(c, userResponse)

}

// POST /zones/:uuid/share
func ShareZone(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	userID := c.GetUint("user_id")
	zoneUUID := c.Param("uuid")
	var req = dto.ShareDTORequest{}
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	service := tenantProvider.GetTenantInfo(tenantCode)
	shareResponse, err := service.ShareService.ShareZone(userID, zoneUUID, req)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	response.Success(c, shareResponse)
}

// PUT /zones/:uuid/share/:user_uuid
func UpdatePermission(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	userID := c.GetUint("user_id")
	zoneUUID := c.Param("uuid")
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

// DELETE /zones/:uuid/share/:user_uuid
func RevokeZone(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	userID := c.GetUint("user_id")
	zoneUUID := c.Param("uuid")
	userUUID := c.Param("user_uuid")
	service := tenantProvider.GetTenantInfo(tenantCode)
	total, err := service.ShareService.RevokeUser(zoneUUID, userUUID, userID)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	response.Success(c, gin.H{"deleted": total})
}
