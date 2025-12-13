package po

import (
	"mappersample/domain/enums"

	"github.com/shopspring/decimal"
)

//go:generate go tool shoot map -path=../../domain/model -type=Order

type Order struct {
	ID     string
	Amount decimal.Decimal   `gorm:"type:decimal(20,2)"`
	Status enums.OrderStatus `gorm:"type:enum('PENDING', 'COMPLETED', 'CANCELED')"`
	City   string
	Street string
	Room   string
	// X      []decimal.Decimal
	// Y      *string
}
