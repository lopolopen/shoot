package model

import "encoding/json"

//go:generate go tool shoot new -type User

type User struct {
	name   string `jsonx:"name"` //param;get
	age    int    `jsonx:"age"`  //param;get;set
	Gender string `json:"gender"`
}

// type User2 struct {
// 	name     string `jsonx:"name" new:"param,get"`
// 	age      int    `jsonx:"age" new:"param,get,set"`
// 	Gender   string `json:"gender"`
// 	_private string
// }

// func NewUser(name string, age int) *User {
// 	return &User{
// 		name: name,
// 		age:  age,
// 	}
// }

// func (u *User) Name() string {
// 	return u.name
// }

// func (u *User) Age() int {
// 	return u.age
// }

// func (u *User) SetAge(age int) {
// 	u.age = age
// }

type _user_ struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Gender string `json:"gender"`
}

func (u *User) MarshalJSON() ([]byte, error) {
	data := _user_{
		Name:   u.Name(),
		Age:    u.Age(),
		Gender: u.Gender,
	}
	return json.Marshal(data)
}

func (u *User) UnmarshalJSON(data []byte) error {
	var user _user_
	if err := json.Unmarshal(data, &user); err != nil {
		return err
	}
	u.name = user.Name
	u.age = user.Age
	u.Gender = user.Gender
	return nil
}
