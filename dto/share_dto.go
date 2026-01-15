package dto

import (
	"golang-rest-user/enums"
	"time"

	"gorm.io/gorm"
)

type ShareDTORequest struct {
	UserID     uint                 `json:"user_id"`
	Permission enums.UserPermission `json:"permission"`
}

type ShareDTOResponse struct {
	UUID       string               `json:"uuid"`
	UserID     uint                 `json:"user_id"`
	ZoneID     uint                 `json:"zone_id"`
	Permission enums.UserPermission `json:"permission"`
	CreatedAt  time.Time            `json:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at"`
	DeletedAt  gorm.DeletedAt       `Gorm:"index" json:"-"`
}
