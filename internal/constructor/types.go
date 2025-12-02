package constructor

import (
	"github.com/lopolopen/shoot/internal/shoot"
)

type Data struct {
	*shoot.BaseData
	// GoFile  string
	Imports string
	//All = Exported + Unexported
	AllList     []string
	NewList     []string
	GetSet      bool
	GetterList  []string
	SetterList  []string
	Option      bool
	DefaultList []string
	Json        bool
	//Marshal: Getteer + Exported
	//Unmarshal: Setter + Exported
	ExportedList []string
	EmbedList    []string
	Self         bool
	Short        bool
}

func NewData() *Data {
	return &Data{
		BaseData: shoot.NewBaseData(),
	}
}

type Flags struct {
	//if true:
	//[ ] = [get;set] => get+set
	//[get] => get-only
	//[set] => set-only
	//if false:
	//[ ] => neither
	//[get] => get-only
	//[set] => set-only
	//[get;set] => get+set
	getset bool
	json   bool
	opt    bool
	short  bool
}
