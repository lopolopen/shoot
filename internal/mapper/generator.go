package mapper

import (
	_ "embed"
	"flag"
	"go/ast"
	"log"
	"path/filepath"
	"strings"

	"github.com/lopolopen/shoot/internal/shoot"
	"golang.org/x/tools/go/packages"
)

const SubCmd = "map"

//go:embed mapper.tmpl
var tmplTxt string

// Generator holds the state of the analysis.
type Generator struct {
	*shoot.GeneratorBase
	flags           *Flags
	data            *TmplData
	destpkg         *packages.Package
	mapperpkg       *packages.Package
	srcExpList      []Field
	destExpList     []Field
	mappingFuncList []Func
	assignedDestSet map[string]bool
	assignedSrcSet  map[string]bool
	tagMap          map[string]string
}

func New() *Generator {
	g := &Generator{
		GeneratorBase: shoot.NewGeneratorBase(SubCmd, tmplTxt),
	}
	return g
}

func (g *Generator) ParseFlags() {
	sub := flag.NewFlagSet(SubCmd, flag.ExitOnError)
	path := sub.String("path", "", "destination package path to map to")
	destNames := sub.String("dest", "", "destination type names to map to (must align to -type)")
	g.ParseCommonFlags(sub)

	typNames := g.CommonFlags().TypeNames
	typMap := make(map[string]string)
	if len(typNames) == 0 {
		if *destNames != "" {
			log.Fatal("❌ -dest only works when -type used")
		}
	} else {
		destTypNames := strings.Split(*destNames, ",")
		if len(destTypNames) != len(typNames) {
			log.Fatal("❌ -dest must align to -type")
		}
		for i, n := range typNames {
			typMap[n] = destTypNames[i]
		}
	}

	g.flags = &Flags{
		destDir:   *path,
		destTypes: typMap,
	}
}

func (g *Generator) LoadPackage() {
	cwd := "."
	destDir := g.flags.destDir
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedSyntax |
			packages.NeedTypes |
			packages.NeedTypesInfo,
	}
	pkgs, err := loadPkgs(cfg, cwd, destDir)
	if err != nil {
		log.Fatalf("❌ %s", err)
	}
	if len(pkgs) != 2 {
		log.Fatalf("❌ error: %d packages found", len(pkgs))
	}

	g.SetPkg(pkgs[cwd])
	g.destpkg = pkgs[destDir]
	g.GeneratorBase.LoadPackage()
}

func loadPkgs(cfg *packages.Config, patterns ...string) (map[string]*packages.Package, error) {
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*packages.Package, len(patterns))

	for _, pat := range patterns {
		isPath := strings.HasPrefix(pat, ".") || strings.HasPrefix(pat, "/")

		for _, pkg := range pkgs {
			if isPath {
				absPat, _ := filepath.Abs(pat)
				for _, f := range pkg.GoFiles {
					absFile, _ := filepath.Abs(f)
					rel, err := filepath.Rel(absPat, absFile)
					if err == nil && !strings.HasPrefix(rel, "..") {
						result[pat] = pkg
						break
					}
				}
			} else {
				if pkg.PkgPath == pat {
					result[pat] = pkg
				}
			}
			if result[pat] != nil {
				break
			}
		}
	}

	return result, nil
}

func (g *Generator) MakeData(srcTypeName string) any {
	g.data = NewTmplData()

	var destTypeName string
	if g.flags.destTypes != nil {
		if name, ok := g.flags.destTypes[srcTypeName]; ok {
			destTypeName = name
		}
	}
	if destTypeName == "" {
		destTypeName = srcTypeName
	}
	g.data.DestTypeName = destTypeName

	g.parseSrcFields(srcTypeName)
	g.parseDestFields(destTypeName)
	g.parseCustomMapMethods(srcTypeName, destTypeName)
	g.makeMatch()
	mapperTypeName := g.loadTypeMapperPkg(srcTypeName)
	if mapperTypeName != "" {
		g.parseMapper(mapperTypeName)
		g.makeTypeMismatch()
	}

	if g.destpkg.PkgPath == g.Pkg().PkgPath { //same package
		g.data.QualifiedDestTypeName = g.data.DestTypeName
	} else {
		g.data.QualifiedDestTypeName = g.destpkg.Name + "." + g.data.DestTypeName
	}
	g.data.DestPkgName = g.destpkg.Name

	g.data.SetTypeName(srcTypeName)
	g.data.SetPackageName(g.Pkg().Name)
	g.data.SetCmd(strings.Join(append([]string{shoot.Cmd}, flag.Args()...), " "))

	g.checkUnassigned()
	return g.data
}

func (g *Generator) ListTypes() []string {
	var typeNames []string
	for _, f := range g.Pkg().Syntax {
		ast.Inspect(f, func(n ast.Node) bool {
			if !g.testNode("", n) {
				return true
			}

			ts, _ := n.(*ast.TypeSpec)
			typeNames = append(typeNames, ts.Name.Name)
			return false
		})
	}
	return typeNames
}

func (g *Generator) testNode(typeName string, node ast.Node) bool {
	ts, ok := node.(*ast.TypeSpec)
	if !ok {
		return false
	}

	if typeName != "" && ts.Name.Name != typeName {
		return false
	}

	_, ok = ts.Type.(*ast.StructType)
	if !ok {
		return false
	}
	return true
}
