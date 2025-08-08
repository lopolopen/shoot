package enumer

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

const SubCmd = "enum"

type Data struct {
	shoot.Meta
	Bitwise bool
}

//go:embed enumer.tmpl
var tmplTxt string

func Gen() error {
	var g shoot.Generator
	dir := "."

	sub := flag.NewFlagSet(SubCmd, flag.ExitOnError)
	typeNames := sub.String("type", "", "comma-separated list of type names")
	bitwise := sub.Bool("bit", false, "generate bitwise enumerations")
	sub.Parse(flag.Args()[1:])
	typs := strings.Split(*typeNames, ",")
	log.Println("Type names:", typs)

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
			PackageName: g.Pkg().Name,
			TypeName:    typs[0],
		},
		Bitwise: *bitwise,
	}

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

	name := strings.ToLower(fmt.Sprintf("%s_%s_%s.go", typs[0], shoot.Cmd, SubCmd))
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
