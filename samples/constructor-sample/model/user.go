package model

import (
	dto "constructorsample/dto"
	dto2 "constructorsample/model/dto"
)

//go:generate go tool shoot new -getset -json -type=User

type User struct {
	id     string
	name   string
	gender int
	age    int
	tel    string
}

//go:generate go tool shoot new -type=Book

type Book struct {
	name    string
	names   []string
	nameMap map[string]string
	userMap map[string]User
	owner   *User
	c       *dto.Class
	c2      *dto2.Class
}

//go:generate go tool shoot new -getset -json -type=Book2

type Book2 struct {
	name  string
	names []string
	owner *User
}
