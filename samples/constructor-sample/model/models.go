package model

//go:generate go tool shoot new -getset -opt -json -file=$GOFILE

type A struct {
	a string
}

type B struct {
	b string
}

type C struct {
	c string
}
