package constructor

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
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

	sub.Parse((flag.Args()[1:])) //e.g. new -getset -type=YourType ./testdata

	if *typeNames == "" && *fileName == "" {
		sub.Usage()
		os.Exit(2)
	}

	var typNames []string
	if *typeNames != "" {
		typNames = strings.Split(*typeNames, ",")
	}

	dir := sub.Arg(0) //e.g. ./testdata
	if dir == "" {
		dir = "."
	}

	if *fileName != "" {
		if !strings.HasSuffix(*fileName, ".go") {
			log.Fatal("file must be a go file")
		}
		fp := filepath.Join(dir, *fileName)
		_, err := os.Stat(fp)
		if !(err == nil || os.IsExist(err)) {
			log.Fatalf("file not exists: %s", fp)
		}
	}

	g.flags = &Flags{
		typeNames: typNames,
		fileName:  *fileName,
		getset:    *getset,
		json:      *json,
		opt:       *opt || *option,
		separate:  *s || *separate || *fileName == "",
		verbose:   *v || *verbose,
		dir:       dir,
	}
}

func (g *Generator) Generate() map[string][]byte {
	pat := g.flags.dir
	if g.flags.fileName != "" {
		pat = filepath.Join(pat, g.flags.fileName)
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

func mergeGoSources(sources ...[]byte) ([]byte, error) {
	if len(sources) == 0 {
		return nil, nil
	}

	fset := token.NewFileSet()
	var files []*ast.File

	for i, src := range sources {
		file, err := parser.ParseFile(fset, fmt.Sprintf("src%d.go", i), src, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("parse src%d.go: %w", i, err)
		}
		files = append(files, file)
	}

	pkgName := files[0].Name.Name

	importMap := map[string]bool{}
	var importDecls []ast.Decl
	for _, f := range files {
		for _, imp := range f.Imports {
			if !importMap[imp.Path.Value] {
				importMap[imp.Path.Value] = true
				importDecls = append(importDecls, &ast.GenDecl{
					Tok:   token.IMPORT,
					Specs: []ast.Spec{imp},
				})
			}
		}
	}

	var otherDecls []ast.Decl
	for _, f := range files {
		for _, decl := range f.Decls {
			if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok == token.IMPORT {
				continue
			}
			otherDecls = append(otherDecls, decl)
		}
	}

	var buf bytes.Buffer
	// header (first line comment)
	if len(files[0].Comments) > 0 {
		fmt.Fprint(&buf, "// ")
		fmt.Fprintln(&buf, files[0].Comments[0].Text())
		fmt.Fprintln(&buf)
	}
	// package
	fmt.Fprintf(&buf, "package %s\n\n", pkgName)
	// imports
	if len(importDecls) > 0 {
		for _, decl := range importDecls {
			printer.Fprint(&buf, fset, decl)
			fmt.Fprintln(&buf)
		}
		fmt.Fprintln(&buf)
	}
	// decls
	for _, decl := range otherDecls {
		printer.Fprint(&buf, fset, decl)
		fmt.Fprintln(&buf)
	}

	// return buf.Bytes(), nil

	out, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("format merged source: %w", err)
	}
	return out, nil
}
