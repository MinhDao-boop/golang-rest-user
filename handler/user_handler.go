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

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(s service.UserService) *UserHandler {
	return &UserHandler{svc: s}
}

// GET /users?page=1&page_size=10&search=...
func (h *UserHandler) ListUsers(c *gin.Context) {
	tenantDB, err := utils.GetTenantDB(c)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
	}
	page, pageSize := utils.GetPageAndPageSize(c)
	search := c.Query("search")

	userResponses, total, err := h.svc.List(tenantDB, page, pageSize, search)
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
func (h *UserHandler) CreateUser(c *gin.Context) {
	tenantDB, err := utils.GetTenantDB(c)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
	}
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	userResponse, err := h.svc.Create(tenantDB, req)
	if err != nil {
		if strings.Contains(err.Error(), "exists") {
			response.Error(c, response.CodeBadRequest, "username already exists", nil, http.StatusConflict)
			return
		}
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	location := c.Request.URL.Path + "/" + strconv.FormatUint(uint64(userResponse.ID), 10)
	c.Header("Location", location)
	response.Success(c, userResponse)
}

// GET /users/:uuid
func (h *UserHandler) GetByUserUUID(c *gin.Context) {
	tenantDB, err := utils.GetTenantDB(c)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
	}
	uuid := c.Param("uuid")
	userResponse, err := h.svc.GetByUUID(tenantDB, uuid)
	if err != nil {
		response.Error(c, response.CodeBadRequest, "user not found", nil, http.StatusBadRequest)
		return
	}
	response.Success(c, userResponse)
}

// PUT /users/:uuid
func (h *UserHandler) UpdateUser(c *gin.Context) {
	tenantDB, err := utils.GetTenantDB(c)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
	}
	uuid := c.Param("uuid")

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	userResponse, err := h.svc.Update(tenantDB, uuid, req)
	if err != nil {
		response.Error(c, response.CodeBadRequest, "user not found", nil, http.StatusNotFound)
		return
	}
	response.Success(c, userResponse)
}

// DELETE /users?uuids=1b0f0fe4-8710-4518-b8bc-7f1e52b280e4,1c8edc4f-b1a0-4252-808b-682eb76551ad,...
func (h *UserHandler) DeleteManyUsers(c *gin.Context) {
	tenantDB, err := utils.GetTenantDB(c)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
	}
	uuidsParam := c.Query("uuids")
	if uuidsParam == "" {
		response.Error(c, response.CodeBadRequest, "ids query param required", nil, http.StatusBadRequest)
		return
	}
	parts := strings.Split(uuidsParam, ",")
	uuids := []string{}
	for _, p := range parts {
		if p == "" {
			continue
		}
		uuids = append(uuids, strings.TrimSpace(p))
	}
	deleted, err := h.svc.DeleteMany(tenantDB, uuids)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, response.CodeBadRequest, err.Error(), gin.H{"deleted": deleted}, http.StatusBadRequest)
		}
		response.Error(c, response.CodeBadRequest, err.Error(), gin.H{"deleted": deleted}, http.StatusInternalServerError)
		return
	}
	response.Success(c, gin.H{"deleted": deleted})
}
