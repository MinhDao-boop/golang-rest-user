package handler

import (
	"golang-rest-user/dto"
	"golang-rest-user/provider/serviceProvider"
	"golang-rest-user/response"

	"github.com/gin-gonic/gin"
)

// GET /tenants?page=1&page_size=10&search=...
func ListTenant(c *gin.Context) {
	appService := serviceProvider.GetInstance()
	var req dto.ListTenantRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, response.VAD0000, err)
		return
	}
	tenantResponses := appService.TenantService.List(req)
	response.Data(c, tenantResponses)
}

// POST /tenants
func CreateTenant(c *gin.Context) {
	appService := serviceProvider.GetInstance()
	var req dto.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.VAD0000, err)
		return
	}
	tenantResponse := appService.TenantService.Create(req)
	response.Data(c, tenantResponse)
}

// GET /tenants/:tenant_uuid
//func GetByTenantCode(c *gin.Context) {
//	appService := serviceProvider.GetInstance()
//	tenantUuid := c.Param("tenant_uuid")
//	tenantResponse, err := appService.TenantService.GetByUUID(tenantUuid)
//	if err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			response.Error(c, response.ERR0001, "tenant not found", nil, http.StatusNotFound)
//			return
//		}
//		response.Error(c, response.ERR0001, err.Error(), nil, http.StatusInternalServerError)
//		return
//	}
//	response.Success(c, tenantResponse)
//}

// PUT /tenants/:tenant_uuid
func UpdateTenant(c *gin.Context) {
	appService := serviceProvider.GetInstance()
	tenantUuid := c.Param("tenant_uuid")
	var req dto.UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.VAD0000, err)
		return
	}
	tenantResponse := appService.TenantService.Update(tenantUuid, req)
	response.Data(c, tenantResponse)
}

// DELETE /tenants/:tenant_uuid
func DeleteTenant(c *gin.Context) {
	appService := serviceProvider.GetInstance()
	tenantUuid := c.Param("tenant_uuid")
	tenantResponse := appService.TenantService.Delete(tenantUuid)
	response.Data(c, tenantResponse)
}
