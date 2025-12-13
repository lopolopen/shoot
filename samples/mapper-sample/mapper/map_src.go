package mapper

//go:generate go run tool map -path=./dest -type=Order
type Order struct {
	// *testmapper.Mapper
	ID           string
	Amount       string
	OrderingTime string `map:"OrderTime"`
	// Address      *OrderAddress
}

// type OrderAddress struct {
// 	City string
// }

// type Src struct {
// 	*Embed
// 	SrcName string `map:"DestName"`
// }

// type Embed struct {
// 	ID int
// 	*EmbedEmbed
// }

// type EmbedEmbed struct {
// 	X string
// }
