package models

import (
	"time"
)

type DeviceTelemetry struct {
	BaseModel
	DeviceID  *uint
	Device    Device    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"device,omitempty"`
	Metric    string    `gorm:"size:255" json:"metric"`
	Value     uint      `gorm:"size:255" json:"value"`
	Timestamp time.Time `orm:"auto_now_add;type(datetime)" json:"timestamp"`
}
