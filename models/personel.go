package models

import "gorm.io/gorm"

type Personel struct {
	gorm.Model
	ID           int `gorm:"primaryKey"`
	PolyclinicID int
	Polyclinic   Polyclinic `gorm:"foreignKey:PolyclinicID;references:ID" json:"polyclinic"`
	HospitalID   int        `gorm:"not null" json:"hospital_id"`
	Name         string     `gorm:"type:varchar(255)"json:"name"`
	Surname      string     `gorm:"type:varchar(255)"json:"surname"`
	TcNo         string     `gorm:"type:varchar(11);unique" json:"tc_no"`
	TelNo        string     `gorm:"type:varchar(20);unique"json:"tel_no"`
	JobType      string     `gorm:"type:varchar(100)"json:"job_type"`
	Title        string     `gorm:"type:varchar(100)"json:"title"`
	WorkingDays  string     `gorm:"type:varchar(100)"json:"working_days"`
}
