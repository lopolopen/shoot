package dto

import (
	"mappersample/domain/enums"
	"mappersample/domain/model"
	"mappersample/infra/mapper"

	"github.com/shopspring/decimal"
)

type Order struct {
	*mapper.Mapper
	ID           string            `json:"id"`
	Amount       string            `json:"amount"`
	Status       enums.OrderStatus `json:"status"`
	OrderingTime string            `json:"orderingTime" map:"OrderTime"`
	Address      *OrderAddress     `json:"address"`
	Non          int               `map:"-"`
	Price        decimal.Decimal
	Price2       decimal.Decimal
	Value        int
}

func (o *Order) writeDomain(dest *model.Order) {
	dest.X = 0 //suppress warnings
}

type OrderAddress struct {
	City    string `json:"city"`
	Street  string `json:"street"`
	RoomNum string `json:"roomNum" map:"Room"`
}
