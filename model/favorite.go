package model

import (
	"gorm.io/gorm"
)

type Favorite struct {
	gorm.Model
	User      User    `gorm:"foreignkey:UserID"`
	UserId    uint    `gorm:"not null"`
	Product   Product `gorm:"foreignkey:ProductID"`
	ProductID uint    `gorm:"not null"`
	Boss      User    `gorm:"foreignkey:BossID"`
	BossID    uint    `gorm:"not null"`
}
