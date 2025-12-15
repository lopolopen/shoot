package po

import (
	"mappersample/domain/model"

	"gorm.io/gorm"
)

//go:generate go tool shoot map -path=../../domain/model -alias=domain -type=User -v

type User struct {
	gorm.Model  `map:".ID"`   //???
	FirstName   string        `gorm:"size:64"`
	LastName    string        `gorm:"size:64"`
	Email       string        `gorm:"size:128;uniqueIndex"`
	AddressList []UserAddress `gorm:"-"`
}

func (u *User) writeModel(user *model.User) {
	user.ID = u.ID
}
