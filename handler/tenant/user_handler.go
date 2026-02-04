package tenant

import (
	"golang-rest-user/dto"
	"golang-rest-user/provider/tenantProvider"
	"golang-rest-user/response"

	"github.com/gin-gonic/gin"
)

// GET /users?page=1&page_size=10&search=...
func ListUsers(c *gin.Context) {
	tenantCode := c.GetString("tenant_code")
	if tenantCode == "" {
		return
	}
	service := tenantProvider.GetTenantInfo(tenantCode)
	var req dto.ListUsersRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Error(c, response.VAD0000, err)
		return
	}
	userResponse := service.UserService.List(req)
	response.Data(c, userResponse)
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
		response.Error(c, response.VAD0000, err)
		return
	}
	userResponse := service.UserService.Create(req)
	response.Data(c, userResponse)
}

//// GET /users/:user_uuid
//func GetByUserUUID(c *gin.Context) {
//	tenantCode := c.GetString("tenant_code")
//	if tenantCode == "" {
//		return
//	}
//	service := tenantProvider.GetTenantInfo(tenantCode)
//
//	uuid := c.Param("user_uuid")
//	userResponse, err := service.UserService.GetByUUID(uuid)
//	if err != nil {
//		response.Error(c, response.ERR0001, "user not found", nil, http.StatusBadRequest)
//		return
//	}
//	response.Success(c, userResponse)
//}

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
		response.Error(c, response.VAD0000, err)
		return
	}
	userResponse := service.UserService.Update(uuid, req)
	response.Data(c, userResponse)
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
		response.Error(c, response.VAD0000, err)
		return
	}
	userResponse := service.UserService.DeleteMany(req)
	response.Data(c, userResponse)
}
