package model

import (
	dto2 "constructorsample/model/dto"
)

//go:generate go tool shoot new -getset -json -v $GOFILE

type Base struct {
	id int
}

type Son struct {
	Base //todo
	name string
	dto2.Class
	// dto.Class
}
