package models

import (
	"gorm.io/gorm"
)

type Hospital struct {
	gorm.Model
	HospitalBody

	User        []User
	Polyclinics []Polyclinic
}
type HospitalBody struct {
	Name         string `gorm:"type:varchar(255);not null;default:null"`
	TaxID        string `gorm:"type:varchar(255);unique;not null;default:null"`
	Email        string `gorm:"type:varchar(255);unique;not null;default:null"`
	TelNo        string `gorm:"type:varchar(255);unique;not null;default:null"`
	City         string `gorm:"type:varchar(255);not null;default:null"`
	District     string `gorm:"type:varchar(255);not null;default:null"`
	AdressDetail string `gorm:"type:text;not null;default:null"`
}
