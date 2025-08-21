package constructor

import (
	"go/ast"
	"go/types"
	"regexp"
)

func (g *Generator) makeNew(typeName string) {
	var allList []string
	var unexportedList []string
	var newList []string
	typeMap := make(map[string]string)

	for _, f := range g.pkg.files {
		ast.Inspect(f.file, func(n ast.Node) bool {
			ts, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}
			if ts.Name.Name != typeName {
				return false
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				return true
			}

			for _, field := range st.Fields.List {
				for _, name := range field.Names {
					if name.Obj.Kind != ast.Var {
						continue
					}

					r := g.pkg.defs[name].(*types.Var)
					typeMap[name.Name] = r.Type().String()

					allList = append(allList, name.Name)

					if ast.IsExported(name.Name) {
						continue
					}

					unexportedList = append(unexportedList, name.Name)
					if field.Doc != nil {
						n := parseNew(field.Doc.Text())
						if n {
							newList = append(newList, name.Name)
						}
					}
				}
			}
			return false
		})
	}

	g.data.AllList = allList
	if len(newList) > 0 {
		g.data.NewList = newList
	} else {
		g.data.NewList = unexportedList
	}

	g.data.Register("typeof", func(key string) string {
		return typeMap[key]
	})
}

func parseNew(doc string) bool {
	regNew := regexp.MustCompile(`(?im)^shoot:.*?\Wnew(;.*|\s*)$`)
	new := regNew.Match([]byte(doc))
	return new
}
