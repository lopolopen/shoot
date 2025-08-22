package dto

//go:generate go tool shoot new -getset -json $GOFILE

type Class struct {
	id   int
	name string
}
