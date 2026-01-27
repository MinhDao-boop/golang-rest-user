package models

type ZoneSOS struct {
	BaseModel
	ZoneID       uint       `gorm:"not null;index"`
	SOSContactID uint       `gorm:"not null;index"`
	Zone         Zone       `gorm:"foreignKey:ZoneID"`
	SOSContact   SOSContact `gorm:"foreignKey:SOSContactID"`
}
