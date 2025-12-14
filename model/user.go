package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username       string `gorm:"unique"`
	Email          string
	PasswordDigest string
	NickName       string
	Status         string
	Avatar         string
	Money          string
}
