package handler

import (
	"golang-rest-user/dto"
	"golang-rest-user/response"
	"golang-rest-user/service"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
		//c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
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
		//c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	tenant, err := h.svc.Create(req)
	if err != nil {
		if strings.Contains(err.Error(), "exists") {
			//c.JSON(http.StatusConflict, gin.H{"success": false, "error": err.Error()})
			response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
			return
		}
		//c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
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

// GET /tenants/:id
func (h *TenantHandler) GetByTenantID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	uid := uint(id)
	tenant, err := h.svc.GetByID(uid)
	if err != nil {
		//c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
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

// PUT /tenants/:id
func (h *TenantHandler) UpdateTenant(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	uid := uint(id)
	var req dto.UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	tenant, err := h.svc.Update(uid, req)
	if err != nil {
		//c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
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

// DELETE /tenants/:id
func (h *TenantHandler) DeleteTenant(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	uid := uint(id)
	if err := h.svc.Delete(uid); err != nil {
		//c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error()})
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	//c.Status(http.StatusNoContent)
	response.Success(c, gin.H{"deleted": true})
}
