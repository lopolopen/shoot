package shoot

import (
	"text/template"

	"github.com/lopolopen/shoot/internal/transfer"
)

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
	d.Register("in", func(s string, list []string) bool {
		for _, x := range list {
			if s == x {
				return true
			}
		}
		return false
	})
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
