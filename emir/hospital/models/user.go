package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserBody
	ID           int       `gorm:"primaryKey"`
	Hospital     *Hospital `gorm:"foreignKey:HospitalID;references:ID;"`
	TokenVersion int       `gorm:"type:int; not null;default:0"`
}

type UserBody struct {
	Name       string `gorm:"type:varchar(255);not null;default:null"`
	Surname    string `gorm:"type:varchar(255)";not null;default:null`
	TcNo       string `gorm:"type:varchar(11);unique;not null;default:null"`
	Email      string `gorm:"type:varchar(255);unique;not null;default:null"`
	TelNo      string `gorm:"type:varchar(255);unique;not null;default:null"`
	Password   string `gorm:"type:varchar(255); not nulldefault:null"`
	Role       string `gorm:"type:varchar(50); not nulldefault:null"`
	HospitalID uint   `gorm:"not null"`
}
