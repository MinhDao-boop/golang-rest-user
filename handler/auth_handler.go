package handler

import (
	"golang-rest-user/provider/tenantProvider"
	"net/http"

	"golang-rest-user/dto"
	"golang-rest-user/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	tenantCode := c.GetHeader("X-Tenant-Code")
	if tenantCode == "" {
		return
	}
	service := tenantProvider.GetTenantInfo(tenantCode)
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, "invalid request", nil, http.StatusBadRequest)
		return
	}

	userResponse, err := service.AuthService.Register(req)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	response.Success(c, userResponse)
}

// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	tenantCode := c.GetHeader("X-Tenant-Code")
	if tenantCode == "" {
		return
	}
	service := tenantProvider.GetTenantInfo(tenantCode)
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	tokens, err := service.AuthService.Login(tenantCode, req)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusUnauthorized)
		return
	}

	response.Success(c, tokens)
}

// POST /auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	tenantCode := c.GetHeader("X-Tenant-Code")
	if tenantCode == "" {
		return
	}
	service := tenantProvider.GetTenantInfo(tenantCode)
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	tokens, err := service.AuthService.Refresh(req.RefreshToken)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusUnauthorized)
		return
	}
	response.Success(c, tokens)
}

// POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	tenantCode := c.GetHeader("X-Tenant-Code")
	if tenantCode == "" {
		return
	}
	service := tenantProvider.GetTenantInfo(tenantCode)
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	if err := service.AuthService.Logout(req.RefreshToken); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusUnauthorized)
		return
	}
	response.Success(c, gin.H{"message": "logged out"})
}
