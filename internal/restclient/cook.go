package restclient

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"log"
	"regexp"
	"strings"

	"github.com/lopolopen/shoot/internal/shoot"
)

func (g *Generator) cookClient(typeName string) {
	var methodList []string
	var postList []string
	sigMap := make(map[string]string)
	httpMethodMap := make(map[string]string)
	pathMap := make(map[string]string)
	aliasMap := make(map[string]map[string]string)
	pathParansMap := make(map[string][]string)
	bodyParamMap := make(map[string]string)
	queryParamsMap := make(map[string][]string)
	resultTypeMap := make(map[string]string)

	for _, f := range g.pkg.files {
		ast.Inspect(f.file, func(n ast.Node) bool {
			ts, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}

			if ts.Name.Name != typeName {
				return true
			}

			iface, ok := ts.Type.(*ast.InterfaceType)
			if !ok {
				return true
			}

			var embedIfaces []string
			for _, field := range iface.Methods.List {
				if len(field.Names) == 0 {
					embedIfaces = append(embedIfaces, exprToString(field.Type))
				} else {
					ftype, ok := field.Type.(*ast.FuncType)
					if !ok {
						return true
					}

					doc := field.Doc.Text()
					methodName := field.Names[0].Name
					sigMap[methodName] = methodSignature(field) //full signature

					if field.Doc != nil {
						httpMethod, path, pathParams, ok := parsePath(doc) //pathParams ~ [id]
						if !ok {
							log.Println("[warn:]") //todo:
							continue
						}

						// switch httpMethod {
						// case http.MethodGet:
						// 	getList = append(getList, methodName)
						// case http.MethodPost:
						// 	postList = append(postList, methodName)
						// }

						methodList = append(methodList, methodName)
						httpMethodMap[methodName] = httpMethod //http mehod
						pathMap[methodName] = path             //http path

						alsMap := parseAlias(doc)            //userID -> id
						reversMap := make(map[string]string) //id -> userID
						for k, v := range alsMap {
							reversMap[v] = k
						}
						aliasMap[methodName] = alsMap

						var realPathParams []string
						for _, name := range pathParams {
							if real, ok := reversMap[name]; ok {
								realPathParams = append(realPathParams, real) //fix path param
							} else {
								realPathParams = append(realPathParams, name)
							}
						}
						pathParansMap[methodName] = realPathParams

						//------------Params---------------
						var queryParams []string
						for _, param := range ftype.Params.List {
							for _, name := range param.Names {
								switch t := param.Type.(type) {
								case *ast.SelectorExpr:
									typeName := exprToString(param.Type)
									if typeName == "context.Context" { //todo: ctx exists? import alias?
										continue
									}
								case *ast.StarExpr:
									if _, ok := t.X.(*ast.Ident); ok {
										bodyParamMap[methodName] = name.Name
									}
								case *ast.Ident:
									if isStructType(name.Name, f.file) { //todo:
										bodyParamMap[methodName] = name.Name
									} else {
										if shoot.Contains(realPathParams, name.Name) {
											continue
										}
										queryParams = append(queryParams, name.Name) //basic type
									}
								default:
									log.Fatalf("bad")
								}
							}
						}
						queryParamsMap[methodName] = queryParams

						if len(ftype.Results.List) > 2 {
							log.Fatalf("bad") //todo
						}
						//------------Results---------------
						for _, result := range ftype.Results.List {
							if len(result.Names) == 0 {
								typeName := getUnderlyingTypeName(result.Type)
								if typeName == "error" {
									continue
								}
								resultTypeMap[methodName] = typeName
							} else {
								//todo:
							}
						}

						//todo: check alias={a:alias}, a exists?
					}
				}
			}

			if !shoot.Contains(embedIfaces, "shoot.RestClient") {
				log.Printf("[warn:] interface without shoot.RestClient embed will be ignore")
				return true
			}

			return false
		})
	}

	g.data.MethodList = methodList
	g.data.PostList = postList
	g.data.SigMap = sigMap
	g.data.HTTPMethodMap = httpMethodMap
	g.data.PathMap = pathMap
	g.data.AliasMap = aliasMap
	g.data.PathParamsMap = pathParansMap
	g.data.QueryParamsMap = queryParamsMap
	g.data.ResultTypeMap = resultTypeMap
	g.data.BodyParamMap = bodyParamMap
}

func exprToString(expr ast.Expr) string {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, token.NewFileSet(), expr)
	if err != nil {
		log.Fatalf("print expr: %s", err)
	}
	return buf.String()
}

func methodSignature(method *ast.Field) string {
	funcName := method.Names[0].Name
	funcType, ok := method.Type.(*ast.FuncType)
	if !ok {
		return ""
	}

	params := formatFieldList(funcType.Params)
	results := formatFieldList(funcType.Results)

	if results != "" {
		return fmt.Sprintf("%s(%s) (%s)", funcName, params, results)
	}
	return fmt.Sprintf("%s(%s)", funcName, params)
}

func formatFieldList(fl *ast.FieldList) string {
	if fl == nil || len(fl.List) == 0 {
		return ""
	}
	var parts []string
	for _, f := range fl.List {
		typeStr := exprToString(f.Type)
		if len(f.Names) == 0 {
			parts = append(parts, typeStr)
		} else {
			for _, name := range f.Names {
				parts = append(parts, fmt.Sprintf("%s %s", name.Name, typeStr))
			}
		}
	}
	return strings.Join(parts, ", ")
}

func parsePath(doc string) (string, string, []string, bool) {
	regReq := regexp.MustCompile(`(?im)^shoot:\W+(get|post|put|patch|delete)\((.*)\)\W*;?\W*$`)
	ms := regReq.FindStringSubmatch(doc)
	if len(ms) == 0 {
		return "", "", nil, false
	}
	method := strings.ToUpper(ms[1])
	path := strings.TrimSpace(ms[2])

	regPathParam := regexp.MustCompile(`{(\w+)}`)
	psLst := regPathParam.FindAllStringSubmatch(path, -1)
	var pathParams []string
	for _, ps := range psLst {
		pathParams = append(pathParams, ps[1])
	}
	return method, path, pathParams, true
}

func parseAlias(doc string) map[string]string {
	regAlias := regexp.MustCompile(`(?m)^shoot:.*?\Walias=([^;\n]+)(;.*|\s*)$`)
	ms := regAlias.FindStringSubmatch(doc)
	if len(ms) == 0 {
		return nil
	}
	aliasLst := ms[1] //{userID:id},...
	regKV := regexp.MustCompile(`{(\w+)\W*:\W*(\w+)}`)
	asLst := regKV.FindAllStringSubmatch(aliasLst, -1)
	if len(asLst) == 0 {
		return nil
	}
	aliasMap := make(map[string]string)
	for _, as := range asLst {
		aliasMap[as[1]] = as[2]
	}
	return aliasMap //useID -> id
}

func getUnderlyingTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		// 指针类型，递归获取底层类型
		return getUnderlyingTypeName(t.X)
	case *ast.SelectorExpr:
		// 处理像 context.Context 这样的类型
		return getUnderlyingTypeName(t.X) + "." + t.Sel.Name
	default:
		return fmt.Sprintf("%T", expr)
	}
}

func isStructType(name string, file *ast.File) bool {
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.TYPE {
			continue
		}
		for _, spec := range gen.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || typeSpec.Name.Name != name {
				continue
			}
			_, isStruct := typeSpec.Type.(*ast.StructType)
			return isStruct
		}
	}
	return false
}
