package constructor

import (
	"go/ast"
	"regexp"

	"github.com/lopolopen/shoot/internal/transfer"
)

func (g *Generator) makeJson(typeName string) {
	g.data.RegisterTransfer("jsontagof", transfer.ID)

	if !g.flags.json {
		return
	}

	var exportedList []string
	tagMap := make(map[string]string)

	for _, f := range g.Package().Files() {
		ast.Inspect(f.File(), func(n ast.Node) bool {
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

					var tag string
					if field.Tag != nil {
						tag = jsonTag(field.Tag.Value)
						tagMap[name.Name] = tag
					} else {
						tagMap[name.Name] = name.Name
					}

					if ast.IsExported(name.Name) {
						exportedList = append(exportedList, name.Name)
					}
				}
			}
			return false
		})
	}
	g.data.ExportedList = exportedList
	g.data.RegisterTransfer("jsontagof", func(key string) string {
		return tagMap[key]
	})
	g.data.Json = true
}

func jsonTag(tag string) string {
	reg := regexp.MustCompile(`json:"([^"]*)"`)
	matches := reg.FindStringSubmatch(tag)
	if len(matches) <= 1 {
		return ""
	}
	return matches[1]
}
