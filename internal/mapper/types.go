package mapper

import (
	"go/types"
	"strings"

	"github.com/lopolopen/shoot/internal/shoot"
)

type TmplData struct {
	*shoot.TmplDataBase
	DestTypeName          string            //Order
	DestPkgName           string            //model
	DestPkgAlias          string            //domain
	DestPkgPath           string            //mappersample/domain/model
	QualifiedDestTypeName string            //domain.Order
	SrcFieldList          []string          //[ID, Status, ...]
	ExactMatchMap         map[string]string //ID -> ID
	DestExactMatchMap     map[string]string

	SrcPtrTypeMap  map[string]string
	SrcPtrPathList []string

	SrcNeedReadCheckMap map[string]string

	DestNeedReadCheckMap map[string]string
	DestPtrTypeMap       map[string]string
	DestAccessCondMap    map[string]string

	ConvMatchMap     map[string]string //Quantity -> Quantity
	SrcToDestTypeMap map[string]string //Quantity -> int32
	DestToSrcTypeMap map[string]string //Quantity -> int

	MismatchFuncMap  map[string]string //Amount -> Amount
	SrcToDestFuncMap map[string]string //Amount -> StringToDecimal
	DestToSrcFuncMap map[string]string //Amount -> DecimalToString

	MismatchSubMap     map[string]string //Address -> Address
	DestMismatchSubMap map[string]string
	SrcPtrSet          map[string]bool   //Address -> true
	DestPtrSet         map[string]bool   //Address -> false
	SrcSubTypeMap      map[string]string //Address -> OrderAddress
	DestSubTypeMap     map[string]string //Address -> domain.OrderAddress

	MismatchSubListMap     map[string]string //AddrList -> AddrList
	DestMismatchSubListMap map[string]string

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
	name string //ID
	path string //Model.ID
	typ  types.Type
}

func (f Field) IsEmbeded() bool {
	return strings.Contains(f.path, dot)
}

func (f Field) CoveredBy(path string) bool {
	if f.path == path { //Model.ID
		return true
	}
	if strings.HasPrefix(f.path, path+dot) { //Model.
		return true
	}
	xs := strings.Split(path, dot)
	return strings.HasSuffix(f.path, dot+xs[len(xs)-1]) //.ID
}

type Func struct {
	name   string
	param  types.Type
	result types.Type
}
