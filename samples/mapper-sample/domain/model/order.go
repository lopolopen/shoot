package model

import (
	"mappersample/domain/enums"
	"time"

	"github.com/shopspring/decimal"
)

type Price decimal.Decimal
type Decimal = decimal.Decimal

type Order struct {
	ID        string
	Amount    decimal.Decimal
	Status    enums.OrderStatus
	OrderTime time.Time
	Address   OrderAddress
	X         int
	Price     Price
	Price2    Decimal
	Value     int32
}

func (o *Order) DomainMethod() {
	//domian work
}
