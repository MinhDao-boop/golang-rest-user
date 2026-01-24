package models

import "golang-rest-user/enums"

type SOSContact struct {
	BaseModel
	Name     string                 `gorm:"type:varchar(50);not null" json:"name;required"`
	Position string                 `gorm:"type:varchar(50);not null" json:"position;required"`
	Phone    string                 `gorm:"type:varchar(20);not null;uniqueIndex" json:"phone;required"`
	IsActive enums.SOSContactStatus `gorm:"type:tinyint(1);not null" json:"is_active;required"`
	Note     string                 `gorm:"type:text" json:"note"`
	ZoneID   *uint
	Zone     Zone `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"zone,omitempty"`
}
