package models

import (
	"golang-rest-user/enums"
	"time"

	"gorm.io/gorm"
)

type Tenant struct {
	ID        uint               `gorm:"primaryKey" json:"id"`
	Code      string             `gorm:"size:45; uniqueIndex" json:"code"`
	Name      string             `gorm:"size:255; not null" json:"name"`
	DBUser    string             `gorm:"size:50" json:"db_user"`
	DBPass    string             `gorm:"size:50" json:"db_pass"`
	DBHost    string             `gorm:"size:50" json:"db_host"`
	DBPort    string             `gorm:"size:50" json:"db_port"`
	DBName    string             `gorm:"size:50; uniqueIndex" json:"db_name"`
	Status    enums.TenantStatus `gorm:"type:enum('active', 'inactive'); default:'active'" json:"status"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	DeletedAt gorm.DeletedAt     `gorm:"index" json:"-"`
}
