package tenant

import (
	"golang-rest-user/dto"
	"golang-rest-user/provider/tenantProvider"
	"golang-rest-user/response"

	"github.com/gin-gonic/gin"
)

// POST /auth/register
func Register(c *gin.Context) {
	tenantCode := c.GetString("TENANT_CODE")
	if tenantCode == "" {
		return
	}
	service := tenantProvider.GetTenantInfo(tenantCode)
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.VAD0000, err)
		return
	}

	authResponse := service.AuthService.Register(req)
	response.Data(c, authResponse)
}

// POST /auth/login
func Login(c *gin.Context) {
	tenantCode := c.GetString("TENANT_CODE")
	if tenantCode == "" {
		return
	}
	service := tenantProvider.GetTenantInfo(tenantCode)
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.VAD0000, err)
		return
	}
	authResponse := service.AuthService.Login(req)
	response.Data(c, authResponse)
}

// POST /auth/refresh
func Refresh(c *gin.Context) {
	tenantCode := c.GetString("TENANT_CODE")
	if tenantCode == "" {
		return
	}
	service := tenantProvider.GetTenantInfo(tenantCode)
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.VAD0000, err)
		return
	}
	authResponse := service.AuthService.Refresh(req.RefreshToken)
	response.Data(c, authResponse)
}

// POST /auth/logout
func Logout(c *gin.Context) {
	tenantCode := c.GetString("TENANT_CODE")
	if tenantCode == "" {
		return
	}
	service := tenantProvider.GetTenantInfo(tenantCode)
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.VAD0000, err)
		return
	}
	authResponse := service.AuthService.Logout(req.RefreshToken)
	response.Data(c, authResponse)
}
