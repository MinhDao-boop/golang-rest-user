package dto

import (
	"golang-rest-user/enums"
	"time"
)

type ZoneIdQuery struct {
	ZoneUuid string `form:"zone_uuid" binding:"required"`
}
type CreateSOSContactRequest struct {
	Name  string `json:"name" binding:"required"`
	Role  string `json:"role" binding:"required"`
	Phone string `json:"phone" binding:"required"`
	Note  string `json:"note"`
}

type ListSOSContactRequest struct {
	ZoneUuid string `form:"zone_uuid" binding:"required"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=50"`
	Search   string `form:"search"`
	IsAll    bool   `form:"is_all"`
	IsActive *bool  `form:"is_active"`
}

type UpdateSOSContactRequest struct {
	Name  *string `json:"name" binding:"required" patch:"true"`
	Role  *string `json:"role" binding:"required" patch:"true"`
	Phone *string `json:"phone" binding:"required" patch:"true"`
	Note  string  `json:"note" patch:"true"`
}

type ToggleSOSContactRequest struct {
	IsActive *enums.SOSContactStatus `json:"is_active" patch:"true"`
}

type DeleteSOSContactRequest struct {
	Ids []uint `json:"ids"`
}

type SOSContactResponse struct {
	ID        uint                   `json:"id"`
	UUID      string                 `json:"uuid"`
	Name      string                 `json:"name"`
	Role      string                 `json:"role"`
	Phone     string                 `json:"phone"`
	Note      string                 `json:"note"`
	IsActive  enums.SOSContactStatus `json:"is_active"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}
