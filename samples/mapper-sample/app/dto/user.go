package dto

import (
	"mappersample/domain/model"
	"mappersample/infra/mapper"
	"strings"
)

//-go:generate go tool shoot map -path=../../domain/model -type=User

type User struct {
	*mapper.Mapper
	ID       uint
	FullName string
	Email    string
}

func (u User) toModel(user *model.User) { //or writeModel
	parts := strings.Fields(u.FullName)
	if len(parts) == 1 {
		user.FirstName = parts[0]
	}
	if len(parts) > 1 {
		user.LastName = parts[1]
	}
}

func (u *User) fromModel(user model.User) { //or readModel
	u.FullName = user.FirstName + " " + user.LastName
}

// func (u *User) writeModel(user *model.User) {
// 	//err: found more than one manual write method: writeModel
// }
