package shoot

import (
	"text/template"

	"github.com/lopolopen/shoot/internal/transfer"
)

type Data interface {
	SetCmd(cmd string)

	SetTypeName(typeName string)

	SetPackageName(pkgName string)

	RegisterTransfer(key string, transfer any)

	Transfers() template.FuncMap
}

type BaseData struct {
	Cmd         string
	PackageName string
	TypeName    string
	transfers   template.FuncMap
}

func NewBaseData() *BaseData {
	d := &BaseData{}
	d.preRegister()
	return d
}

func (d *BaseData) preRegister() {
	d.RegisterTransfer("firstLower", transfer.FirstLower)
	d.RegisterTransfer("camelCase", transfer.ToCamelCase)
	d.RegisterTransfer("pascalCase", transfer.ToPascalCase)
	d.RegisterTransfer("in", func(s string, list []string) bool {
		for _, x := range list {
			if s == x {
				return true
			}
		}
		return false
	})
}

func (d *BaseData) SetCmd(cmd string) {
	d.Cmd = cmd
}

func (d *BaseData) SetTypeName(typeName string) {
	d.TypeName = typeName
}

func (d *BaseData) SetPackageName(pkgName string) {
	d.PackageName = pkgName
}

func (d *BaseData) RegisterTransfer(key string, transfer any) {
	if d.transfers == nil {
		d.transfers = make(template.FuncMap)
	}

	d.transfers[key] = transfer
}

func (d *BaseData) Transfers() template.FuncMap {
	return d.transfers
}
