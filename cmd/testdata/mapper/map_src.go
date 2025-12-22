package mapper

import "github.com/lopolopen/shoot/cmd/testdata/mapper/testmapper"

type Order struct {
	*testmapper.Mapper
	ID           string
	Amount       string
	OrderingTime string `map:"OrderTime"`
	Address      *OrderAddress
}

type OrderAddress struct {
	City string
}

type Src struct {
	*Embed
	ID      int
	SrcName string `map:"DestName"`
}

type Embed struct {
	ID int
	*EmbedEmbed
}

type EmbedEmbed struct {
	X string
}
