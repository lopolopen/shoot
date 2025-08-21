package model

//go:generate go tool shoot new -getset -json -type=User

type User struct {
	id     string
	name   string
	gender int
	age    int
	tel    string
}
