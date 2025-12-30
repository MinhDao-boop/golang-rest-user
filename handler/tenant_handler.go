package handler

import (
	"errors"
	"golang-rest-user/dto"
	"golang-rest-user/response"
	"golang-rest-user/service"
	"golang-rest-user/utils"
	"net/http"
	"strconv"
	"strings"

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
func (h *TenantHandler) ListTenant(c *gin.Context) {
	page, pageSize := utils.GetPageAndPageSize(c)
	search := c.Query("search")

	tenantResponses, total, err := h.svc.List(page, pageSize, search)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusInternalServerError)
	}

	response.Success(c, gin.H{
		"data":      tenantResponses,
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}

// POST /tenants
func (h *TenantHandler) CreateTenant(c *gin.Context) {
	var req dto.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	tenantResponse, err := h.svc.Create(req)
	if err != nil {
		if strings.Contains(err.Error(), "exists") {
			response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusConflict)
			return
		}
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	location := c.Request.URL.Path + "/" + strconv.Itoa(int(tenantResponse.ID))
	c.Header("Location", location)
	response.Success(c, tenantResponse)
}

// GET /tenants/:code
func (h *TenantHandler) GetByTenantCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.Error(c, response.CodeBadRequest, "tenant code is required", nil, http.StatusBadRequest)
	}
	tenantResponse, err := h.svc.GetByTenantCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, response.CodeBadRequest, "tenant not found", nil, http.StatusNotFound)
			return
		}
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusInternalServerError)
		return
	}
	response.Success(c, tenantResponse)
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

	tenantResponse, err := h.svc.Update(code, req)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusNotFound)
		return
	}

	response.Success(c, tenantResponse)
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
