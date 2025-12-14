package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserId    uint `gorm:"not null"`
	ProductId uint `gorm:"not null"`
	BossId    uint `gorm:"not null"`
	AddressId uint `gorm:"not null"`
	Num       int
	OrderNum  int
	Type      uint //1没有支付 2已经支付
	Money     float64
}
