package constructor

import (
	_ "embed"
	"flag"
	"go/ast"

	"github.com/lopolopen/shoot/internal/shoot"
)

const SubCmd = "new"

//go:embed constructor.tmpl
var tmplTxt string

// Generator holds the state of the analysis.
type Generator struct {
	*shoot.GenBase
	flags *Flags
	data  *Data
}

func New() *Generator {
	g := &Generator{
		GenBase: &shoot.GenBase{},
		data:    NewData(),
	}
	g.SetWorker(g)
	return g
}

func (g *Generator) SubCmd() string {
	return SubCmd
}

func (g *Generator) TmplTxt() string {
	return tmplTxt
}

func (g *Generator) Data() shoot.Data {
	return g.data
}

func (g *Generator) ParseFlags() {
	sub := flag.NewFlagSet(SubCmd, flag.ExitOnError)
	getset := sub.Bool("getset", false, "generate Get/Set method for the type")
	json := sub.Bool("json", false, "generate MarshalJSON/UnmarshalJSON method for the type")
	option := sub.Bool("option", false, "generate functional option pattern constructor")
	opt := sub.Bool("opt", false, "generate functional option pattern constructor (alias for -option)")
	short := sub.Bool("short", false, "generate short config function name (no 'OfType' suffix)")

	g.ParseCommonFlags(sub)

	g.flags = &Flags{
		getset: *getset,
		json:   *json,
		opt:    *opt || *option,
		short:  *short,
	}
}

func (g *Generator) Do(typeName string) bool {
	g.makeNew(typeName)
	g.makeOpt(typeName)
	g.makeGetSet(typeName)
	g.makeJson(typeName)

	return true
}

func (g *Generator) TypeNames() []string {
	var typeNames []string
	for _, f := range g.Package().Files() {
		ast.Inspect(f.File(), func(n ast.Node) bool {
			ts, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}

			_, ok = ts.Type.(*ast.StructType)
			if !ok {
				return true
			}

			typeNames = append(typeNames, ts.Name.Name)
			return false
		})
	}
	return typeNames
}
