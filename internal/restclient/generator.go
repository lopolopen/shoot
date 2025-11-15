package restclient

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"go/ast"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/lopolopen/shoot/internal/shoot"
	"golang.org/x/tools/go/packages"
)

const SubCmd = "rest"

//go:embed restclient.tmpl
var tmplTxt string

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
	verbose := sub.Bool("verbose", false, "verbose output")
	v := sub.Bool("v", false, "verbose output (alias for verbose)")
	separate := sub.Bool("separate", false, "each type has its own go file")
	s := sub.Bool("s", false, "each type has its own go file (alias for separate)")
	raw := sub.Bool("raw", false, "raw source")
	r := sub.Bool("r", false, "raw source (alias for raw)")

	sub.Parse((flag.Args()[1:])) //e.g. rest -type=YourType ./testdata

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
		BaseFlags: shoot.BaseFlags{
			TypeNames: typNames,
			FileName:  *fileName,
			Dir:       dir,
			Verbose:   *v || *verbose,
			Separate:  *s || *separate || *fileName == "",
			Raw:       *r || *raw,
		},
	}
}

func (g *Generator) Generate() map[string][]byte {
	pat := g.flags.Dir
	if g.flags.FileName != "" {
		pat = filepath.Join(pat, g.flags.FileName)
	}

	g.parsePackage([]string{pat})

	g.data.BaseData = shoot.BaseData{
		Cmd:         strings.Join(append([]string{shoot.Cmd}, flag.Args()...), " "),
		PackageName: g.pkg.name,
	}

	if len(g.flags.TypeNames) == 0 {
		g.parseTypeNames()
	} else {
		for _, typName := range g.flags.TypeNames {
			gofile := getGoFile(g.pkg.pkg, typName)
			if g.flags.FileName == "" {
				g.flags.FileName = gofile
			} else if g.flags.FileName != gofile {
				log.Fatalf("types are not in the same file")
			}
		}
	}

	srcMap := make(map[string][]byte)
	if g.flags.Separate {
		//each type has its own separate file
		for _, typName := range g.flags.TypeNames {
			srcMap[g.fileName(typName, false)] = g.generate(typName)
		}
	} else {
		//types in one file
		var srcList [][]byte
		for _, typName := range g.flags.TypeNames {
			srcList = append(srcList, g.generate(typName))
		}
		src, err := shoot.MergeSources(srcList...)
		if err != nil {
			log.Fatalf("merge sources error: %s", err)
		}
		srcMap[g.fileName("", false)] = src
	}
	return srcMap
}

func (g *Generator) generate(typeName string) []byte {
	g.data.PreRegister()
	g.cookClient(typeName)

	var buff bytes.Buffer
	tmpl, err := template.New(SubCmd).Funcs(g.data.Transfers()).Parse(tmplTxt)
	if err != nil {
		log.Fatalf("parsing template: %s", err)
	}
	g.data.TypeName = typeName
	if g.flags.Verbose {
		log.Printf("[debug]:\n%+v", g.data)
	}
	err = tmpl.Execute(&buff, g.data)
	if err != nil {
		log.Fatalf("executing template: %s", err)
	}
	src := buff.Bytes()
	if g.flags.Verbose {
		log.Printf("[debug:]\n%s", string(src))
	}

	if g.flags.Raw {
		return src //typically used for debugging
	}

	src, err = shoot.FormatSrc(src)
	if err != nil {
		log.Fatalf("format source: %s", err)
	}
	return src
}

func (g *Generator) FileName(typeName string) string {
	return strings.ToLower(fmt.Sprintf("%s_%s_%s.go", typeName, shoot.Cmd, SubCmd))
}

// parsePackage analyzes the single package constructed from the patterns and tags.
func (g *Generator) parsePackage(patterns []string) {
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedSyntax |
			packages.NeedTypes |
			packages.NeedTypesInfo,
		Tests: false,
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		log.Fatal(err)
	}
	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found", len(pkgs))
		fmt.Println(pkgs)
	}
	g.addPackage(pkgs[0])
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

func (g *Generator) parseTypeNames() {
	var typeNames []string
	for _, f := range g.pkg.files {
		ast.Inspect(f.file, func(n ast.Node) bool {
			if !g.testNode("", n) {
				return true
			}

			ts, _ := n.(*ast.TypeSpec)
			typeNames = append(typeNames, ts.Name.Name)
			return false
		})
	}
	g.flags.TypeNames = typeNames
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

func (g *Generator) fileName(name string, pkgScope bool) string {
	if pkgScope {
		return fmt.Sprintf("%s%s_%s.go", shoot.Cmd, SubCmd, name)
	}
	gofile := g.flags.FileName
	if name == "" {
		return fmt.Sprintf("%s_%s%s.go", gofile, shoot.Cmd, SubCmd)
	}
	if !ast.IsExported(name) {
		name = name + "_"
	}
	return fmt.Sprintf("%s_%s.go", gofile, strings.ToLower(name))
}

func (g *Generator) testNode(typeName string, node ast.Node) bool {
	ts, ok := node.(*ast.TypeSpec)
	if !ok {
		return false
	}

	if typeName != "" && ts.Name.Name != typeName {
		return false
	}

	iface, ok := ts.Type.(*ast.InterfaceType)
	if !ok {
		return false
	}

	isRestClient := false
	for _, field := range iface.Methods.List {
		if len(field.Names) > 0 {
			continue
		}

		typ := g.pkg.pkg.TypesInfo.Types[field.Type].Type
		named, ok := typ.(*types.Named)
		if !ok {
			continue
		}
		obj := named.Obj()
		pkgPath := obj.Pkg().Path()
		if pkgPath == shoot.SelfPkgPath && obj.Name() == "RestClient" {
			isRestClient = true
			break
		}
	}
	return isRestClient
}
