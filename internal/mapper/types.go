package mapper

import (
	"go/types"

	"github.com/lopolopen/shoot/internal/shoot"
)

type TmplData struct {
	*shoot.TmplDataBase
	DestTypeName          string
	DestPkgName           string
	QualifiedDestTypeName string
	SrcFieldList          []string
	ExactMatchMap         map[string]string
	ConvMatchMap          map[string]string
	SrcToDestTypeMap      map[string]string
	DestToSrcTypeMap      map[string]string
	MismatchMap           map[string]string
	SrcToDestFuncMap      map[string]string
	DestToSrcFuncMap      map[string]string
	ReadMethodName        string
	ReadParamPrefix       string
	WriteMethodName       string
}

func NewTmplData(cmdline string) *TmplData {
	return &TmplData{
		TmplDataBase: shoot.NewTmplDataBase(cmdline),
	}
}

type Flags struct {
	destDir   string
	destTypes map[string]string
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
