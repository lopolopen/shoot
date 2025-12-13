package model

type UserAddress struct { //entity
	ID        uint
	City      string
	Street    string
	Room      string
	Tag       string
	IsDefault bool
}

type OrderAddress struct { //value object
	City   string
	Street string
	Room   string
}
