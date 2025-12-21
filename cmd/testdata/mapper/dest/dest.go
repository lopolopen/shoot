package dest

import (
	"time"

	"github.com/shopspring/decimal"
)

type Order struct {
	ID        string
	Amount    decimal.Decimal
	OrderTime time.Time
	Address   OrderAddress
}
