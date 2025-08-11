package enumer

import (
	_ "embed"
	"flag"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/lopolopen/shoot/shoot"
	"github.com/lopolopen/shoot/shoot/transfer"
)

const SubCmd = "enum"

//go:embed enumer.tmpl
var tmplTxt string

type Data struct {
	shoot.Meta
	String        bool
	NoChangeGuard string
	NameIndexDecl string
	Bitwise       bool
	ValueNames    []string
}

type Flags struct {
	TypeNames []string
	String    bool
	Bitwise   bool
}

func parse(sub *flag.FlagSet) Flags {
	str := sub.Bool("str", false, "generate String() method for the type (alias for -string)")
	string := sub.Bool("string", false, "generate String() method for the type")
	bit := sub.Bool("bit", false, "generate bitwise enumerations (alias for -bitwise)")
	bitwise := sub.Bool("bitwise", false, "generate bitwise enumerations")
	typeNames := sub.String("type", "", "comma-separated list of type names")

	sub.Parse(flag.Args()[1:])
	return Flags{
		TypeNames: strings.Split(*typeNames, ","),
		String:    *string || *str,
		Bitwise:   *bitwise || *bit,
	}
}

func Gen(sub *flag.FlagSet) error {
	flags := parse(sub)
	var g shoot.Generator
	dir := "."

	g.ParsePackage([]string{dir}, []string{})

	// index, name := createIndexAndNameDecl(g.Pkg().Files()[0].Values(), flags.TypeNames[0], "")

	data := &Data{
		Meta: shoot.Meta{
			Cmd:         strings.Join(append([]string{shoot.Cmd}, flag.Args()...), " "),
			PackageName: g.Pkg().Name(),
			TypeName:    flags.TypeNames[0],
		},
		String:        flags.String,
		NoChangeGuard: "",
		// NameIndexDecl: fmt.Sprintf("const %s\n%s", name, index),
		Bitwise: flags.Bitwise,
	}

	fileName := strings.ToLower(fmt.Sprintf("%s_%s_%s.go", data.TypeName, shoot.Cmd, SubCmd))
	output := filepath.Join(dir, fileName)

	funcs := template.FuncMap{
		"camel":  transfer.ToCamelCase,
		"pascal": transfer.ToPascalCase,
		"firstL": transfer.FirtLower,
		"typeof": func(name string) string {
			m := map[string]string{
				"name": "string",
				"age":  "int",
			}
			return m[name]
		},
	}
	shoot.Render(dir, output, SubCmd, tmplTxt, funcs, data)

	// outFile, err := os.Create(output)
	// if err != nil {
	// 	log.Fatalf("creating output file: %s", err)
	// }
	// defer outFile.Close()
	// _, err = outFile.Write(src)
	// if err != nil {
	// 	log.Fatal()
	// }

	return nil
}
