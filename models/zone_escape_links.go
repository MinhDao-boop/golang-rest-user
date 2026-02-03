package models

import "gorm.io/datatypes"

type ZoneEscapeLink struct {
	BaseModel
	ZoneID  uint
	Zone    Zone           `gorm:"foreignKey:ZoneID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	WebView datatypes.JSON `gorm:"type:json;not null"`
}
