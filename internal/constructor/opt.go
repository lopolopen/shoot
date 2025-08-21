package constructor

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"log"
	"regexp"
	"slices"
	"strings"
	"text/template"

	"github.com/lopolopen/shoot/internal"
	"github.com/lopolopen/shoot/shoot"
)

func (g *Generator) makeOpt(typeName string) {
	g.data.Register("defaultof", internal.ID)

	var defList []string
	defMap := make(map[string]string)

	for _, f := range g.pkg.files {
		ast.Inspect(f.file, func(n ast.Node) bool {
			ts, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}
			if ts.Name.Name != typeName {
				return false
			}
			st, ok := ts.Type.(*ast.StructType)
			if !ok {
				return true
			}

			for _, field := range st.Fields.List {
				for _, name := range field.Names {
					if name.Obj.Kind != ast.Var {
						continue
					}
					if ast.IsExported(name.Name) {
						continue
					}

					if field.Doc != nil {
						v, ok := parseDefault(field.Doc.Text())
						if ok {
							if !g.flags.opt {
								log.Fatalf("should not use instruction def(ault) when -opt(ion) disabled")
							}

							if len(g.data.NewList) != len(g.data.AllList)-len(g.data.ExportedList) {
								if slices.Contains(g.data.NewList, name.Name) {
									log.Fatalf("should not apply both instruction new and def(ault) to field %s of type %s", name.Name, typeName)
								}
							}

							defList = append(defList, name.Name)
							defMap[name.Name] = strings.TrimSpace(v)
						}
					}
				}
			}
			return false
		})
	}

	if !g.flags.opt {
		return
	}

	g.data.Register("defaultof", func(key string) string {
		return defMap[key]
	})
	g.data.DefaultList = defList
	g.data.Option = true
}

func (g *Generator) generateOpt() []byte {
	cmd := g.data.Cmd
	g.data.Cmd = fmt.Sprintf("%s %s -opt", shoot.Cmd, SubCmd)
	defer func() {
		g.data.Cmd = cmd
	}()
	var buff bytes.Buffer
	tmpl, err := template.New("opt").Funcs(g.data.Transfers()).Parse(tmplTxtOpt)
	if err != nil {
		log.Fatalf("parsing template: %s", err)
	}

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

func parseDefault(doc string) (string, bool) {
	regDef := regexp.MustCompile(`(?im)^shoot:.*?\Wdef(ault)?=([^;\n]+)(;.*|\s*)$`)
	ms := regDef.FindStringSubmatch(doc)
	for idx, m := range ms {
		if (m == "" || m == "ault") && idx+1 < len(ms) {
			return ms[idx+1], true
		}
	}
	return "", false
}
