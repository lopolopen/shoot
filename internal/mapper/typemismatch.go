package mapper

import (
	"go/ast"
	"log"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

func (g *Generator) loadTypeMapperPkg(typeName string) string {
	var mapperTypeName string
	for _, f := range g.Pkg().Syntax {
		ast.Inspect(f, func(n ast.Node) bool {
			if !g.testNode(typeName, n) {
				return true
			}

			ts, _ := n.(*ast.TypeSpec)
			stru, _ := ts.Type.(*ast.StructType)
			for _, field := range stru.Fields.List {
				if len(field.Names) > 0 {
					continue
				}

				//may be ident or selector
				x := field.Type
				if star, ok := field.Type.(*ast.StarExpr); ok {
					x = star.X
				}

				_, ok := x.(*ast.Ident)
				if ok {
					g.mapperpkg = g.Pkg()
					return false
				}

				sel, ok := x.(*ast.SelectorExpr)
				if !ok {
					continue
				}
				imp := findImportForSelector(f, sel)
				if imp == nil {
					continue
				}
				impPath := strings.Trim(imp.Path.Value, `"`)
				cfg := &packages.Config{
					Mode: packages.NeedName |
						packages.NeedFiles |
						packages.NeedSyntax |
						packages.NeedTypes |
						packages.NeedTypesInfo,
				}
				pkgs, err := loadPkgs(cfg, impPath)
				if err != nil {
					log.Fatalf("❌ %s", err)
				}
				g.mapperpkg = pkgs[impPath]
				mapperTypeName = sel.Sel.Name
			}

			return false
		})
	}
	return mapperTypeName

}

func (g *Generator) parseMapper(mapperTypeName string) {
	if mapperTypeName == "" {
		return
	}

	var expFuncList []Func
	for _, f := range g.mapperpkg.Syntax {
		ast.Inspect(f, func(n ast.Node) bool {
			if !g.testNode(mapperTypeName, n) {
				return true
			}

			for _, decl := range f.Decls {
				if fn, ok := decl.(*ast.FuncDecl); ok && fn.Recv != nil {
					if len(fn.Recv.List) == 0 {
						continue
					}

					r := fn.Recv.List[0]
					switch expr := r.Type.(type) {
					case *ast.Ident:
						if expr.Name != mapperTypeName {
							continue
						}
					case *ast.StarExpr:
						if ident, ok := expr.X.(*ast.Ident); ok && ident.Name != mapperTypeName {
							continue
						}
					}

					if fn.Type.Params == nil {
						continue
					}
					params := fn.Type.Params.List
					if len(params) != 1 {
						continue
					}
					if fn.Type.Results == nil {
						continue
					}
					results := fn.Type.Results.List
					if len(results) != 1 {
						continue
					}

					expFuncList = append(expFuncList, Func{
						name:   fn.Name.Name,
						param:  g.mapperpkg.TypesInfo.TypeOf(params[0].Type),
						result: g.mapperpkg.TypesInfo.TypeOf(results[0].Type),
					})
				}
			}

			return false
		})
	}
	g.mappingFuncList = expFuncList
}

func (g *Generator) makeTypeMismatch() {
	g.data.SrcToDestFuncMap = make(map[string]string) //Fint -> IntToString
	g.data.DestToSrcFuncMap = make(map[string]string) //Fint -> StringToInt
	g.data.MismatchMap = map[string]string{}          //Fint -> FString

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

			if f1.typ.String() == f2.typ.String() {
				continue
			}

			//TODO:  !!!!!!!!!优先级？？？？

			for _, fn := range g.mappingFuncList {
				//in ToXxx, mapping func's param type is src filed type
				if fn.param.String() == f1.typ.String() && fn.result.String() == f2.typ.String() {
					g.data.SrcToDestFuncMap[f1.name] = fn.name
					g.data.MismatchMap[f1.name] = f2.name
				}

				//in FromXxx, the opposite applies
				if fn.param.String() == f2.typ.String() && fn.result.String() == f1.typ.String() {
					g.data.DestToSrcFuncMap[f1.name] = fn.name
					g.data.MismatchMap[f1.name] = f2.name
				}
			}
		}
	}
}

func findImportForSelector(file *ast.File, sel *ast.SelectorExpr) *ast.ImportSpec {
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return nil
	}

	name := ident.Name
	for _, imp := range file.Imports {
		if imp.Name != nil {
			if imp.Name.Name == name {
				return imp
			}
		} else {
			path := strings.Trim(imp.Path.Value, `"`)
			base := filepath.Base(path)
			if base == name {
				return imp
			}
		}
	}
	return nil
}
