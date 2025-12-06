package enumer

import (
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"log"
	"sort"
	"strings"

	"github.com/lopolopen/shoot/internal/shoot"
)

func (g *Generator) makeStr(typeName string) {
	var values []shoot.Value
	for _, f := range g.Package().Files() {
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

					//todo:

					if ts.Assign.IsValid() {
						if typeName == ts.Name.Name {
							log.Fatalf("type %s should not be an alias", typeName)
						}
					}
				}
				return true
			}

			if decl.Tok != token.CONST {
				// We only care about const declarations.
				return true
			}
			// The name of the type of the constants we are declaring.
			// Can change if this is a multi-element declaration.
			typ := ""
			// Loop over the elements of the declaration. Each element is a ValueSpec:
			// a list of names possibly followed by a type, possibly followed by values.
			// If the type and value are both missing, we carry down the type (and value,
			// but the "go/types" package takes care of that).
			for _, spec := range decl.Specs {
				vspec := spec.(*ast.ValueSpec) // Guaranteed to succeed as this is CONST.
				if vspec.Type == nil && len(vspec.Values) > 0 {
					// "X = 1". With no type but a value, the constant is untyped.
					// Skip this vspec and reset the remembered type.
					typ = ""
					continue
				}
				if vspec.Type != nil {
					// "X T". We have a type. Remember it.
					ident, ok := vspec.Type.(*ast.Ident)
					if !ok {
						continue
					}
					typ = ident.Name
				}
				if typ != typeName {
					// This is not the type we're looking for.
					continue
				}
				// We now have a list of names (from one line of source code) all being
				// declared with the desired type.
				// Grab their names and actual values and store them in f.values.
				for _, n := range vspec.Names {
					if n.Name == "_" {
						continue
					}
					// This dance lets the type checker find the values for us. It's a
					// bit tricky: look up the object declared by the n, find its
					// types.Const, and extract its value.
					obj, ok := f.Pkg().Defs()[n]
					if !ok {
						log.Fatalf("no value for constant %s", n)
					}
					info := obj.Type().Underlying().(*types.Basic).Info()
					if info&types.IsInteger == 0 {
						log.Fatalf("can't handle non-integer constant type %s", typ)
					}
					value := obj.(*types.Const).Val() // Guaranteed to succeed as this is CONST.
					if value.Kind() != constant.Int {
						log.Fatalf("can't happen: constant is not an integer %s", n)
					}
					i64, isInt := constant.Int64Val(value)
					u64, isUint := constant.Uint64Val(value)
					if !isInt && !isUint {
						log.Fatalf("internal error: value of %s is not an integer: %s", n, value.String())
					}
					if !isInt {
						u64 = uint64(i64)
					}
					v := *shoot.NewValue(n.Name, n.Name, u64, info&types.IsUnsigned == 0, value.String())

					if c := vspec.Comment; f.LineComment() && c != nil && len(c.List) == 1 {
						v.SetName(strings.TrimSpace(c.Text()))
					}
					values = append(values, v)
				}
			}
			return false
		})
	}

	var nameList []string
	valueMap := make(map[string]int64)
	strMap := make(map[string]string)
	sort.Slice(values, func(i, j int) bool {
		return values[i].Value() < values[j].Value()
	})
	for _, v := range values {
		nameList = append(nameList, v.Name())
		valueMap[v.Name()] = int64(v.Value())
		shortName := strings.TrimPrefix(v.Name(), typeName)
		strMap[v.Name()] = shortName
	}

	g.data.NameList = nameList

	g.RegisterTransfer("valueof", func(key string) interface{} {
		return valueMap[key]
	})

	g.RegisterTransfer("strof", func(key string) interface{} {
		return strMap[key]
	})
}
