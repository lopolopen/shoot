package model

import (
	"mappersample/domain/enums"

	"github.com/shopspring/decimal"
)

type Price decimal.Decimal
type Decimal = decimal.Decimal

type Order struct {
	// *Test
	ID string
	// Amount decimal.Decimal
	// Status    enums.OrderStatus
	// OrderTime time.Time
	// NullTime  *time.Time
	// Address   OrderAddress
	// X         int
	// Price     Price
	// Price2    Decimal
	// Value     int32
	X  *OrderAddress
	Xs []OrderAddress
}

func (o *Order) DomainMethod() {
	//domian work
}

type Test struct {
	Status enums.OrderStatus
}
