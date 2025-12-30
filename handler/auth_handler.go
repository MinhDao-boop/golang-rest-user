package handler

import (
	"golang-rest-user/utils"
	//"log"
	"net/http"

	"golang-rest-user/dto"
	"golang-rest-user/response"
	"golang-rest-user/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authSvc service.AuthService
}

func NewAuthHandler(authSvc service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	tenantDB, err := utils.GetTenantDB(c)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, "invalid request", nil, http.StatusBadRequest)
		return
	}

	user, err := h.authSvc.Register(tenantDB, req)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	response.Success(c, user)
}

// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	tenantDB, err := utils.GetTenantDB(c)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	tenantCode := c.GetHeader("X-Tenant-Code")

	tokens, err := h.authSvc.Login(tenantDB, tenantCode, req)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusUnauthorized)
		return
	}

	response.Success(c, tokens)
}

// POST /auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	tenantDB, err := utils.GetTenantDB(c)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	tokens, err := h.authSvc.Refresh(tenantDB, req.RefreshToken)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusUnauthorized)
		return
	}
	response.Success(c, tokens)
}

// POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	if err := h.authSvc.Logout(req.RefreshToken); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusUnauthorized)
		return
	}
	response.Success(c, gin.H{"message": "logged out"})
}
