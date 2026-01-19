package models

import (
	"time"
)

type DeviceTelemetry struct {
	BaseModel
	DeviceID  *uint
	Device    Device    `Gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"device,omitempty"`
	Metric    string    `Gorm:"size:255" json:"metric"`
	Value     uint      `Gorm:"size:255" json:"value"`
	Timestamp time.Time `orm:"auto_now_add;type(datetime)" json:"timestamp"`
}
