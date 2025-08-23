package dto

//go:generate go tool shoot new -getset -json -file=$GOFILE

type Class struct {
	id   int
	name string
}
