package constructor

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"log"
	"regexp"
)

func (g *Generator) makeNew(typeName string) {
	var imports string
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

			imports = buildImports(f.file.Imports)

			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				return true
			}

			for _, field := range st.Fields.List {
				for _, name := range field.Names {
					if name.Obj.Kind != ast.Var {
						continue
					}

					// r := g.pkg.defs[name].(*types.Var)
					// if r != nil {
					// 	typeMap[name.Name] = stripPkgPrefix(r.Type())
					// }

					fs := token.NewFileSet()
					typeMap[name.Name] = exprString(fs, field.Type)

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

	g.data.Imports = imports

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

func buildImports(imports []*ast.ImportSpec) string {
	var buff bytes.Buffer
	for _, imp := range imports {
		if imp.Name != nil {
			buff.WriteString(imp.Name.Name)
			buff.WriteString(" ")
		}
		if imp.Path != nil {
			buff.WriteString(imp.Path.Value)
			buff.WriteString("\n")
		}
	}
	return buff.String()
}

func parseNew(doc string) bool {
	regNew := regexp.MustCompile(`(?im)^shoot:.*?\Wnew(;.*|\s*)$`)
	new := regNew.Match([]byte(doc))
	return new
}

func exprString(fset *token.FileSet, expr ast.Expr) string {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, fset, expr)
	if err != nil {
		log.Fatalf("print expr: %s", err)
	}
	return buf.String()
}

// func stripPkgPrefix(t types.Type) string {
// 	switch tt := t.(type) {
// 	case *types.Named:
// 		return tt.Obj().Name()
// 	case *types.Pointer:
// 		return "*" + stripPkgPrefix(tt.Elem())
// 	case *types.Slice:
// 		return "[]" + stripPkgPrefix(tt.Elem())
// 	case *types.Array:
// 		return fmt.Sprintf("[%d]%s", tt.Len(), stripPkgPrefix(tt.Elem()))
// 	case *types.Map:
// 		return fmt.Sprintf("map[%s]%s",
// 			stripPkgPrefix(tt.Key()), stripPkgPrefix(tt.Elem()))
// 	case *types.Chan:
// 		dir := ""
// 		if tt.Dir() == types.SendOnly {
// 			dir = "chan<- "
// 		} else if tt.Dir() == types.RecvOnly {
// 			dir = "<-chan "
// 		} else {
// 			dir = "chan "
// 		}
// 		return dir + stripPkgPrefix(tt.Elem())
// 	case *types.Basic:
// 		return tt.Name()
// 	default:
// 		return t.String()
// 	}
// }
