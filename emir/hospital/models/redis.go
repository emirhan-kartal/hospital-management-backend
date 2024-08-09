package models

type City struct {
	ID        int        `json:"id"`
	Name      string     `json:"name";gorm:"type:varchar(255);not null;default:null;unique"`
	Districts []District `json:"districts"`
}
type District struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type JobType struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Titles []Title `json:"subtitles"`
}
type Title struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
