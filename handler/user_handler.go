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
	db := c.MustGet("TENANT_DB").(*gorm.DB)
	userRepo := repository.NewUserRepo(db)
	userSvc := service.NewUserService(userRepo)
	return userSvc
}

// GET /users?page=1&page_size=10&search=...
func (h *UserHandler) ListUsersResponse(c *gin.Context) {
	svc := getUserService(c)
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
		//c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	resp := []dto.UserResponse{}
	for _, u := range users {
		resp = append(resp, dto.UserResponse{
			ID: u.ID, Username: u.Username, FullName: u.FullName, Phone: u.Phone, Position: u.Position,
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
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	user, err := svc.Create(req)
	if err != nil {
		if strings.Contains(err.Error(), "exists") {
			//c.JSON(http.StatusConflict, gin.H{"success": false, "error": "username already exists"})
			response.Error(c, response.CodeBadRequest, "username already exists", nil, http.StatusBadRequest)
			return
		}
		//c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	location := c.Request.URL.Path + "/" + strconv.FormatUint(uint64(user.ID), 10)
	c.Header("Location", location)
	response.Success(c, dto.UserResponse{
		ID: user.ID, Username: user.Username, FullName: user.FullName, Phone: user.Phone, Position: user.Position,
		CreatedAt: user.CreatedAt.Format(time.RFC3339), UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	})
}

// GET /users/:id
func (h *UserHandler) GetByUserID(c *gin.Context) {
	svc := getUserService(c)
	id64, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	id := uint(id64)
	user, err := svc.GetByID(id)
	if err != nil {
		//c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "user not found"})
		response.Error(c, response.CodeBadRequest, "user not found", nil, http.StatusBadRequest)
		return
	}
	response.Success(c, dto.UserResponse{
		ID: user.ID, Username: user.Username, FullName: user.FullName, Phone: user.Phone, Position: user.Position,
		CreatedAt: user.CreatedAt.Format(time.RFC3339), UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	})
}

// PUT /users/:id
func (h *UserHandler) UpdateUserRequest(c *gin.Context) {
	svc := getUserService(c)
	id64, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	id := uint(id64)

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		//c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}

	user, err := svc.Update(id, req)
	if err != nil {
		//c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "user not found"})
		response.Error(c, response.CodeBadRequest, "user not found", nil, http.StatusBadRequest)
		return
	}
	response.Success(c, dto.UserResponse{
		ID: user.ID, Username: user.Username, FullName: user.FullName, Phone: user.Phone, Position: user.Position,
		CreatedAt: user.CreatedAt.Format(time.RFC3339), UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	})
}

// DELETE /users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	svc := getUserService(c)
	id64, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	id := uint(id64)
	if err := svc.Delete(id); err != nil {
		//c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "user not found"})
		response.Error(c, response.CodeBadRequest, "user not found", nil, http.StatusBadRequest)
		return
	}
	//c.Status(http.StatusNoContent)
	response.Success(c, gin.H{"deleted": true})
}

// DELETE /users?ids=1,2,3
func (h *UserHandler) DeleteManyUsers(c *gin.Context) {
	svc := getUserService(c)
	idsParam := c.Query("ids")
	if idsParam == "" {
		//c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "ids query param required"})
		response.Error(c, response.CodeBadRequest, "ids query param required", nil, http.StatusBadRequest)
		return
	}
	parts := strings.Split(idsParam, ",")
	ids := []uint{}
	for _, p := range parts {
		if p == "" {
			continue
		}
		v, err := strconv.ParseUint(strings.TrimSpace(p), 10, 64)
		if err != nil {
			//c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "invalid id in ids"})
			response.Error(c, response.CodeBadRequest, "invalid id in ids", nil, http.StatusBadRequest)
			return
		}
		ids = append(ids, uint(v))
	}
	deleted, err := svc.DeleteMany(ids)
	if err != nil {
		//c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		response.Error(c, response.CodeBadRequest, err.Error(), nil, http.StatusBadRequest)
		return
	}
	if deleted == 0 {
		//c.Status(http.StatusNoContent)
		response.Error(c, response.CodeBadRequest, "no users deleted", nil, http.StatusBadRequest)
		return
	}
	//c.JSON(http.StatusOK, gin.H{"success": true, "deleted": deleted})
	response.Success(c, gin.H{"deleted": deleted})
}
