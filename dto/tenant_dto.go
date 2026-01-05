package dto

import (
	"golang-rest-user/enums"
)

type CreateTenantRequest struct {
	Code   string `json:"code" binding:"required"`
	Name   string `json:"name" binding:"required"`
	DBUser string `json:"db_user" binding:"required"`
	DBPass string `json:"db_pass" binding:"required"`
	DBHost string `json:"db_host" binding:"required"`
	DBPort string `json:"db_port" binding:"required"`
	DBName string `json:"db_name" binding:"required"`
}

type UpdateTenantRequest struct {
	Name   string `json:"name" binding:"required"`
	DBUser string `json:"db_user" binding:"required"`
	DBPass string `json:"db_pass" binding:"required"`
	DBHost string `json:"db_host" binding:"required"`
	DBPort string `json:"db_port" binding:"required"`
}

type TenantResponse struct {
	ID        uint               `json:"id"`
	Code      string             `json:"code"`
	Name      string             `json:"name"`
	DBHost    string             `json:"db_host"`
	DBPort    string             `json:"db_port"`
	DBName    string             `json:"db_name"`
	Status    enums.TenantStatus `json:"status"`
	CreatedAt string             `json:"created_at"`
	UpdatedAt string             `json:"updated_at"`
}
