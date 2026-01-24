package models

import (
	"time"
)

type DeviceEvents struct {
	BaseModel
	DeviceID  *uint
	Device    Device    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"device,omitempty"`
	EventType string    `gorm:"size:255" json:"event_type"`
	Payload   string    `gorm:"size:255" json:"payload"`
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)" json:"created_at"`
}
