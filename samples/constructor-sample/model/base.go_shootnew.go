package model

import (
	"constructorsample/model/dto"
	"encoding/json"
)

// NewBase constructs a new instance of type Base
func NewBase(id int) *Base {
	return &Base{
		id: id,
	}
}

// Id gets the value of field id
func (b *Base) Id() int {
	return b.id
}

// SetId sets the value of field id
func (b *Base) SetId(id_ int) {
	b.id = id_
}

type _Base_marshal struct {
	Id int `json:"id"`
}

type _Base_unmarshal struct {
	Id int `json:"id"`
}

// MarshalJSON serializes type Base to json bytes
func (b Base) MarshalJSON() ([]byte, error) {
	data := _Base_marshal{

		Id: b.Id(),
	}
	return json.Marshal(data)
}

// UnmarshalJSON deserializes json bytes to type Base
func (b *Base) UnmarshalJSON(data []byte) error {
	var base_ _Base_unmarshal
	if err := json.Unmarshal(data, &base_); err != nil {
		return nil
	}
	b.SetId(base_.Id)

	return nil
}

// NewSon constructs a new instance of type Son
func NewSon(name string) *Son {
	return &Son{
		name: name,
	}
}

// Name gets the value of field name
func (s *Son) Name() string {
	return s.name
}

// SetName sets the value of field name
func (s *Son) SetName(name_ string) {
	s.name = name_
}

type _Son_marshal struct {
	Base
	dto.Class

	Name string `json:"name"`
}

type _Son_unmarshal struct {
	Base
	dto.Class

	Name string `json:"name"`
}

// MarshalJSON serializes type Son to json bytes
func (s Son) MarshalJSON() ([]byte, error) {
	data := _Son_marshal{
		Base:  s.Base,
		Class: s.Class,

		Name: s.Name(),
	}
	return json.Marshal(data)
}

// UnmarshalJSON deserializes json bytes to type Son
func (s *Son) UnmarshalJSON(data []byte) error {
	var son_ _Son_unmarshal
	if err := json.Unmarshal(data, &son_); err != nil {
		return nil
	}
	s.SetName(son_.Name)

	s.Base = son_.Base
	s.Class = son_.Class

	return nil
}
