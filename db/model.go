package db

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"unique;not null"`
	Phone    string `gorm:"unique;not null"`
	Password string
	Cookie   string
	Key      string
}

type Site struct {
	gorm.Model
	Img  string
	Name string
	Uid  uint
}

type Device struct {
	gorm.Model
	PosX   float64
	PosY   float64
	Zoom   float64
	Img    string
	Name   string
	Status string
	Type   uint
	Sid    uint
}

type Log struct {
	gorm.Model
	Sid          uint
	Did          uint
	Name         string
	StatusBefore string
	StatusAfter  string
}
