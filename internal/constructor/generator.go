package constructor

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
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
	getset := sub.Bool("getset", false, "generate Get/Set method for the type")
	json := sub.Bool("json", false, "generate MarshalJSON/UnmarshalJSON method for the type")
	opt := sub.Bool("opt", false, "generate functional option pattern constructor (alias for -option)")
	verbose := sub.Bool("v", false, "verbose outpot for debug")

	var typNames []string
	sub.Parse((flag.Args()[1:]))
	if (*typeNames) == "" {
		gofile := sub.Arg(0)
		if !strings.HasSuffix(gofile, ".go") {
			sub.Usage()
			os.Exit(2)
		}
		g.data.GoFile = gofile
	} else {
		typNames = strings.Split(*typeNames, ",")
	}

	g.flags = &Flags{
		typeNames: typNames,
		getset:    *getset,
		json:      *json,
		opt:       *opt,
		verbose:   *verbose,
	}
}

func (g *Generator) Generate() map[string][]byte {
	pat := g.data.GoFile
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
	}

	srcMap := make(map[string][]byte)
	for _, typName := range g.flags.typeNames {
		srcMap[typName] = g.generate(typName)
	}
	if g.data.Option {
		srcMap["@new_opt"] = g.generateOpt()
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

func (g *Generator) FileName(typeName string) string {
	prefix := shoot.Cmd
	postfix := ""
	if strings.HasPrefix(typeName, "@") {
		typeName = typeName[1:]
	} else {
		if !strings.HasPrefix(typeName, "_") {
			prefix = getGoFile(g.pkg.pkg, typeName)
		}
		if !ast.IsExported(typeName) {
			postfix = "_x"
		}
	}
	fileName := strings.ToLower(fmt.Sprintf("%s_%s%s", prefix, typeName, postfix))
	return fmt.Sprintf("%s.go", fileName)
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
