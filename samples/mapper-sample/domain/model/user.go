package model

type User struct {
	ID          uint
	FirstName   string
	LastName    string
	Email       string
	AddressList []UserAddress
}
