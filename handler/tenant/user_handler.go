package tenant

import (
	"golang-rest-user/dto"
	"golang-rest-user/provider/tenantProvider"
	"golang-rest-user/response"
	"golang-rest-user/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// GET /users?page=1&page_size=10&search=...
func ListUsers(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	service := tenantProvider.GetTenantInfo(tenantCode)

	page, pageSize := utils.GetPageAndPageSize(c)
	search := c.Query("search")

	userResponses, total, err := service.UserService.List(page, pageSize, search)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	response.Success(c, gin.H{
		"data":      userResponses,
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	})
}

// POST /users
func CreateUser(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	service := tenantProvider.GetTenantInfo(tenantCode)

	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	userResponse, err := service.UserService.Create(req)
	if err != nil {
		if strings.Contains(err.Error(), "exists") {
			response.Error(c, response.CodeBadRequest, "username already exists", nil, http.StatusConflict)
			return
		}
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	//location := c.Request.URL.Path + "/" + strconv.FormatUint(uint64(userResponse.ID), 10)
	//c.Header("Location", location)
	response.Success(c, userResponse)
}

// GET /users/:user_uuid
func GetByUserUUID(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	service := tenantProvider.GetTenantInfo(tenantCode)

	uuid := c.Param("user_uuid")
	userResponse, err := service.UserService.GetByUUID(uuid)
	if err != nil {
		response.Error(c, response.CodeBadRequest, "user not found", nil, http.StatusBadRequest)
		return
	}
	response.Success(c, userResponse)
}

// PUT /users/:user_uuid
func UpdateUser(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	service := tenantProvider.GetTenantInfo(tenantCode)
	uuid := c.Param("user_uuid")

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	userResponse, err := service.UserService.Update(uuid, req)
	if err != nil {
		response.Error(c, response.CodeBadRequest, "user not found", nil, http.StatusNotFound)
		return
	}
	response.Success(c, userResponse)
}

// DELETE /users
func DeleteManyUsers(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	service := tenantProvider.GetTenantInfo(tenantCode)
	var req dto.DeleteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	deleted, err := service.UserService.DeleteMany(req)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), gin.H{"deleted": deleted}, http.StatusBadRequest)
		return
	}
	response.Success(c, gin.H{"deleted": deleted})
}
