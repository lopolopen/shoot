package mapper

import (
	_ "embed"
	"flag"
	"fmt"
	"go/ast"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lopolopen/shoot/internal/shoot"
	"github.com/lopolopen/shoot/internal/tools/logx"
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
			logx.Fatal("-dest only works when -type used")
		}
	} else {
		destTypNames := strings.Split(*destNames, ",")
		if len(destTypNames) != len(typNames) {
			logx.Fatal("-dest must align to -type")
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
		logx.Fatalf("%s", err)
	}

	//todo: each path has only one pkg

	pkg := pkgs[cwd]
	if g.CommonFlags().FileName != "" {
		var fs []*ast.File
		for _, f := range pkg.Syntax {
			filename := pkg.Fset.File(f.Pos()).Name()
			if filepath.Base(filename) == g.CommonFlags().FileName {
				fs = append(fs, f)
				break
			}
		}
		pkg.Syntax = fs
	} else if shoot.Contains(g.CommonFlags().TypeNames, "*") {
	end:
		for i, f := range pkg.Syntax {
			for _, cg := range f.Comments {
				for _, c := range cg.List {
					if !findCmdLine(c.Text, g.CommonFlags().CmdLine) {
						continue
					}
					filename := filepath.Base(pkg.GoFiles[i])
					g.AllInOne = filename
					break end
				}
			}
		}
	}

	g.SetPkg(pkg)
	g.destpkg = pkgs[destDir]
	g.GeneratorBase.LoadPackage()
}

func findCmdLine(doc string, cmdline string) bool {
	pat := fmt.Sprintf("(?im)^//go:generate.*%s$", regexp.QuoteMeta(cmdline))
	regAll := regexp.MustCompile(pat)
	new := regAll.Match([]byte(doc))
	return new
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
	g.data = NewTmplData(g.CommonFlags().CmdLine)

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
	g.parseManual(srcTypeName, destTypeName)
	mapperTypeName := g.loadTypeMapperPkg(srcTypeName)
	if mapperTypeName != "" {
		g.parseMapper(mapperTypeName)
		g.makeMismatch() //priority: > makeMatch
	}
	g.makeMatch()

	if g.destpkg.PkgPath == g.Pkg().PkgPath { //same package
		g.data.QualifiedDestTypeName = g.data.DestTypeName
	} else {
		g.data.QualifiedDestTypeName = g.destpkg.Name + "." + g.data.DestTypeName
	}
	g.data.DestPkgName = g.destpkg.Name

	g.data.SetTypeName(srcTypeName)
	g.data.SetPackageName(g.Pkg().Name)

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
