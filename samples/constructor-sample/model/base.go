package model

//go:generate go tool shoot new -getset -json $GOFILE

type Base struct {
	id int
}

type Son struct {
	Base //todo
	name string
}

func test() {
	s := NewSon("")
	s.Base = *NewBase(1)
}
