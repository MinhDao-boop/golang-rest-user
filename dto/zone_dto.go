package dto

import (
	"time"

	"gorm.io/datatypes"
)

type ZoneDTORequest struct {
	Name     string         `json:"name"`
	Type     string         `json:"type"`
	Metadata datatypes.JSON `Gorm:"type:json"`
	ParentID uint           `json:"parent_id"`
}

type ZoneDTOResponse struct {
	ID        uint           `json:"id"`
	Name      string         `json:"name"`
	Type      string         `json:"type"`
	Path      string         `json:"path"`
	Level     int            `json:"level"`
	Metadata  datatypes.JSON `Gorm:"type:json"`
	CreatedAt time.Time      `Gorm:"type:datetime"`
	UpdatedAt time.Time      `Gorm:"type:datetime"`
}
