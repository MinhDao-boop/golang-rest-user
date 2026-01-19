package models

import (
	"time"
)

type DeviceCredential struct {
	BaseModel
	DeviceID   *uint
	Device     Device     `Gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"device,omitempty"`
	AuthType   string     `Gorm:"size:255" json:"auth_type"`
	SecretHash string     `Gorm:"size:255" json:"secret_hash"`
	ExpiredAt  *time.Time `orm:"auto_now_add;type(datetime)" json:"expired_at"`
}
