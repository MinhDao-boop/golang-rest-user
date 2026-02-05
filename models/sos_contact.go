package models

type SOSContact struct {
	BaseModel
	Name     string `gorm:"type:varchar(50);not null" json:"name;required"`
	Role     string `gorm:"type:varchar(50);not null" json:"role;required"`
	Phone    string `gorm:"type:varchar(20);not null" json:"phone;required"`
	IsActive bool   `gorm:"type:tinyint(1);not null" json:"is_active;required"`
	Note     string `gorm:"type:text" json:"note"`
}
