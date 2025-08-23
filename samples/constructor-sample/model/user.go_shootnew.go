package model

import (
	"constructorsample/model/dto"
	"encoding/json"
)

// NewUser constructs a new instance of type User
func NewUser(id string, name string, gender int, age int, tel string) *User {
	return &User{
		id:     id,
		name:   name,
		gender: gender,
		age:    age,
		tel:    tel,
	}
}

// Id gets the value of field id
func (u *User) Id() string {
	return u.id
}

// Name gets the value of field name
func (u *User) Name() string {
	return u.name
}

// Gender gets the value of field gender
func (u *User) Gender() int {
	return u.gender
}

// Age gets the value of field age
func (u *User) Age() int {
	return u.age
}

// Tel gets the value of field tel
func (u *User) Tel() string {
	return u.tel
}

// SetId sets the value of field id
func (u *User) SetId(id_ string) {
	u.id = id_
}

// SetName sets the value of field name
func (u *User) SetName(name_ string) {
	u.name = name_
}

// SetGender sets the value of field gender
func (u *User) SetGender(gender_ int) {
	u.gender = gender_
}

// SetAge sets the value of field age
func (u *User) SetAge(age_ int) {
	u.age = age_
}

// SetTel sets the value of field tel
func (u *User) SetTel(tel_ string) {
	u.tel = tel_
}

type _User_marshal struct {
	Id string `json:"id"`

	Name string `json:"name"`

	Gender int `json:"gender"`

	Age int `json:"age"`

	Tel string `json:"tel"`
}

type _User_unmarshal struct {
	Id string `json:"id"`

	Name string `json:"name"`

	Gender int `json:"gender"`

	Age int `json:"age"`

	Tel string `json:"tel"`
}

// MarshalJSON serializes type User to json bytes
func (u User) MarshalJSON() ([]byte, error) {
	data := _User_marshal{

		Id:     u.Id(),
		Name:   u.Name(),
		Gender: u.Gender(),
		Age:    u.Age(),
		Tel:    u.Tel(),
	}
	return json.Marshal(data)
}

// UnmarshalJSON deserializes json bytes to type User
func (u *User) UnmarshalJSON(data []byte) error {
	var user_ _User_unmarshal
	if err := json.Unmarshal(data, &user_); err != nil {
		return nil
	}
	u.SetId(user_.Id)
	u.SetName(user_.Name)
	u.SetGender(user_.Gender)
	u.SetAge(user_.Age)
	u.SetTel(user_.Tel)

	return nil
}

// NewBook constructs a new instance of type Book
func NewBook(name string, names []string, nameMap map[string]string, userMap map[string]User, owner *User, c *dto.Class) *Book {
	return &Book{
		name:    name,
		names:   names,
		nameMap: nameMap,
		userMap: userMap,
		owner:   owner,
		c:       c,
	}
}

// Name gets the value of field name
func (b *Book) Name() string {
	return b.name
}

// Names gets the value of field names
func (b *Book) Names() []string {
	return b.names
}

// NameMap gets the value of field nameMap
func (b *Book) NameMap() map[string]string {
	return b.nameMap
}

// UserMap gets the value of field userMap
func (b *Book) UserMap() map[string]User {
	return b.userMap
}

// Owner gets the value of field owner
func (b *Book) Owner() *User {
	return b.owner
}

// C gets the value of field c
func (b *Book) C() *dto.Class {
	return b.c
}

// SetName sets the value of field name
func (b *Book) SetName(name_ string) {
	b.name = name_
}

// SetNames sets the value of field names
func (b *Book) SetNames(names_ []string) {
	b.names = names_
}

// SetNameMap sets the value of field nameMap
func (b *Book) SetNameMap(nameMap_ map[string]string) {
	b.nameMap = nameMap_
}

// SetUserMap sets the value of field userMap
func (b *Book) SetUserMap(userMap_ map[string]User) {
	b.userMap = userMap_
}

// SetOwner sets the value of field owner
func (b *Book) SetOwner(owner_ *User) {
	b.owner = owner_
}

// SetC sets the value of field c
func (b *Book) SetC(c_ *dto.Class) {
	b.c = c_
}

type _Book_marshal struct {
	Name string `json:"name"`

	Names []string `json:"names"`

	NameMap map[string]string `json:"nameMap"`

	UserMap map[string]User `json:"userMap"`

	Owner *User `json:"owner"`

	C *dto.Class `json:"c"`
}

type _Book_unmarshal struct {
	Name string `json:"name"`

	Names []string `json:"names"`

	NameMap map[string]string `json:"nameMap"`

	UserMap map[string]User `json:"userMap"`

	Owner *User `json:"owner"`

	C *dto.Class `json:"c"`
}

// MarshalJSON serializes type Book to json bytes
func (b Book) MarshalJSON() ([]byte, error) {
	data := _Book_marshal{

		Name:    b.Name(),
		Names:   b.Names(),
		NameMap: b.NameMap(),
		UserMap: b.UserMap(),
		Owner:   b.Owner(),
		C:       b.C(),
	}
	return json.Marshal(data)
}

// UnmarshalJSON deserializes json bytes to type Book
func (b *Book) UnmarshalJSON(data []byte) error {
	var book_ _Book_unmarshal
	if err := json.Unmarshal(data, &book_); err != nil {
		return nil
	}
	b.SetName(book_.Name)
	b.SetNames(book_.Names)
	b.SetNameMap(book_.NameMap)
	b.SetUserMap(book_.UserMap)
	b.SetOwner(book_.Owner)
	b.SetC(book_.C)

	return nil
}

// NewBook2 constructs a new instance of type Book2
func NewBook2(name string, names []string, owner *User) *Book2 {
	return &Book2{
		name:  name,
		names: names,
		owner: owner,
	}
}

// Name gets the value of field name
func (b *Book2) Name() string {
	return b.name
}

// Names gets the value of field names
func (b *Book2) Names() []string {
	return b.names
}

// Owner gets the value of field owner
func (b *Book2) Owner() *User {
	return b.owner
}

// SetName sets the value of field name
func (b *Book2) SetName(name_ string) {
	b.name = name_
}

// SetNames sets the value of field names
func (b *Book2) SetNames(names_ []string) {
	b.names = names_
}

// SetOwner sets the value of field owner
func (b *Book2) SetOwner(owner_ *User) {
	b.owner = owner_
}

type _Book2_marshal struct {
	Name string `json:"name"`

	Names []string `json:"names"`

	Owner *User `json:"owner"`
}

type _Book2_unmarshal struct {
	Name string `json:"name"`

	Names []string `json:"names"`

	Owner *User `json:"owner"`
}

// MarshalJSON serializes type Book2 to json bytes
func (b Book2) MarshalJSON() ([]byte, error) {
	data := _Book2_marshal{

		Name:  b.Name(),
		Names: b.Names(),
		Owner: b.Owner(),
	}
	return json.Marshal(data)
}

// UnmarshalJSON deserializes json bytes to type Book2
func (b *Book2) UnmarshalJSON(data []byte) error {
	var book2_ _Book2_unmarshal
	if err := json.Unmarshal(data, &book2_); err != nil {
		return nil
	}
	b.SetName(book2_.Name)
	b.SetNames(book2_.Names)
	b.SetOwner(book2_.Owner)

	return nil
}
