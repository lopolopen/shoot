package mapper

import (
	"fmt"
	"go/ast"
	"log"
	"strings"

	"github.com/lopolopen/shoot/internal/transfer"
)

func (g *Generator) parseCustomMapMethods(srcTypeName, destTypeName string) []string {
	g.assignedSrcSet = make(map[string]bool)
	g.assignedDestSet = make(map[string]bool)

	pkg := g.Pkg()
	for _, f := range pkg.Syntax {
		ast.Inspect(f, func(n ast.Node) bool {
			if !g.testNode(srcTypeName, n) {
				return true
			}

			for _, decl := range f.Decls {
				if fn, ok := decl.(*ast.FuncDecl); ok && fn.Recv != nil {
					if len(fn.Recv.List) == 0 {
						continue
					}

					isWrite := isWriteMethod(fn.Name.Name, g.destpkg.Name)
					isRead := !isWrite && isReadMethod(fn.Name.Name, g.destpkg.Name)

					var recvTypeName string
					recv := fn.Recv.List[0]
					switch expr := recv.Type.(type) {
					case *ast.Ident: //value receiver
						if expr.Name != srcTypeName {
							continue
						}
						recvTypeName = srcTypeName
						if isRead {
							log.Fatalf("❌ method (%s).%s must use a pointer receiver", recvTypeName, fn.Name.Name)
						}
					case *ast.StarExpr: //pointer receiver
						if ident, ok := expr.X.(*ast.Ident); ok && ident.Name != srcTypeName {
							continue
						}
						recvTypeName = "*" + srcTypeName
						if isWrite {
							log.Printf("💡 method (%s).%s should use a value receiver (recommended)", recvTypeName, fn.Name.Name)
						}
					}

					if fn.Type.Params == nil {
						continue
					}
					params := fn.Type.Params.List
					if len(params) != 1 {
						continue
					}
					if fn.Type.Results != nil && len(fn.Type.Results.List) != 0 {
						//warn
						continue
					}

					paramTypeExpr := params[0].Type
					switch expr := paramTypeExpr.(type) {
					case *ast.Ident:
						//todo:
					case *ast.SelectorExpr: //value dest param
						paramTypeExpr = expr.Sel
						if isWrite {
							log.Fatalf("❌ method (%s).%s must has a pointer parameter", recvTypeName, fn.Name.Name)
						} else {
							g.data.ReadParamPrefix = "*"
						}
					case *ast.StarExpr: //pointer dest param
						paramTypeExpr = expr.X
						if isRead {
							log.Printf("💡 method (%s).%s should use a value parameter (recommended)", recvTypeName, fn.Name.Name)
						}
					}

					//checking param type needs full type name
					destFullName := fmt.Sprintf("%s.%s", g.destpkg.PkgPath, destTypeName)
					paramType := pkg.TypesInfo.TypeOf(paramTypeExpr)
					if paramType.String() != destFullName {
						continue
					}

					if isWrite {
						if g.data.WriteMethodName == "" {
							g.data.WriteMethodName = fn.Name.Name
						} else {
							log.Fatalf("❌ found more than one custom write method: (%s).%s", recvTypeName, fn.Name.Name)
						}
						names := findAssignedFields(fn, params[0].Names[0].Name)
						for _, n := range names {
							g.assignedDestSet[n] = true
						}
					} else if isRead {
						if g.data.ReadMethodName == "" {
							g.data.ReadMethodName = fn.Name.Name
						} else {
							log.Fatalf("❌ found more than one custom read method: (%s).%s", recvTypeName, fn.Name.Name)
						}
						names := findAssignedFields(fn, recv.Names[0].Name)
						for _, n := range names {
							g.assignedSrcSet[n] = true
						}
					}

				}
			}

			return false
		})
	}
	return nil
}

func isWriteMethod(methodName string, destPkgName string) bool {
	keys := []string{"to", "write"}
	for _, key := range keys {
		if methodName == key+strings.ToLower(destPkgName) {
			return true
		}

		if methodName == key+transfer.ToPascalCase(destPkgName) {
			return true
		}
	}
	return false
}

func isReadMethod(methodName string, destPkgName string) bool {
	keys := []string{"from", "read"}
	for _, key := range keys {
		if methodName == key+strings.ToLower(destPkgName) {
			return true
		}

		if methodName == key+transfer.ToPascalCase(destPkgName) {
			return true
		}
	}
	return false
}

func findAssignedFields(funcDecl *ast.FuncDecl, name string) []string {
	//todo: recursive
	var fields []string
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}
		for _, lhs := range assign.Lhs {
			if sel, ok := lhs.(*ast.SelectorExpr); ok {
				if ident, ok := sel.X.(*ast.Ident); ok && ident.Name == name {
					fields = append(fields, sel.Sel.Name)
				}
			}
		}
		return true
	})
	return fields
}
