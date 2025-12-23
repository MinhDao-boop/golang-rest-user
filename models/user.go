package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"size:255;uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"size:255;not null" json:"password"`
	FullName  string         `gorm:"size:255" json:"full_name"`
	Phone     string         `gorm:"size:50" json:"phone"`
	Position  string         `gorm:"size:255" json:"position"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
