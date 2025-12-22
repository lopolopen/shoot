package testmapper

import (
	"time"

	"github.com/lopolopen/shoot/cmd/testdata/mapper/dest"
)

type Mapper struct {
	//shoot: map
}

func (Mapper) StringToDecimal(s string) dest.Decimal {
	return dest.Decimal{}
}

func (Mapper) DecimalToString(d dest.Decimal) string {
	return ""
}

func (Mapper) StringToTime(s string) time.Time {
	t, err := time.ParseInLocation(time.DateTime, s, time.Local)
	if err != nil {
		panic(err)
	}
	return t
}

func (Mapper) TimeToString(t time.Time) string {
	return t.Format(time.DateTime)
}
