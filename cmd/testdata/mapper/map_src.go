package mapper

import "github.com/lopolopen/shoot/cmd/testdata/mapper/mapper"

type Order struct {
	*mapper.Mapper
	ID           string        `json:"id"`
	Amount       string        `json:"amount"`
	OrderingTime string        `json:"orderingTime" map:"OrderTime"`
	Address      *OrderAddress `json:"address"`
}
