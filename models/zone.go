package models

import (
	"gorm.io/datatypes"
)

type Zone struct {
	BaseModel
	Name     string         `gorm:"primaryKey" json:"name"`
	Type     string         `gorm:"size:255" json:"type"`
	Path     string         `gorm:"size:255" json:"path"`
	Level    int            `gorm:"index"`
	ParentID *uint          `gorm:"foreignKey:ParentID; references:ID" json:"parent_id"`
	Metadata datatypes.JSON `Gorm:"type:json"`
	Children []Zone         `gorm:"-" json:"children,omitempty"`
}
