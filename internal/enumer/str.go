package enumer

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"log"
	"sort"
	"strings"
)

func (g *Generator) makeStr(typeName string) {
	var values []Value
	for _, f := range g.pkg.files {
		ast.Inspect(f.file, func(n ast.Node) bool {
			decl, ok := n.(*ast.GenDecl)
			if !ok || decl.Tok != token.CONST {
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
					obj, ok := f.pkg.defs[n]
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
					v := Value{
						originalName: n.Name,
						name:         n.Name,
						value:        u64,
						signed:       info&types.IsUnsigned == 0,
						str:          value.String(),
					}
					if c := vspec.Comment; f.lineComment && c != nil && len(c.List) == 1 {
						v.name = strings.TrimSpace(c.Text())
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
		return values[i].value < values[j].value
	})
	for _, v := range values {
		nameList = append(nameList, v.name)
		valueMap[v.name] = int64(v.value)
		shortName := strings.TrimPrefix(v.name, typeName)
		strMap[v.name] = shortName
	}

	g.data.NameList = nameList

	g.data.Register("valueof", func(key string) interface{} {
		return valueMap[key]
	})

	g.data.Register("strof", func(key string) interface{} {
		return strMap[key]
	})
}

// usize returns the number of bits of the smallest unsigned integer
// type that will hold n. Used to create the smallest possible slice of
// integers to use as indexes into the concatenated strings.
func usize(n int) int {
	switch {
	case n < 1<<8:
		return 8
	case n < 1<<16:
		return 16
	default:
		// 2^32 is enough constants for anyone.
		return 32
	}
}

// createIndexAndNameDecl returns the pair of declarations for the run.
func createIndexAndNameDecl(run []Value, typeName string, suffix string) (string, string) {
	if len(run) == 0 {
		log.Fatalln("empty run")
	}
	b := new(bytes.Buffer)
	indexes := make([]int, len(run))
	for i := range run {
		b.WriteString(run[i].name)
		indexes[i] = b.Len()
	}
	nameConst := fmt.Sprintf("_%sName%s = %q", typeName, suffix, b.String())
	nameLen := b.Len()
	b.Reset()
	_, _ = fmt.Fprintf(b, "_%sIndex%s = [...]uint%d{0, ", typeName, suffix, usize(nameLen))
	for i, v := range indexes {
		if i > 0 {
			_, _ = fmt.Fprintf(b, ", ")
		}
		_, _ = fmt.Fprintf(b, "%d", v)
	}
	_, _ = fmt.Fprintf(b, "}")
	return b.String(), nameConst
}
