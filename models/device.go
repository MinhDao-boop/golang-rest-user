package models

import (
	"time"

	"gorm.io/datatypes"
)

type Device struct {
	BaseModel
	Code            string `gorm:"size:50; uniqueIndex" json:"code"`
	Name            string `gorm:"size:255" json:"name"`
	Type            string `gorm:"size:255" json:"type"`
	Model           string `gorm:"size:255" json:"model"`
	Manufacturer    string `gorm:"size:255" json:"manufacturer"`
	ZoneID          *uint
	Zone            Zone           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"zone,omitempty"`
	Status          string         `gorm:"size:50; index" json:"status,omitempty"`
	ConnectStatus   string         `gorm:"size:50; index" json:"connect_status,omitempty"`
	LastSeenAt      *time.Time     `gorm:"index" json:"last_seen_at,omitempty"`
	Protocol        string         `gorm:"size:255" json:"protocol"`
	FirmwareVersion string         `gorm:"size:255" json:"firmware_version"`
	HardwareVersion string         `gorm:"size:255" json:"hardware_version"`
	Config          datatypes.JSON `gorm:"type:json" json:"config"`
	DesiredConfig   datatypes.JSON `gorm:"type:json" json:"desired_config,omitempty"`
}
