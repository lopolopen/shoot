package po

import (
	"mappersample/domain/model"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model  `map:".ID"`    //???
	FirstName   string         `gorm:"size:64"`
	LastName    string         `gorm:"size:64"`
	Email       string         `gorm:"size:128;uniqueIndex"`
	AddressList []*UserAddress `gorm:"-"`
}

func (u *User) writeDomain(user *model.User) {
	user.ID = u.ID
}
