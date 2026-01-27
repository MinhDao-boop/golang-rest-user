package models

type ZoneSOS struct {
	BaseModel
	ZoneID       uint       `gorm:"not null;index"`
	SOSContactID uint       `gorm:"not null;index"`
	Zone         Zone       `gorm:"foreignKey:ZoneID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	SOSContact   SOSContact `gorm:"foreignKey:SOSContactID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
