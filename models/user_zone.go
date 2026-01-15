package models

import "golang-rest-user/enums"

type UserZone struct {
	BaseModel
	UserID     uint                 `gorm:"primaryKey" json:"user_id"`
	ZoneID     uint                 `gorm:"primaryKey" json:"zone_id"`
	Permission enums.UserPermission `json:"permission"`
}
