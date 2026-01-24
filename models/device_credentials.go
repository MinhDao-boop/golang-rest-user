package models

import (
	"time"
)

type DeviceCredential struct {
	BaseModel
	DeviceID   *uint
	Device     Device     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"device,omitempty"`
	AuthType   string     `gorm:"size:255" json:"auth_type"`
	SecretHash string     `gorm:"size:255" json:"secret_hash"`
	ExpiredAt  *time.Time `orm:"auto_now_add;type(datetime)" json:"expired_at"`
}
