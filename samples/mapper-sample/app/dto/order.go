package dto

import (
	"mappersample/domain/enums"
	"mappersample/domain/model"
	"mappersample/infra/mapper"

	"github.com/shopspring/decimal"
)

//go:generate go tool shoot map -path=../../domain/model -type=Order

type Order struct {
	*mapper.Mapper
	ID           string            `json:"id"`
	Amount       string            `json:"amount"`
	Status       enums.OrderStatus `json:"status"`
	OrderingTime string            `json:"orderingTime" map:"OrderTime"`
	Address      OrderAddress      `json:"address"`
	Non          int
	Price        decimal.Decimal
	Int          int
}

// func (o Order) newModel() *model.Order {
// 	return new(model.Order)
// }

func (o Order) writeModel(dest *model.Order) {
	dest.ID = o.ID + "test"
	dest.X = 0 //suppress warnings
	// dest.Price = model.Price(o.Price)
}

func (o *Order) readModel(dest model.Order) {
	o.ID = dest.ID + "test"
	o.Non = 0 //suppress warnings or map:"-"
}

type OrderAddress struct {
	City    string `json:"city"`
	Street  string `json:"street"`
	RoomNum string `json:"roomNum" map:"Room"`
}
