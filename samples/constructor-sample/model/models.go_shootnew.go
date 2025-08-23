package model

import "encoding/json"

// NewA constructs a new instance of type A
func NewA(a string) *A {
	return &A{
		a: a,
	}
}

// With initializes this instance using the functional options pattern
func (a *A) With(opts ..._opt_[A, *A]) *A {
	a._def_()
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func AOfA(a_ string) _opt_[A, *A] {
	return func(a *A) {
		a.a = a_
	}
}

func (a *A) _def_() {

}

// A gets the value of field a
func (a *A) A() string {
	return a.a
}

// SetA sets the value of field a
func (a *A) SetA(a_ string) {
	a.a = a_
}

type _A_marshal struct {
	A string `json:"a"`
}

type _A_unmarshal struct {
	A string `json:"a"`
}

// MarshalJSON serializes type A to json bytes
func (a A) MarshalJSON() ([]byte, error) {
	data := _A_marshal{

		A: a.A(),
	}
	return json.Marshal(data)
}

// UnmarshalJSON deserializes json bytes to type A
func (a *A) UnmarshalJSON(data []byte) error {
	var a_ _A_unmarshal
	if err := json.Unmarshal(data, &a_); err != nil {
		return nil
	}
	a.SetA(a_.A)

	return nil
}

// NewB constructs a new instance of type B
func NewB(b string) *B {
	return &B{
		b: b,
	}
}

// With initializes this instance using the functional options pattern
func (b *B) With(opts ..._opt_[B, *B]) *B {
	b._def_()
	for _, opt := range opts {
		opt(b)
	}
	return b
}

func BOfB(b_ string) _opt_[B, *B] {
	return func(b *B) {
		b.b = b_
	}
}

func (b *B) _def_() {

}

// B gets the value of field b
func (b *B) B() string {
	return b.b
}

// SetB sets the value of field b
func (b *B) SetB(b_ string) {
	b.b = b_
}

type _B_marshal struct {
	B string `json:"b"`
}

type _B_unmarshal struct {
	B string `json:"b"`
}

// MarshalJSON serializes type B to json bytes
func (b B) MarshalJSON() ([]byte, error) {
	data := _B_marshal{

		B: b.B(),
	}
	return json.Marshal(data)
}

// UnmarshalJSON deserializes json bytes to type B
func (b *B) UnmarshalJSON(data []byte) error {
	var b_ _B_unmarshal
	if err := json.Unmarshal(data, &b_); err != nil {
		return nil
	}
	b.SetB(b_.B)

	return nil
}

// NewC constructs a new instance of type C
func NewC(c string) *C {
	return &C{
		c: c,
	}
}

// With initializes this instance using the functional options pattern
func (c *C) With(opts ..._opt_[C, *C]) *C {
	c._def_()
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func COfC(c_ string) _opt_[C, *C] {
	return func(c *C) {
		c.c = c_
	}
}

func (c *C) _def_() {

}

// C gets the value of field c
func (c *C) C() string {
	return c.c
}

// SetC sets the value of field c
func (c *C) SetC(c_ string) {
	c.c = c_
}

type _C_marshal struct {
	C string `json:"c"`
}

type _C_unmarshal struct {
	C string `json:"c"`
}

// MarshalJSON serializes type C to json bytes
func (c C) MarshalJSON() ([]byte, error) {
	data := _C_marshal{

		C: c.C(),
	}
	return json.Marshal(data)
}

// UnmarshalJSON deserializes json bytes to type C
func (c *C) UnmarshalJSON(data []byte) error {
	var c_ _C_unmarshal
	if err := json.Unmarshal(data, &c_); err != nil {
		return nil
	}
	c.SetC(c_.C)

	return nil
}
