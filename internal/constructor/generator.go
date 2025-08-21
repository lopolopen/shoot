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

	sub.Parse((flag.Args()[1:]))
	if (*typeNames) == "" {
		sub.Usage()
		os.Exit(2)
	}

	g.flags = &Flags{
		typeNames: strings.Split(*typeNames, ","),
		getset:    *getset,
		json:      *json,
		opt:       *opt,
		verbose:   *verbose,
	}
}

func (g *Generator) Generate() map[string][]byte {
	g.parsePackage([]string{"."}, []string{"."})

	g.data.BaseData = shoot.BaseData{
		Cmd:         strings.Join(append([]string{shoot.Cmd}, flag.Args()...), " "),
		PackageName: g.pkg.name,
	}

	srcMap := make(map[string][]byte)
	for _, typName := range g.flags.typeNames {
		srcMap[typName] = g.generate(typName)
	}
	if g.data.Option {
		srcMap["_opt_"] = g.generateOpt()
	}
	return srcMap
}

// parsePackage analyzes the single package constructed from the patterns and tags.
// parsePackage exits if there is an error.
func (g *Generator) parsePackage(patterns []string, tags []string) {
	cfg := &packages.Config{
		Mode:  packages.LoadSyntax,
		Tests: false,
		Fset:  token.NewFileSet(),
		ParseFile: func(fset *token.FileSet, filename string, src []byte) (*ast.File, error) {
			return parser.ParseFile(fset, filename, src, parser.ParseComments)
		},
		// 关键：禁止用 export data，强制读源码
		Overlay: nil, // 确保不被替换
		Env: append(os.Environ(),
			"GOPACKAGESDRIVER=off", // 禁用外部 driver
		),
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
	fileName := strings.ToLower(fmt.Sprintf("%s_%s_%s", shoot.FilePrefix, SubCmd, typeName))
	if ast.IsExported(typeName) || strings.HasPrefix(typeName, "_") {
		return fileName + ".go"
	}
	return fileName + "_.go"
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
