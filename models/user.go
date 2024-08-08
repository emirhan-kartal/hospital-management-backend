package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           int       `gorm:"primaryKey"`
	HospitalID   int       `gorm:"unique;not null"`
	Hospital     *Hospital `gorm:"foreignKey:HospitalID;references:ID;"`
	Name         string    `gorm:"type:varchar(255);not null;default:null"`
	Surname      string    `gorm:"type:varchar(255)";not null;default:null`
	TcNo         string    `gorm:"type:varchar(11);unique;not null;default:null"`
	Email        string    `gorm:"type:varchar(255);unique;not null;default:null"`
	TelNo        string    `gorm:"type:varchar(255);unique;not null;default:null"`
	Password     string    `gorm:"type:varchar(255); not nulldefault:null"`
	Role         string    `gorm:"type:varchar(50); not nulldefault:null"`
	TokenVersion int       `gorm:"type:int; not null;default:0"`
}
