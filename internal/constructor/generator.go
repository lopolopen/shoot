package constructor

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/lopolopen/shoot/shoot"
	"golang.org/x/tools/go/packages"
)

const SubCmd = "new"

//go:embed constructor.tmpl
var tmplTxt string

//go:embed option.tmpl
var tmplTxtOpt string

// Generator holds the state of the analysis.
type Generator struct {
	flags *Flags
	pkg   *Package
	data  *Data
}

func New() *Generator {
	return &Generator{
		data: &Data{},
	}
}

func (g *Generator) ParseFlags() {
	sub := flag.NewFlagSet(SubCmd, flag.ExitOnError)
	typeNames := sub.String("type", "", "comma-separated list of type names")
	fileName := sub.String("file", "", "the targe go file to generate, typical value: $GOFILE")
	getset := sub.Bool("getset", false, "generate Get/Set method for the type")
	json := sub.Bool("json", false, "generate MarshalJSON/UnmarshalJSON method for the type")
	option := sub.Bool("option", false, "generate functional option pattern constructor")
	opt := sub.Bool("opt", false, "generate functional option pattern constructor (alias for -option)")
	separate := sub.Bool("separate", false, "each type has its own go file")
	s := sub.Bool("s", false, "each type has its own go file (alias for -separate)")
	verbose := sub.Bool("verbose", false, "verbose output")
	v := sub.Bool("v", false, "verbose output (alias for -separate)")

	sub.Parse((flag.Args()[1:]))

	if *typeNames == "" && *fileName == "" {
		sub.Usage()
		os.Exit(2)
	}

	var typNames []string
	if *typeNames != "" {
		typNames = strings.Split(*typeNames, ",")
	}
	if *fileName != "" && !strings.HasSuffix(*fileName, ".go") {
		log.Fatal("file must be a go file")
	}

	g.flags = &Flags{
		typeNames: typNames,
		fileName:  *fileName,
		getset:    *getset,
		json:      *json,
		opt:       *opt || *option,
		separate:  *s || *separate || *fileName == "",
		verbose:   *v || *verbose,
	}
}

func (g *Generator) Generate() map[string][]byte {
	pat := g.flags.fileName
	if pat == "" {
		pat = "."
	}

	g.parsePackage([]string{pat})

	g.data.BaseData = shoot.BaseData{
		Cmd:         strings.Join(append([]string{shoot.Cmd}, flag.Args()...), " "),
		PackageName: g.pkg.name,
	}

	if len(g.flags.typeNames) == 0 {
		g.parseTypeNames()
	} else {
		for _, typName := range g.flags.typeNames {
			gofile := getGoFile(g.pkg.pkg, typName)
			if g.flags.fileName == "" {
				g.flags.fileName = gofile
			} else if g.flags.fileName != gofile {
				log.Fatalf("types are not in the same file")
			}
		}
	}

	srcMap := make(map[string][]byte)
	if g.flags.separate {
		//each type has its own separate file
		for _, typName := range g.flags.typeNames {
			srcMap[g.fileName(typName, false)] = g.generate(typName)
		}
	} else {
		//types in one file
		var srcList [][]byte
		for _, typName := range g.flags.typeNames {
			srcList = append(srcList, g.generate(typName))
		}
		src, err := mergeGoSources(srcList...)
		if err != nil {
			log.Fatalf("merge sources error: %s", err)
		}
		srcMap[g.fileName("", false)] = src
	}
	if g.data.Option {
		srcMap[g.fileName("opt", true)] = g.generateOpt()
	}
	return srcMap
}

// parsePackage analyzes the single package constructed from the patterns and tags.
// parsePackage exits if there is an error.
func (g *Generator) parsePackage(patterns []string) {
	cfg := &packages.Config{
		Mode:  packages.LoadSyntax,
		Tests: false,
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found", len(pkgs))
	}

	g.addPackage(pkgs[0])
}

func (g *Generator) generate(typeName string) []byte {
	g.data.PreRegister()
	g.makeNew(typeName)
	g.makeOpt(typeName)
	g.makeGetSet(typeName)
	g.makeJson(typeName)

	var buff bytes.Buffer
	tmpl, err := template.New(SubCmd).Funcs(g.data.Transfers()).Parse(tmplTxt)
	if err != nil {
		log.Fatalf("parsing template: %s", err)
	}
	g.data.TypeName = typeName
	err = tmpl.Execute(&buff, g.data)
	if err != nil {
		log.Fatalf("executing template: %s", err)
	}

	src := buff.Bytes()
	if g.flags.verbose {
		log.Printf("[debug]:\n%s", string(src))
	}
	src, err = format.Source(src)
	if err != nil {
		log.Fatalf("format source: %s", err)
	}
	return src
}

func (g *Generator) fileName(name string, pkgScope bool) string {
	if pkgScope {
		return fmt.Sprintf("%s%s_%s.go", shoot.Cmd, SubCmd, name)
	}
	gofile := g.flags.fileName
	if name == "" {
		return fmt.Sprintf("%s_%s%s.go", gofile, shoot.Cmd, SubCmd)
	}
	if !ast.IsExported(name) {
		name = name + "_"
	}
	return fmt.Sprintf("%s_%s.go", gofile, strings.ToLower(name))
}

// addPackage adds a type checked Package and its syntax files to the generator.
func (g *Generator) addPackage(pkg *packages.Package) {
	g.pkg = &Package{
		pkg:   pkg,
		name:  pkg.Name,
		defs:  pkg.TypesInfo.Defs,
		files: make([]*File, len(pkg.Syntax)),
	}

	for i, file := range pkg.Syntax {
		g.pkg.files[i] = &File{
			file: file,
			pkg:  g.pkg,
		}
	}
}

func getGoFile(pkg *packages.Package, typeName string) string {
	for _, obj := range pkg.TypesInfo.Defs {
		if obj != nil && obj.Name() == typeName {
			pos := pkg.Fset.Position(obj.Pos())
			return filepath.Base(pos.Filename)
		}
	}
	return ""
}

func (g *Generator) parseTypeNames() {
	var typeNames []string
	for _, f := range g.pkg.files {
		ast.Inspect(f.file, func(n ast.Node) bool {
			ts, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}
			typeNames = append(typeNames, ts.Name.Name)
			return false
		})
	}
	g.flags.typeNames = typeNames
}

// MergeGoSources 接收任意数量的 Go 源码字符串，合并为一个格式化后的源文件
func mergeGoSources(sources ...[]byte) ([]byte, error) {
	if len(sources) == 0 {
		return nil, fmt.Errorf("没有传入任何源码")
	}

	fset := token.NewFileSet()
	var files []*ast.File

	// 解析所有源码
	for i, src := range sources {
		file, err := parser.ParseFile(fset, fmt.Sprintf("file%d.go", i), src, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("解析第 %d 段源码失败: %w", i+1, err)
		}
		files = append(files, file)
	}

	// 用第一段源码的包名作为统一包名
	pkgName := files[0].Name.Name

	// 准备合并后的文件
	merged := &ast.File{
		Name:  ast.NewIdent(pkgName),
		Decls: []ast.Decl{},
	}

	// import 去重
	importMap := map[string]bool{}
	for _, file := range files {
		file.Name.Name = pkgName // 强制统一包名
		for _, imp := range file.Imports {
			if !importMap[imp.Path.Value] {
				importMap[imp.Path.Value] = true
				merged.Decls = append(merged.Decls, &ast.GenDecl{
					Tok:   token.IMPORT,
					Specs: []ast.Spec{imp},
				})
			}
		}
	}

	// 合并除 import 外的其他声明
	for _, file := range files {
		for _, decl := range file.Decls {
			if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok == token.IMPORT {
				continue // 跳过 import，已处理
			}
			merged.Decls = append(merged.Decls, decl)
		}
	}

	// 格式化输出
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, merged); err != nil {
		return nil, fmt.Errorf("格式化失败: %w", err)
	}

	return buf.Bytes(), nil
}
