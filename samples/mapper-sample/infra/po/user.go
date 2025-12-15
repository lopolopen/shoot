package po

import "gorm.io/gorm"

//go:generate go tool shoot map -path=../../domain/model -type=User

type User struct {
	gorm.Model  `map:".ID"`   //???
	FirstName   string        `gorm:"size:64"`
	LastName    string        `gorm:"size:64"`
	Email       string        `gorm:"size:128;uniqueIndex"`
	AddressList []UserAddress `gorm:"-"`
}
