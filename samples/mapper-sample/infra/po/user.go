package po

import (
	"mappersample/domain/model"

	"gorm.io/gorm"
)

//go:generate go tool shoot map -path=../../domain/model -alias=domain -type=User

type User struct {
	*gorm.Model                //`map:"..."`
	FirstName   string         `gorm:"size:64"`
	LastName    string         `gorm:"size:64"`
	Email       string         `gorm:"size:128;uniqueIndex"`
	AddressList []*UserAddress `gorm:"-"`
	NoMap       NoMap
	B           int
	C           int
}

func (u *User) writeDomain(user *model.User) {
	// user.ID = u.ID
	// user.B = 2
}

func (u *User) readDomain(user model.User) {
	// u.Model = &gorm.Model{
	// 	ID: user.ID,
	// }

	// u.Model.CreatedAt = time.Time{}
	// u.Model.UpdatedAt = time.Time{}
	// u.Model.DeletedAt = gorm.DeletedAt{}

	// u.B = 100
}

type NoMap struct{}
