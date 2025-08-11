package constructor

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"go/format"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/lopolopen/shoot/shoot"
	"github.com/lopolopen/shoot/shoot/transfer"
)

const SubCmd = "new"

//go:embed constructor.tmpl
var tmplTxt string

type Data struct {
	shoot.Meta
	NewList        []string
	AssignmentList []string
	GetterList     []string
	SetterList     []string
}

type Flags struct {
	TypeNames []string
}

func parse(sub *flag.FlagSet) Flags {
	typeNames := sub.String("type", "", "comma-separated list of type names")
	sub.Parse(flag.Args()[1:])
	return Flags{
		TypeNames: strings.Split(*typeNames, ","),
	}
}

func usage(sub *flag.FlagSet) {
	//TODO:
	sub.Usage = func() {
		log.Printf("Usage: %s %s [options]", shoot.Cmd, SubCmd)
		log.Println()
		sub.PrintDefaults()
	}
}

func Gen(sub *flag.FlagSet) error {
	usage(sub)
	flags := parse(sub)

	var g shoot.Generator
	dir := "."

	tmpl, err := template.New(SubCmd).Funcs(template.FuncMap{
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
	}).Parse(tmplTxt)
	if err != nil {
		log.Fatalf("parsing template: %s", err)
	}

	g.ParsePackage([]string{dir}, []string{})

	data := &Data{
		Meta: shoot.Meta{
			Cmd:         strings.Join(append([]string{shoot.Cmd}, flag.Args()...), " "),
			PackageName: g.Pkg().Name(),
			TypeName:    flags.TypeNames[0],
		},
		NewList: []string{
			"name string",
			"age int",
		},
		AssignmentList: []string{
			"name",
			"age",
		},
		GetterList: []string{
			"name",
			"age",
		},
		SetterList: []string{
			"age",
		},
	}

	// log.Printf("%+v", data)

	var buff bytes.Buffer
	err = tmpl.Execute(&buff, data)
	if err != nil {
		log.Fatalf("executing template: %s", err)
	}

	src, err := format.Source(buff.Bytes())
	if err != nil {
		log.Fatalf("format source: %s", err)
	}

	// src := buff.Bytes()

	name := strings.ToLower(fmt.Sprintf("%s_%s_%s.go", data.TypeName, shoot.Cmd, SubCmd))
	output := filepath.Join(dir, name)

	outFile, err := os.Create(output)
	if err != nil {
		log.Fatalf("creating output file: %s", err)
	}
	defer outFile.Close()
	_, err = outFile.Write(src)
	if err != nil {
		log.Fatal()
	}

	return nil
}
