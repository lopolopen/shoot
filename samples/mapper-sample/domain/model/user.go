package model

type User struct { //Aggregate Root
	ID          uint
	FirstName   string
	LastName    string
	Email       string
	AddressList []*UserAddress
	*Model1
}

type Model1 struct {
	A int
	Model2
}

type Model2 struct {
	B int
	*Model3
}

type Model3 struct {
	C int
}
