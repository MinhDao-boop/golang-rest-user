package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang-rest-user/dto"
	"golang-rest-user/repository"
	"golang-rest-user/response"
	"golang-rest-user/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func getUserService(c *gin.Context) service.UserService {
	dbAny, ok := c.Get("TENANT_DB")
	if !ok {
		return nil
	}
	db := dbAny.(*gorm.DB)
	userRepo := repository.NewUserRepo(db)
	userSvc := service.NewUserService(userRepo)
	return userSvc
}

// GET /users?page=1&page_size=10&search=...
func (h *UserHandler) ListUsersResponse(c *gin.Context) {
	svc := getUserService(c)
	if svc == nil {
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	search := c.Query("search")

	users, total, err := svc.List(page, pageSize, search)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	resp := []dto.UserResponse{}
	for _, u := range users {
		resp = append(resp, dto.UserResponse{
			ID: u.ID, Uuid: u.Uuid, Username: u.Username, Password: u.Password, FullName: u.FullName, Phone: u.Phone, Position: u.Position,
			CreatedAt: u.CreatedAt.Format(time.RFC3339), UpdatedAt: u.UpdatedAt.Format(time.RFC3339),
		})
	}

	response.Success(c, dto.ListUsersResponse{
		Data: resp, Page: page, PageSize: pageSize, Total: total,
	})
}

// POST /users
func (h *UserHandler) CreateUser(c *gin.Context) {
	svc := getUserService(c)
	if svc == nil {
		return
	}
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	user, err := svc.Create(req)
	if err != nil {
		if strings.Contains(err.Error(), "exists") {
			response.Error(c, response.CodeBadRequest, "username already exists", nil, http.StatusConflict)
			return
		}
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	location := c.Request.URL.Path + "/" + strconv.FormatUint(uint64(user.ID), 10)
	c.Header("Location", location)
	response.Success(c, dto.UserResponse{
		ID: user.ID, Uuid: user.Uuid, Username: user.Username, Password: user.Password, FullName: user.FullName, Phone: user.Phone, Position: user.Position,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	})
}

// GET /users/:uuid
func (h *UserHandler) GetByUserUUID(c *gin.Context) {
	svc := getUserService(c)
	if svc == nil {
		return
	}
	uuid := c.Param("uuid")
	user, err := svc.GetByUUID(uuid)
	if err != nil {
		response.Error(c, response.CodeBadRequest, "user not found", nil, http.StatusBadRequest)
		return
	}
	response.Success(c, dto.UserResponse{
		ID: user.ID, Uuid: user.Uuid, Username: user.Username, Password: user.Password, FullName: user.FullName, Phone: user.Phone, Position: user.Position,
		CreatedAt: user.CreatedAt.Format(time.RFC3339), UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	})
}

// PUT /users/:uuid
func (h *UserHandler) UpdateUserRequest(c *gin.Context) {
	svc := getUserService(c)
	if svc == nil {
		return
	}
	uuid := c.Param("uuid")

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	user, err := svc.Update(uuid, req)
	if err != nil {
		response.Error(c, response.CodeBadRequest, "user not found", nil, http.StatusNotFound)
		return
	}
	response.Success(c, dto.UserResponse{
		ID: user.ID, Uuid: user.Uuid, Username: user.Username, Password: user.Password, FullName: user.FullName, Phone: user.Phone, Position: user.Position,
		CreatedAt: user.CreatedAt.Format(time.RFC3339), UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	})
}

// DELETE /users/:uuid
func (h *UserHandler) DeleteUser(c *gin.Context) {
	svc := getUserService(c)
	if svc == nil {
		return
	}
	uuid := c.Param("uuid")
	if err := svc.Delete(uuid); err != nil {
		response.Error(c, response.CodeBadRequest, "user not found", nil, http.StatusNotFound)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

// DELETE /users?uuids=1b0f0fe4-8710-4518-b8bc-7f1e52b280e4,1c8edc4f-b1a0-4252-808b-682eb76551ad,...
func (h *UserHandler) DeleteManyUsers(c *gin.Context) {
	svc := getUserService(c)
	if svc == nil {
		return
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
		uuids = append(uuids, p)
	}
	deleted, err := svc.DeleteMany(uuids)
	if err != nil {
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusInternalServerError)
		return
	}
	if deleted == 0 {
		response.Error(c, response.CodeBadRequest, "no users deleted", nil, http.StatusNoContent)
		return
	}
	response.Success(c, gin.H{"deleted": deleted})
}
