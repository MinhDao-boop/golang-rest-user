package models

import (
	"gorm.io/datatypes"
)

type Zone struct {
	BaseModel
	Name     string         `Gorm:"size:255;uniqueIndex" json:"name"`
	Type     string         `Gorm:"size:255" json:"type"`
	Path     string         `Gorm:"size:255;uniqueIndex" json:"path"`
	Level    int            `Gorm:"index"`
	ParentID *uint          `Gorm:"foreignKey:ParentID; references:ID; index; default:NULL" json:"parent_id"`
	Metadata datatypes.JSON `Gorm:"type:json"`
	//Children []Zone         `Gorm:"-" json:"children,omitempty"`
}
