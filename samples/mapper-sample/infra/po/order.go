package po

import (
	"database/sql"
	"mappersample/domain/enums"
	"mappersample/domain/model"
	"mappersample/infra/mapper"
	"time"

	"github.com/shopspring/decimal"
)

//-go:generate go tool shoot map -path=../../domain/model -type=Order

type Order struct {
	*mapper.SQLMapper
	ID        string
	Amount    decimal.Decimal   `gorm:"type:decimal(20,2)"`
	Status    enums.OrderStatus `gorm:"type:enum('PENDING', 'COMPLETED', 'CANCELED')"`
	OrderTime time.Time
	NullTime  sql.Null[time.Time]
	City      string `map:"Addres."`
	Street    string `map:"Addres."`
	Room      string `map:"Addres."` //todo:
}

func (o *Order) writeModel(order *model.Order) {
	order.X = 0
	order.Price = model.Price{}
	order.Price2 = decimal.Zero
	order.Value = 0
}
