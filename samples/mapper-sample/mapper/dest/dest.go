package dest

import (
	"time"
)

type Decimal struct{}

type Order struct {
	ID string
	// Amount    Decimal
	OrderTime time.Time
	// Address   OrderAddress
}

// type OrderAddress struct {
// 	City string
// }

// type Dest struct {
// 	ID       int
// 	DestName string
// 	X        string
// }
