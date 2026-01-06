package handler

import (
	"errors"
	"golang-rest-user/dto"
	"golang-rest-user/provider/serviceProvider"
	"golang-rest-user/response"
	"golang-rest-user/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GET /tenants?page=1&page_size=10&search=...
func ListTenant(c *gin.Context) {
	appService := serviceProvider.GetInstance()
	page, pageSize := utils.GetPageAndPageSize(c)
	search := c.Query("search")

	tenantResponses, total, err := appService.TenantService.List(page, pageSize, search)
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
func CreateTenant(c *gin.Context) {
	appService := serviceProvider.GetInstance()
	var req dto.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	tenantResponse, err := appService.TenantService.Create(req)
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
func GetByTenantCode(c *gin.Context) {
	appService := serviceProvider.GetInstance()
	code := c.Param("code")
	if code == "" {
		response.Error(c, response.CodeBadRequest, "tenant code is required", nil, http.StatusBadRequest)
	}
	tenantResponse, err := appService.TenantService.GetByTenantCode(code)
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
func UpdateTenant(c *gin.Context) {
	appService := serviceProvider.GetInstance()
	code := c.Param("code")
	if code == "" {
		response.Error(c, response.CodeBadRequest, "tenant code is required", nil, http.StatusBadRequest)
	}
	var req dto.UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	tenantResponse, err := appService.TenantService.Update(code, req)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusNotFound)
		return
	}

	response.Success(c, tenantResponse)
}

// DELETE /tenants/:code
func DeleteTenant(c *gin.Context) {
	appService := serviceProvider.GetInstance()
	code := c.Param("code")
	if code == "" {
		response.Error(c, response.CodeBadRequest, "tenant code is required", nil, http.StatusBadRequest)
	}
	if err := appService.TenantService.Delete(code); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusNotFound)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}
