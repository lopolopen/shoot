package mapper

import (
	"time"

	"github.com/shopspring/decimal"
)

type Mapper struct{}

func (m *Mapper) StringToDecimal(s string) decimal.Decimal {
	return decimal.RequireFromString(s)
}

func (m *Mapper) DecimalToString(d decimal.Decimal) string {
	return d.String()
}

func (m *Mapper) StringToTime(s string) time.Time {
	t, err := time.ParseInLocation(time.DateTime, s, time.Local)
	if err != nil {
		panic(err)
	}
	return t
}

func (m *Mapper) TimeToString(t time.Time) string {
	return t.Format(time.DateTime)
}
