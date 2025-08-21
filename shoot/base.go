package shoot

import (
	"text/template"

	"github.com/lopolopen/shoot/shoot/transfer"
)

const Cmd = "shoot"
const FilePrefix = "0"

type BaseData struct {
	Cmd         string
	PackageName string
	TypeName    string
	transfers   template.FuncMap
}

func (d *BaseData) PreRegister() {
	d.Register("firstLower", transfer.FirstLower)
	d.Register("camelCase", transfer.ToCamelCase)
	d.Register("pascalCase", transfer.ToPascalCase)
}

func (d *BaseData) Register(key string, transfer any) {
	if d.transfers == nil {
		d.transfers = make(template.FuncMap)
	}

	d.transfers[key] = transfer
}

func (d *BaseData) Transfers() template.FuncMap {
	return d.transfers
}
