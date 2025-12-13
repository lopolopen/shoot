package po

import "gorm.io/gorm"

//-go:generate go tool shoot map -path=../../domain/model -type=UserAddress

type UserAddress struct {
	gorm.Model
	City      string `gorm:"size:64"`
	Street    string `gorm:"size:64"`
	Room      string `gorm:"size:64"`
	Tag       string `gorm:"size:32"`
	IsDefault bool   `gorm:"default:false"`
	UserID    uint   `gorm:"index"`
}
