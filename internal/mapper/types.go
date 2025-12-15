package mapper

import (
	"go/types"

	"github.com/lopolopen/shoot/internal/shoot"
)

type TmplData struct {
	*shoot.TmplDataBase
	DestTypeName          string            //Order
	DestPkgName           string            //model
	DestPkgAlias          string            //domain
	DestPkgPath           string            //
	QualifiedDestTypeName string            //model.Order
	SrcFieldList          []string          //ID, Status ...
	ExactMatchMap         map[string]string //ID -> ID

	ConvMatchMap     map[string]string //Quantity -> Quantity
	SrcToDestTypeMap map[string]string //Quantity -> int32
	DestToSrcTypeMap map[string]string //Quantity -> int

	MismatchFuncMap  map[string]string //Amount -> Amount
	SrcToDestFuncMap map[string]string //Amount -> StringToDecimal
	DestToSrcFuncMap map[string]string //Amount -> DecimalToString

	MismatchSubMap map[string]string //Address -> Address
	SrcPtrSet      map[string]bool   //Address -> true
	DestPtrSet     map[string]bool   //Address -> false
	SrcSubTypeMap  map[string]string //Address -> OrderAddress

	ReadMethodName  string //fromModel
	ReadParamPrefix string
	WriteMethodName string //toModel
}

func NewTmplData(cmdline string) *TmplData {
	return &TmplData{
		TmplDataBase: shoot.NewTmplDataBase(cmdline),
	}
}

type Flags struct {
	destDir   string
	destTypes map[string]string
	alias     string
}

type Field struct {
	name string
	typ  types.Type
}

type Func struct {
	name   string
	param  types.Type
	result types.Type
}
