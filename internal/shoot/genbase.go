package shoot

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
)

type GenBase struct {
	commonFlags *CommonFlags
	pkg         *Package
	worker      GenWorker
}

func (g *GenBase) SetWorker(worker GenWorker) {
	g.worker = worker
}

func (g *GenBase) CommonFlags() *CommonFlags {
	return g.commonFlags
}

func (g *GenBase) Package() *Package {
	return g.pkg
}

func (g *GenBase) ParseCommonFlags(sub *flag.FlagSet) {
	typeNames := sub.String("type", "", "comma-separated list of type names")
	fileName := sub.String("file", "", "the targe go file to generate, typical value: $GOFILE")
	separate := sub.Bool("separate", false, "each type has its own go file")
	s := sub.Bool("s", false, "each type has its own go file (alias for separate)")
	verbose := sub.Bool("verbose", false, "verbose output")
	v := sub.Bool("v", false, "verbose output (alias for verbose)")
	raw := sub.Bool("raw", false, "raw source")
	r := sub.Bool("r", false, "raw source (alias for raw)")

	sub.Parse(flag.Args()[1:]) //e.g. enum -bit type=YourType ./testdata

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

	g.commonFlags = &CommonFlags{
		TypeNames: typNames,
		FileName:  *fileName,
		Separate:  *s || *separate || *fileName == "",
		Dir:       dir,
		Verbose:   *v || *verbose,
		Raw:       *r || *raw,
	}
}

func (g *GenBase) fileName(name string, pkgScope bool) string {
	if pkgScope {
		return fmt.Sprintf("%s%s.%s.go", Cmd, g.worker.SubCmd(), name)
	}
	gofile := g.commonFlags.FileName
	if name == "" {
		return fmt.Sprintf("%s.%s%s.go", gofile, Cmd, g.worker.SubCmd())
	}
	if !ast.IsExported(name) {
		name = "_" + name
	}
	return fmt.Sprintf("%s.%s.go", gofile, strings.ToLower(name))
}

// parsePackage analyzes the single package constructed from the patterns and tags.
func (g *GenBase) parsePackage(patterns []string) {
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
func (g *GenBase) addPackage(pkg *packages.Package) {
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

func (g *GenBase) Generate() map[string][]byte {
	pat := g.commonFlags.Dir
	if g.commonFlags.FileName != "" {
		pat = filepath.Join(pat, g.commonFlags.FileName)
	}

	g.parsePackage([]string{pat})

	typeNames := g.commonFlags.TypeNames
	if len(typeNames) > 0 {
		for _, typName := range typeNames {
			gofile := getGoFile(g.pkg.pkg, typName)
			if g.commonFlags.FileName == "" {
				g.commonFlags.FileName = gofile
			} else if g.commonFlags.FileName != gofile {
				log.Fatalf("types are not in the same file")
			}
		}
	} else {
		typeNames = g.worker.TypeNames()
	}

	srcMap := make(map[string][]byte)
	var srcList [][]byte
	for _, typName := range typeNames {
		if ok := g.worker.Do(typName); !ok {
			continue
		}
		src := g.generate(typName)
		if len(src) == 0 {
			continue
		}
		if g.commonFlags.Separate {
			srcMap[g.fileName(typName, false)] = src
		} else {
			srcList = append(srcList, src)
		}
	}
	if len(srcList) > 0 {
		src, err := MergeSources(srcList...)
		if err != nil {
			log.Fatalf("merge sources error: %s", err)
		}
		if len(src) > 0 {
			srcMap[g.fileName("", false)] = src
		}
	}

	return srcMap
}

func (g *GenBase) generate(typeName string) []byte {
	data := g.worker.Data()
	cflags := g.CommonFlags()
	pkg := g.Package()

	var buff bytes.Buffer
	tmpl, err := template.New(g.worker.SubCmd()).Funcs(data.Transfers()).Parse(g.worker.TmplTxt())
	if err != nil {
		log.Fatalf("parsing template: %s", err)
	}

	data.SetCmd(strings.Join(append([]string{Cmd}, flag.Args()...), " "))
	data.SetTypeName(typeName)
	data.SetPackageName(pkg.Name())

	if cflags.Verbose {
		log.Printf("[debug]:\n%+v", data)
	}

	err = tmpl.Execute(&buff, data)
	if err != nil {
		log.Fatalf("executing template: %s", err)
	}
	src := buff.Bytes()
	if cflags.Verbose {
		log.Printf("[debug]:\n%s", string(src))
	}

	if cflags.Raw {
		return src //typically used for debugging
	}

	src, err = FormatSrc(src)
	if err != nil {
		log.Fatalf("format source: %s", err)
	}
	return src
}
