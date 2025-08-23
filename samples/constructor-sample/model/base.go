package model

import "constructorsample/model/dto"

//go:generate go tool shoot new -getset -json -file=$GOFILE

type Base struct {
	id int
}

type Son struct {
	Base //todo
	name string
	dto.Class
	// dto.Class
}
