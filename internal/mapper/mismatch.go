package mapper

import (
	"go/ast"
	"go/types"
	"path/filepath"
	"strings"

	"github.com/lopolopen/shoot/internal/tools/logx"
	"golang.org/x/tools/go/packages"
)

func (g *Generator) loadTypeMapperPkg(typeName string) string {
	//todo: recursive
	var mappers string
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

				if field.Tag != nil {
					tag := mapTag(field.Tag.Value)
					if tag == "-" {
						continue
					}
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
					logx.Fatalf("%s", err)
				}
				g.mapperpkg = pkgs[impPath]
				mappers = sel.Sel.Name
			}
			return false
		})
	}
	return mappers
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

func (g *Generator) makeMismatch() {
	g.data.SrcToDestFuncMap = make(map[string]string)
	g.data.DestToSrcFuncMap = make(map[string]string)
	g.data.MismatchFuncMap = make(map[string]string)
	g.data.MismatchSubMap = make(map[string]string)
	g.data.SrcPtrSet = make(map[string]bool)
	g.data.DestPtrSet = make(map[string]bool)
	g.data.SrcSubTypeMap = make(map[string]string)

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

			g.makeFuncMap(f1, f2)
			g.makeSubMap(f1, f2)
		}
	}
}

func (g *Generator) makeFuncMap(f1, f2 Field) {
	for _, fn := range g.mappingFuncList {
		//in ToXxx, mapping func's param type is src filed type
		if fn.param.String() == f1.typ.String() && fn.result.String() == f2.typ.String() {
			g.data.SrcToDestFuncMap[f1.name] = fn.name
			g.data.MismatchFuncMap[f1.name] = f2.name
		}

		//in FromXxx, the opposite applies
		if fn.param.String() == f2.typ.String() && fn.result.String() == f1.typ.String() {
			g.data.DestToSrcFuncMap[f1.name] = fn.name
			g.data.MismatchFuncMap[f1.name] = f2.name
		}
	}
}

func (g *Generator) makeSubMap(sub1, sub2 Field) {
	if _, ok := g.data.MismatchFuncMap[sub1.name]; ok {
		return
	}

	typ1 := sub1.typ
	typ2 := sub2.typ
	if p1, ok := typ1.(*types.Pointer); ok {
		g.data.SrcPtrSet[sub1.name] = true
		typ1 = p1.Elem()
	}
	if p2, ok := typ2.(*types.Pointer); ok {
		g.data.DestPtrSet[sub2.name] = true
		typ2 = p2.Elem()
	}

	if n1, ok := typ1.(*types.Named); ok {
		pkgpath1 := n1.Obj().Pkg().Path()
		g.data.SrcSubTypeMap[sub1.name] = n1.Obj().Name()

		if named2, ok := typ2.(*types.Named); ok {
			pkgpath2 := named2.Obj().Pkg().Path()
			if pkgpath1 == g.Pkg().PkgPath && pkgpath2 == g.destpkg.PkgPath {
				g.data.MismatchSubMap[sub1.name] = sub2.name
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
