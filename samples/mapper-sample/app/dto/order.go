package dto

import (
	"mappersample/domain/model"
	"mappersample/infra/mapper"
)

// type inMapper struct {
// 	//shoot: map
// }

// func (inMapper) StringToDecimal(s string) decimal.Decimal {
// 	return decimal.RequireFromString(s)
// }

//go:generate go tool shoot map -path=../../domain/model -alias=domain -type=Order

type Order struct {
	// inMapper
	mapper.Mapper

	ID string `json:"id"`
	// Amount string `json:"amount"`
	// Status       enums.OrderStatus `json:"status"`
	// OrderingTime string            `json:"orderingTime" map:"OrderTime"`
	// Address      *OrderAddress     `json:"address"`
	// Non          int               `map:"-"`
	// Price        decimal.Decimal
	// Price2       decimal.Decimal
	// Value        int
	X  *OrderAddress
	Xs []OrderAddress
}

func (o *Order) writeDomain(dest *model.Order) {
	// dest.X = 0 //suppress warnings
	// dest.Amount = decimal.Zero
}

func (o *Order) readDomain(dest *model.Order) {
	// o.Amount = ""
	o.X = nil
	o.Xs = nil
	// o.ID = ""
}

// func (*Order) StringToDecimal(s string) decimal.Decimal {
// 	return decimal.RequireFromString(s)
// }

type OrderAddress struct {
	City    string `json:"city"`
	Street  string `json:"street"`
	RoomNum string `json:"roomNum" map:"Room"`
}
