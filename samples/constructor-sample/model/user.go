package model

import "constructorsample/model/dto"

//go:generate go tool shoot new -getset -json -file=$GOFILE

type User struct {
	id     string
	name   string
	gender int
	age    int
	tel    string
}

type Book struct {
	name    string
	names   []string
	nameMap map[string]string
	userMap map[string]User
	owner   *User
	c       *dto.Class
}

type Book2 struct {
	name  string
	names []string
	owner *User
}
