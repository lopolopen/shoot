package mapper

import (
	"go/ast"
	"go/types"
	"regexp"

	"github.com/lopolopen/shoot/internal/tools/logx"
	"golang.org/x/tools/go/packages"
)

func (g *Generator) parseSrcFields(srcTypeName string) {
	g.tagMap = make(map[string]string)

	ptrTypeMap := make(map[string]string)
	var exportedFields []Field
	for _, f := range g.Pkg().Syntax {
		ast.Inspect(f, func(n ast.Node) bool {
			if !g.testNode(srcTypeName, n) {
				return true
			}

			ts, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				return true
			}

			g.extractExportedTopFiels(g.Pkg(), st, ptrTypeMap, &exportedFields)
			return false
		})
	}

	g.ptrTypeMap = ptrTypeMap
	g.exportedFields = exportedFields
	for _, f := range exportedFields {
		g.data.SrcFieldList = append(g.data.SrcFieldList, f.name)
	}
}

func (g *Generator) parseDestFields(destTypeName string) bool {
	destExists := false

	ptrTypeMap := make(map[string]string)
	var exportedFields []Field
	for _, f := range g.destPkg.Syntax {
		ast.Inspect(f, func(n ast.Node) bool {
			if !g.testNode(destTypeName, n) {
				return true
			}

			destExists = true

			ts, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				return true
			}

			g.extractExportedTopFiels(g.destPkg, st, ptrTypeMap, &exportedFields)
			return false
		})
	}
	g.destExportedFields = exportedFields
	g.destPtrTypeMap = ptrTypeMap
	return destExists
}

func (g *Generator) makeMatch() {
	g.data.ExactMatchMap = make(map[string]string)
	g.data.ConvMatchMap = map[string]string{}
	g.data.SrcToDestTypeMap = make(map[string]string)
	g.data.DestToSrcTypeMap = make(map[string]string)

	for _, f1 := range g.exportedFields {
		if g.assignedSrcSet[f1.name] {
			continue
		}

		if _, ok := g.data.MismatchFuncMap[f1.name]; ok {
			continue
		}

		for _, f2 := range g.destExportedFields {
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
				//in ToXxx, type converter needs desc type
				g.data.SrcToDestTypeMap[f1.name] = qualifiedTypeName(f2.typ, g.flags.alias)
				//in FromXxx, type converter needs src type
				g.data.DestToSrcTypeMap[f2.name] = qualifiedTypeName(f1.typ, g.flags.alias)
			}
		}
	}
}

func qualifiedTypeName(t types.Type, alias string) string {
	qualifier := func(pkg *types.Package) string {
		if alias != "" {
			return alias
		}
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

func getMapTag(tag string) string {
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

func (g *Generator) extractExportedTopFiels(pkg *packages.Package, st *ast.StructType, ptrTypeMap map[string]string, fields *[]Field) {
	for _, f := range st.Fields.List {
		if len(f.Names) == 0 {
			//embedded: gorm.Model
			typ := pkg.TypesInfo.TypeOf(f.Type)
			name := typeName(typ)
			if !ast.IsExported(name) {
				continue
			}
			expandIfStruct(pkg, name, typ, ptrTypeMap, fields)
			continue
		}
		//named:
		name := f.Names[0]
		if f.Tag != nil {
			tag := getMapTag(f.Tag.Value)
			if tag == "-" {
				continue
			}
			if tag != "" {
				g.tagMap[name.Name] = tag
			}
		}
		for _, name := range f.Names {
			if !ast.IsExported(name.Name) {
				continue
			}

			if obj, ok := pkg.TypesInfo.Defs[name].(*types.Var); ok {
				*fields = append(*fields, Field{
					name: name.Name,
					path: name.Name,
					typ:  obj.Type(),
				})
			}
		}
	}
}

func expandIfStruct(pkg *packages.Package, pre string, t types.Type, ptrTypeMap map[string]string, fields *[]Field) {
	switch tt := t.(type) {
	case *types.Pointer:
		e := tt.Elem()
		if st, ok := e.Underlying().(*types.Struct); ok {
			ptrTypeMap[pre] = qualifiedTypeName(e, "")
			if n, ok := e.(*types.Named); ok {
				same := n.Obj().Pkg().Path() == pkg.PkgPath
				ptrTypeMap[pre] = qualifiedName(n.Obj().Pkg().Name(), "", n.Obj().Name(), same)
			}
			extractStructFields(pkg, pre, st, ptrTypeMap, fields)
		}
	case *types.Named:
		if st, ok := tt.Underlying().(*types.Struct); ok {
			extractStructFields(pkg, pre, st, ptrTypeMap, fields)
		}
	case *types.Struct: //todo: embeded struct
		logx.Pinln("!!!!!!!!!")
		extractStructFields(pkg, pre, tt, ptrTypeMap, fields)
	}
}

func extractStructFields(pkg *packages.Package, pre string, st *types.Struct, ptrSet map[string]string, fields *[]Field) {
	for i := 0; i < st.NumFields(); i++ {
		f := st.Field(i)
		if !ast.IsExported(f.Name()) {
			continue
		}
		if f.Embedded() {
			name := typeName(f.Type())
			expandIfStruct(pkg, pre+"."+name, f.Type(), ptrSet, fields)
			continue
		}
		*fields = append(*fields, Field{
			name: f.Name(),
			path: pre + "." + f.Name(),
			typ:  f.Type(),
		})
	}
}

func typeName(t types.Type) string {
	switch tt := t.(type) {
	case *types.Named:
		return tt.Obj().Name()
	case *types.Pointer:
		if named, ok := tt.Elem().(*types.Named); ok {
			return named.Obj().Name()
		}
	}
	return ""
}
