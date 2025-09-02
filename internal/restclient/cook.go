package restclient

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/printer"
	"go/token"
	"go/types"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/lopolopen/shoot/internal/shoot"
	"github.com/lopolopen/shoot/internal/transfer"
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
	queryDictMap := make(map[string]string)
	queryParamsMap := make(map[string][]string)
	resultTypeMap := make(map[string]string)
	ctxParamMap := make(map[string]string)
	defHeaders := map[string]map[string]string{
		http.MethodGet: {
			"Accept": "application/json",
		},
		http.MethodPost: {
			"Accept":       "application/json",
			"Content-Type": "application/json",
		},
		http.MethodPut: {
			"Accept":       "application/json",
			"Content-Type": "application/json",
		},
		http.MethodPatch: {
			"Accept":       "application/json",
			"Content-Type": "application/json",
		},
		http.MethodDelete: {},
	}

	for _, f := range g.pkg.files {
		ast.Inspect(f.file, func(n ast.Node) bool {
			if !g.testNode(typeName, n) {
				return true
			}

			ts, _ := n.(*ast.TypeSpec)
			iface, _ := ts.Type.(*ast.InterfaceType)
			for _, field := range iface.Methods.List {
				if len(field.Names) == 0 {
					if field.Doc != nil {
						headers := parseHeaders(field.Doc.Text())
						for k, v := range headers {
							for _, headers := range defHeaders {
								headers[k] = v
							}
						}
					}
				} else {
					ftype, ok := field.Type.(*ast.FuncType)
					if !ok {
						continue
					}

					doc := field.Doc.Text()
					methodName := field.Names[0].Name
					sigMap[methodName] = methodSignature(field) //full signature

					if field.Doc == nil {
						log.Printf("[warn:] method %s without comments will be ignored", methodName)
						continue
					}

					httpMethod, path, pathParams, ok := parsePath(doc) //pathParams ~ [id]
					if !ok {
						log.Printf("[warn:] method %s with bad comments will be ignored", methodName)
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

					asMap := parseAlias(doc)             //userID -> id
					reversMap := make(map[string]string) //id -> userID
					for k, v := range asMap {
						reversMap[v] = k
					}
					aliasMap[methodName] = asMap

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
					if ftype.Params != nil {
						for _, param := range ftype.Params.List {
							for _, name := range param.Names {
								switch t := param.Type.(type) {
								case *ast.SelectorExpr:
									typ := g.pkg.pkg.TypesInfo.Types[param.Type].Type
									named, ok := typ.(*types.Named)
									if !ok {
										continue
									}
									obj := named.Obj()
									pkgPath := obj.Pkg().Path()
									if pkgPath == "context" && obj.Name() == "Context" {
										ctxParamMap[methodName] = name.Name
									} else {
										bodyParamMap[methodName] = name.Name
										fullPath, err := getPkgDir(pkgPath)
										if err != nil {
											log.Fatalf("get pkg dir: %s", err)
										}
										fields, err := extractStructFields(fullPath, t.Sel.Name)
										if err != nil {
											log.Fatalf("extract struct fields: %s", err)
										}
										for _, f := range fields {
											key := f.Name //name
											value := f.Name
											if f.IsExported {
												key = transfer.ToCamelCase(f.Name)
												value = fmt.Sprintf("%s.%s", name.Name, f.Name)
											} else {
												value = fmt.Sprintf("%s.%s()", name.Name, transfer.ToPascalCase(f.Name))
											}
											queryParams = append(queryParams, value)

											if aliasMap[methodName] == nil {
												aliasMap[methodName] = make(map[string]string)
											}
											if f.Alias != "" {
												aliasMap[methodName][value] = f.Alias
											} else {
												aliasMap[methodName][value] = key
											}
										}
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
								case *ast.MapType:
									if httpMethod == http.MethodGet {
										queryDictMap[methodName] = name.Name
									}
								default:
									log.Fatalf("bad")
								}
							}
						}
					}
					queryParamsMap[methodName] = queryParams

					hasErr := false
					if ftype.Results != nil {
						if len(ftype.Results.List) > 2 {
							log.Fatalf("method %s must not return more than two values", methodName)
						}

						//------------Results---------------
						for _, result := range ftype.Results.List {
							if len(result.Names) == 0 {
								typeName := getUnderlyingTypeName(result.Type)
								if typeName == "error" {
									hasErr = true
									continue
								}
								resultTypeMap[methodName] = typeName
							} else {
								//todo:
							}
						}
					}
					if !hasErr {
						log.Fatalf("method %s must return an error", methodName)
					}
					//todo: check alias={a:alias}, a exists?
				}
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
	g.data.QueryDictMap = queryDictMap
	g.data.QueryParamsMap = queryParamsMap
	g.data.DefaultHeaders = defHeaders
	g.data.CtxParamMap = ctxParamMap
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

func parseHeaders(doc string) map[string]string {
	headers := make(map[string]string)
	regHeaders := regexp.MustCompile(`shoot:.*?\Wheaders=((?:\s*{[^\n]+},?)+)`)
	ms := regHeaders.FindStringSubmatch(doc)
	if len(ms) > 0 {
		kvMap := parseKV(ms[1])
		for k, v := range kvMap {
			headers[k] = v
		}
	}
	return headers
}

func parseKV(str string) map[string]string {
	regKV := regexp.MustCompile(`{([\w|-]+)\W*:\W*([^}]+)}`)
	kvLst := regKV.FindAllStringSubmatch(str, -1)
	if len(kvLst) == 0 {
		return nil
	}
	kvMap := make(map[string]string)
	for _, kv := range kvLst {
		kvMap[kv[1]] = kv[2]
	}
	return kvMap
}

func parsePath(doc string) (string, string, []string, bool) {
	regReq := regexp.MustCompile(`(?im)^shoot:\W+(get|post|put|patch|delete)\((.*)\)\W*;?\W*$`)
	ms := regReq.FindStringSubmatch(doc)
	if len(ms) == 0 {
		return "", "", nil, false
	}
	method := strings.ToUpper(ms[1])
	path := strings.TrimSpace(ms[2])

	regPath := regexp.MustCompile(`^("[^"]+"|[^"]+)$`)
	if !regPath.MatchString(path) {
		log.Fatalf("bad path format: %s", path)
	}

	path = strings.Trim(path, `"`)

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
	kvMap := parseKV(aliasLst)
	return kvMap
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

func parsePreq() {

}

type fieldInfo struct {
	Name       string
	Type       string
	Tag        string
	Alias      string
	IsExported bool
}

func extractStructFields(pkgPath, typeName string) ([]fieldInfo, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkgPath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse package: %w", err)
	}

	var fields []fieldInfo
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}
				if genDecl.Tok != token.TYPE {
					continue
				}
				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					if typeSpec.Name.Name != typeName {
						continue
					}
					structType, ok := typeSpec.Type.(*ast.StructType)
					if !ok {
						continue
					}
					for _, field := range structType.Fields.List {
						var alias string
						var rawTag string
						if field.Tag != nil {
							rawTag = field.Tag.Value
							alias = parseFieldAlias(rawTag)
						}
						for _, name := range field.Names {
							fields = append(fields, fieldInfo{
								Name:       name.Name,
								Type:       exprToString(field.Type),
								Tag:        rawTag,
								Alias:      alias,
								IsExported: name.IsExported(),
							})
						}
						// Handle anonymous fields (embedded structs)
						if len(field.Names) == 0 {
							fields = append(fields, fieldInfo{
								Name:       exprToString(field.Type),
								Type:       exprToString(field.Type),
								Tag:        rawTag,
								Alias:      alias,
								IsExported: true,
							})
						}
					}
				}
			}
		}
	}
	return fields, nil
}

func getPkgDir(importPath string) (string, error) {
	pkg, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		return "", err
	}
	return pkg.Dir, nil
}

func parseFieldAlias(tag string) string {
	t := reflect.StructTag(strings.Trim(tag, "`"))
	aliasReg := regexp.MustCompile(`alias=(\w+)`)
	ms := aliasReg.FindStringSubmatch(t.Get("shoot"))
	if len(ms) == 0 {
		return ""
	}
	return ms[1]
}
