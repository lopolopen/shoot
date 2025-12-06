package enumer

import (
	"github.com/lopolopen/shoot/internal/shoot"
)

type TmplData struct {
	shoot.TmplDataBase
	NameList []string
	Bitwise  bool
	Json     bool
	Text     bool
}

func NewTmplData() *TmplData {
	return &TmplData{
		TmplDataBase: *shoot.NewTmplDataBase(),
	}
}

type Flags struct {
	bitwise bool
	json    bool
	text    bool
}
