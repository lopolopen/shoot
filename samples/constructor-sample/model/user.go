package model

//go:generate go tool shoot new -getset -json -file=$GOFILE -s

type User struct {
	//shoot: get
	id     string
	name   string
	gender int
	age    int
	tel    string
}

type Book struct {
	//shoot: new
	name    string
	writers []string
	Remarks string
	owner   *User
}
