package handler

import (
	//"log"
	"net/http"

	"golang-rest-user/dto"
	"golang-rest-user/repository"
	"golang-rest-user/response"
	"golang-rest-user/security"
	"golang-rest-user/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func getAuthService(c *gin.Context) service.AuthService {
	dbAny, ok := c.Get("TENANT_DB")
	if !ok {
		return nil
	}
	db := dbAny.(*gorm.DB)
	userRepo := repository.NewUserRepo(db)
	refreshTokenRepo := repository.NewRefreshTokenRepo(db)
	jwtConfig := security.LoadJWTConfig()
	jwtManager := security.NewManager(jwtConfig)
	authSvc := service.NewAuthService(userRepo, refreshTokenRepo, jwtManager)
	return authSvc
}

// POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	authSvc := getAuthService(c)
	if authSvc == nil {
		return
	}
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, "invalid request", nil, http.StatusBadRequest)
		return
	}

	user, err := authSvc.Register(req)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	response.Success(c, user)
}

// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	authSvc := getAuthService(c)
	if authSvc == nil {
		return
	}
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	tenantCode := c.GetHeader("X-Tenant-Code")

	tokens, err := authSvc.Login(tenantCode, req)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusUnauthorized)
		return
	}

	response.Success(c, tokens)
}

// POST /auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	authSvc := getAuthService(c)
	if authSvc == nil {
		return
	}
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	tokens, err := authSvc.Refresh(req.RefreshToken)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusUnauthorized)
		return
	}
	response.Success(c, tokens)
}

// POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	authSvc := getAuthService(c)
	if authSvc == nil {
		return
	}
	var req dto.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	if err := authSvc.Logout(req.RefreshToken); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusUnauthorized)
		return
	}
	response.Success(c, gin.H{"message": "logged out"})
}
