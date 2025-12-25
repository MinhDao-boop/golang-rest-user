package handler

import (
	"errors"
	"golang-rest-user/dto"
	"golang-rest-user/response"
	"golang-rest-user/service"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TenantHandler struct {
	svc service.TenantService
}

func NewTenantHandler(s service.TenantService) *TenantHandler {
	return &TenantHandler{svc: s}
}

// GET /tenants?page=1&page_size=10&search=...
func (h *TenantHandler) ListTenantResponse(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	search := c.Query("search")

	tenants, total, err := h.svc.List(page, pageSize, search)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusInternalServerError)
	}

	resp := []dto.TenantResponse{}
	for _, t := range tenants {
		resp = append(resp, dto.TenantResponse{
			ID:        t.ID,
			Code:      t.Code,
			Name:      t.Name,
			DBHost:    t.DBHost,
			DBPort:    t.DBPort,
			DBName:    t.DBName,
			CreatedAt: t.CreatedAt.Format(time.RFC3339),
			UpdatedAt: t.UpdatedAt.Format(time.RFC3339),
		})
	}
	response.Success(c, dto.ListTenantResponse{
		Data: resp, Page: page, PageSize: pageSize, Total: total,
	})
}

// POST /tenants
func (h *TenantHandler) CreateTenantRequest(c *gin.Context) {
	var req dto.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	tenant, err := h.svc.Create(req)
	if err != nil {
		if strings.Contains(err.Error(), "exists") {
			response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusConflict)
			return
		}
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	location := c.Request.URL.Path + "/" + strconv.Itoa(int(tenant.ID))
	c.Header("Location", location)
	response.Success(c, dto.TenantResponse{
		ID:        tenant.ID,
		Code:      tenant.Code,
		Name:      tenant.Name,
		DBHost:    tenant.DBHost,
		DBPort:    tenant.DBPort,
		DBName:    tenant.DBName,
		CreatedAt: tenant.CreatedAt.Format(time.RFC3339),
		UpdatedAt: tenant.UpdatedAt.Format(time.RFC3339),
	})
}

// GET /tenants/:code
func (h *TenantHandler) GetByTenantCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.Error(c, response.CodeBadRequest, "tenant code is required", nil, http.StatusBadRequest)
	}
	tenant, err := h.svc.GetByTenantCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, response.CodeBadRequest, "tenant not found", nil, http.StatusNotFound)
			return
		}
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusInternalServerError)
		return
	}
	response.Success(c, dto.TenantResponse{
		ID:        tenant.ID,
		Code:      tenant.Code,
		Name:      tenant.Name,
		DBHost:    tenant.DBHost,
		DBPort:    tenant.DBPort,
		DBName:    tenant.DBName,
		CreatedAt: tenant.CreatedAt.Format(time.RFC3339),
		UpdatedAt: tenant.UpdatedAt.Format(time.RFC3339),
	})
}

// PUT /tenants/:code
func (h *TenantHandler) UpdateTenant(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.Error(c, response.CodeBadRequest, "tenant code is required", nil, http.StatusBadRequest)
	}
	var req dto.UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	tenant, err := h.svc.Update(code, req)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusNotFound)
		return
	}

	response.Success(c, dto.TenantResponse{
		ID:        tenant.ID,
		Code:      tenant.Code,
		Name:      tenant.Name,
		DBHost:    tenant.DBHost,
		DBPort:    tenant.DBPort,
		DBName:    tenant.DBName,
		CreatedAt: tenant.CreatedAt.Format(time.RFC3339),
		UpdatedAt: tenant.UpdatedAt.Format(time.RFC3339),
	})
}

// DELETE /tenants/:code
func (h *TenantHandler) DeleteTenant(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.Error(c, response.CodeBadRequest, "tenant code is required", nil, http.StatusBadRequest)
	}
	if err := h.svc.Delete(code); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusNotFound)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

// PUT /tenants/deleted/:code
func (h *TenantHandler) RecoverDeleted(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.Error(c, response.CodeBadRequest, "tenant code is required", nil, http.StatusBadRequest)
	}
	tenant, err := h.svc.RecoverDeleted(code)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
	}
	response.Success(c, gin.H{"recovered": dto.TenantResponse{
		ID:        tenant.ID,
		Code:      tenant.Code,
		Name:      tenant.Name,
		DBHost:    tenant.DBHost,
		DBPort:    tenant.DBPort,
		DBName:    tenant.DBName,
		CreatedAt: tenant.CreatedAt.Format(time.RFC3339),
		UpdatedAt: tenant.UpdatedAt.Format(time.RFC3339),
	}})
}
