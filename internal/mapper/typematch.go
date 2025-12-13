package mapper

import (
	"go/ast"
	"go/types"
	"regexp"
)

func (g *Generator) parseSrcFields(srcTypeName string) {
	g.tagMap = make(map[string]string)

	var expList []Field
	for _, f := range g.Pkg().Syntax {
		ast.Inspect(f, func(n ast.Node) bool {
			if !g.testNode(srcTypeName, n) {
				return true
			}

			ts, _ := n.(*ast.TypeSpec)
			stru, _ := ts.Type.(*ast.StructType)
			for _, field := range stru.Fields.List {
				if len(field.Names) == 0 {
					continue
				}

				name := field.Names[0].Name
				if !ast.IsExported(name) {
					continue
				}

				//`map:"DestFieldName"`
				if field.Tag != nil {
					tag := mapTag(field.Tag.Value)
					if tag == "-" {
						continue
					}
					if tag != "" {
						g.tagMap[name] = tag
					}
				}

				expList = append(expList, Field{
					name: name,
					typ:  g.Pkg().TypesInfo.TypeOf(field.Type),
				})
			}
			return false
		})
	}

	g.srcExpList = expList
}

func (g *Generator) parseDestFields(destTypeName string) {
	var destExpList []Field
	for _, f := range g.destpkg.Syntax {
		ast.Inspect(f, func(n ast.Node) bool {
			if !g.testNode(destTypeName, n) {
				return true
			}

			ts, _ := n.(*ast.TypeSpec)
			stru, _ := ts.Type.(*ast.StructType)
			for _, field := range stru.Fields.List {
				if len(field.Names) == 0 {
					//embedded field
				} else {
					name := field.Names[0].Name

					if !ast.IsExported(name) {
						continue
					}

					destExpList = append(destExpList, Field{
						name: name,
						typ:  g.destpkg.TypesInfo.TypeOf(field.Type),
					})
				}
			}
			return false
		})
	}
	g.destExpList = destExpList
}

func (g *Generator) makeMatch() {
	g.data.ExactMatchMap = make(map[string]string)

	g.data.ConvMatchMap = map[string]string{}
	g.data.SrcToDestTypeMap = make(map[string]string)
	g.data.DestToSrcTypeMap = make(map[string]string)

	for _, f1 := range g.srcExpList {
		if g.assignedSrcSet[f1.name] {
			continue
		}

		for _, f2 := range g.destExpList {
			if g.assignedDestSet[f2.name] {
				continue
			}

			if !canNameMatch(f1.name, f2.name, g.tagMap) {
				continue
			}

			same, conv := matchType(f1.typ, f2.typ)
			if same {
				g.data.ExactMatchMap[f1.name] = f2.name
			} else if conv {
				g.data.ConvMatchMap[f1.name] = f2.name
				//in ToXxx, convert's type is desc type
				g.data.SrcToDestTypeMap[f1.name] = qualifiedName(f2.typ)
				//in FromXxx, the opposite applies
				g.data.DestToSrcTypeMap[f2.name] = qualifiedName(f1.typ)
			}
		}
	}
}

func qualifiedName(t types.Type) string {
	qualifier := func(pkg *types.Package) string {
		if pkg == nil {
			return ""
		}
		return pkg.Name()
	}
	return types.TypeString(t, qualifier)
}

func canNameMatch(name1, name2 string, tagMap map[string]string) bool {
	if tagMap == nil {
		tagMap = make(map[string]string)
	}

	n, ok := tagMap[name1]
	if !ok {
		n = name1
	}

	return n == name2
}

func mapTag(tag string) string {
	reg := regexp.MustCompile(`map:"([^"]*)"`)
	matches := reg.FindStringSubmatch(tag)
	if len(matches) <= 1 {
		return ""
	}
	return matches[1]
}

func matchType(type1, type2 types.Type) (bool, bool) {
	same := types.Identical(type1, type2)
	conv := types.ConvertibleTo(type1, type2)
	return same, conv
}
