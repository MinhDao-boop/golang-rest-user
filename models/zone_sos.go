package models

type ZoneSOS struct {
	BaseModel
	ZoneID       uint       `gorm:"index"`
	SOSContactID uint       `gorm:"index"`
	Zone         Zone       `gorm:"foreignKey:ZoneID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SOSContact   SOSContact `gorm:"foreignKey:SOSContactID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
