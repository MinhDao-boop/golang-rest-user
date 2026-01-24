package dto

import (
	"time"

	"gorm.io/datatypes"
)

type ZoneDTORequest struct {
	Name     string         `json:"name"`
	Type     string         `json:"type"`
	Metadata datatypes.JSON `Gorm:"type:json"`
	ParentID *uint          `json:"parent_id"`
}

type ZoneDTOResponse struct {
	ID        uint           `json:"id"`
	UUID      string         `json:"uuid"`
	Name      string         `json:"name"`
	Type      string         `json:"type"`
	Path      string         `json:"path"`
	Level     int            `json:"level"`
	Metadata  datatypes.JSON `json:"metadata"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	ParentID  *uint          `json:"parent_id"`
}
