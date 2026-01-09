package models

type User struct {
	BaseModel
	UUID     string `gorm:"size:255; not null" json:"uuid"`
	Username string `gorm:"size:255;uniqueIndex;not null" json:"username"`
	Password string `gorm:"size:255;not null" json:"password"`
	FullName string `gorm:"size:255" json:"full_name"`
	Phone    string `gorm:"size:50" json:"phone"`
	Position string `gorm:"size:255" json:"position"`
}
