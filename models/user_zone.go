package models

import (
	"golang-rest-user/enums"
	"time"
)

type UserZone struct {
	BaseModel
	UserID     uint                 `gorm:"index"`
	ZoneID     uint                 `gorm:"index"`
	User       User                 `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Zone       Zone                 `gorm:"foreignKey:ZoneID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Permission enums.UserPermission `gorm:"tinyint;not null" json:"permission"`
	ExpiredAt  time.Time            `gorm:"index;default:null"`
	ShareScope uint                 `gorm:"tinyint;default:0;" json:"share_scope"`
	Status     uint                 `gorm:"tinyint;default:0" json:"status"`
}
