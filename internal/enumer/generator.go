package enumer

import (
	_ "embed"
	"flag"
	"go/ast"
	"go/token"
	"go/types"
	"log"

	"github.com/lopolopen/shoot/internal/shoot"
)

const SubCmd = "enum"

//go:embed enumer.tmpl
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
	bit := sub.Bool("bit", false, "generate bitwise enumerations (alias for -bitwise)")
	bitwise := sub.Bool("bitwise", false, "generate bitwise enumerations")
	json := sub.Bool("json", false, "generate MarshalJSON/UnmarshalJSON method for the type")
	text := sub.Bool("text", false, "generate MarshaText/UnmarshalText method for the type")

	g.ParseCommonFlags(sub)

	g.flags = &Flags{
		bitwise: *bitwise || *bit,
		json:    *json,
		text:    *text,
	}
}

func (g *Generator) Do(typeName string) bool {
	g.makeStr(typeName)
	g.makeBitwize()
	g.makeJson()
	g.makeText()

	if len(g.data.NameList) == 0 {
		return false
	}
	return true
}

func (g *Generator) TypeNames() []string {
	var typeNames []string
	pkg := g.Package()
	for _, f := range pkg.Files() {
		ast.Inspect(f.File(), func(n ast.Node) bool {
			decl, ok := n.(*ast.GenDecl)
			if !ok {
				return true
			}

			if decl.Tok == token.TYPE {
				for _, spec := range decl.Specs {
					ts, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}

					obj := pkg.Pkg().TypesInfo.Defs[ts.Name]
					if obj == nil {
						continue
					}

					typ := obj.Type()
					under := typ.Underlying()

					basic, ok := under.(*types.Basic)
					if !ok {
						continue
					}

					kind := basic.Kind()
					if kind != types.Int && kind != types.Uint &&
						kind != types.Int32 && kind != types.Uint32 {
						continue
					}

					if ts.Assign.IsValid() {
						log.Printf("[warn:] alias type %s will be ignored", ts.Name.Name)
					} else {
						typeNames = append(typeNames, ts.Name.Name)
					}
				}
			}

			return false
		})
	}
	return typeNames
}
