package restclient

import (
	_ "embed"
	"flag"
	"go/ast"
	"go/types"

	"github.com/lopolopen/shoot/internal/shoot"
)

const SubCmd = "rest"

//go:embed restclient.tmpl
var tmplTxt string

type Generator struct {
	*shoot.GenBase
	// flags *Flags
	data *Data
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

	g.ParseCommonFlags(sub)
}

func (g *Generator) Do(typeName string) bool {
	g.cookClient(typeName)

	return true
}

func (g *Generator) TypeNames() []string {
	var typeNames []string
	for _, f := range g.Package().Files() {
		ast.Inspect(f.File(), func(n ast.Node) bool {
			if !g.testNode("", n) {
				return true
			}

			ts, _ := n.(*ast.TypeSpec)
			typeNames = append(typeNames, ts.Name.Name)
			return false
		})
	}
	return typeNames
}

func (g *Generator) testNode(typeName string, node ast.Node) bool {
	ts, ok := node.(*ast.TypeSpec)
	if !ok {
		return false
	}

	if typeName != "" && ts.Name.Name != typeName {
		return false
	}

	iface, ok := ts.Type.(*ast.InterfaceType)
	if !ok {
		return false
	}

	isRestClient := false
	for _, field := range iface.Methods.List {
		if len(field.Names) > 0 {
			continue
		}

		typ := g.Package().Pkg().TypesInfo.Types[field.Type].Type
		named, ok := typ.(*types.Named)
		if !ok {
			continue
		}
		obj := named.Obj()
		pkgPath := obj.Pkg().Path()
		if pkgPath == shoot.SelfPkgPath && obj.Name() == "RestClient" {
			isRestClient = true
			break
		}
	}
	return isRestClient
}
