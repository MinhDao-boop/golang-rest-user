package models

import (
	"gorm.io/datatypes"
)

type Zone struct {
	BaseModel
	Name     string         `gorm:"size:255" json:"name"`
	Type     string         `gorm:"size:255" json:"type"`
	Path     string         `gorm:"size:255;uniqueIndex" json:"path"`
	Level    int            `gorm:"index"`
	ParentID *uint          `gorm:"foreignKey:ParentID; references:ID; index; default:NULL; constraint:OnDelete:SET NULL" json:"parent_id"`
	Metadata datatypes.JSON `gorm:"type:json"`
}
