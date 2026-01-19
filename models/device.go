package models

import (
	"time"

	"gorm.io/datatypes"
)

type Device struct {
	BaseModel
	Code            string `Gorm:"size:50; uniqueIndex" json:"code"`
	Name            string `Gorm:"size:255" json:"name"`
	Type            string `Gorm:"size:255" json:"type"`
	Model           string `Gorm:"size:255" json:"model"`
	Manufacturer    string `Gorm:"size:255" json:"manufacturer"`
	ZoneID          *uint
	Zone            Zone           `Gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"zone,omitempty"`
	Status          string         `Gorm:"size:50; index" json:"status,omitempty"`
	ConnectStatus   string         `Gorm:"size:50; index" json:"connect_status,omitempty"`
	LastSeenAt      *time.Time     `Gorm:"index" json:"last_seen_at,omitempty"`
	Protocol        string         `Gorm:"size:255" json:"protocol"`
	FirmwareVersion string         `Gorm:"size:255" json:"firmware_version"`
	HardwareVersion string         `Gorm:"size:255" json:"hardware_version"`
	Config          datatypes.JSON `Gorm:"type:json" json:"config"`
	DesiredConfig   datatypes.JSON `Gorm:"type:json" json:"desired_config,omitempty"`
}
