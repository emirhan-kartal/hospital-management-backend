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
	PolyclinicBody
	ID         int `gorm:"primaryKey;unique"`
	HospitalID int
	Hospital   Hospital `gorm:"foreignKey:HospitalID;references:ID"`
	Personels  []Personel
}
type PolyclinicBody struct {
	Name     string `json:"polyclinic_name"`
	City     string `json:"city"`
	District string `json:"district"`
}
