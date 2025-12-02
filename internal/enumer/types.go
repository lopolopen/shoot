package enumer

import (
	"github.com/lopolopen/shoot/internal/shoot"
)

type Data struct {
	shoot.BaseData
	NameList []string
	Bitwise  bool
	Json     bool
	Text     bool
}

func NewData() *Data {
	return &Data{
		BaseData: *shoot.NewBaseData(),
	}
}

type Flags struct {
	bitwise bool
	json    bool
	text    bool
}
