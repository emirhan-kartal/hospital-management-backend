package models

import (
	"gorm.io/gorm"
)

type RedisPolyclinic struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type Polyclinic struct {
	gorm.Model
	ID         int `gorm:"primaryKey;unique"`
	HospitalID int
	Hospital   Hospital `gorm:"foreignKey:HospitalID;references:ID"`
	City       string   `gorm:"type:varchar(100)"`
	District   string   `gorm:"type:varchar(100)"`
	Name       string   `gorm:"type:varchar(255)"`
	Personels  []Personel
}
