package mapper

import (
	"strings"

	"github.com/lopolopen/shoot/internal/tools/logx"
)

func (g *Generator) checkUnassigned() {
	var src []string
	for _, f := range g.srcExpList {
		if g.assignedSrcSet[f.name] {
			continue
		}

		if _, ok := g.data.ExactMatchMap[f.name]; ok {
			continue
		}

		if _, ok := g.data.DestToSrcTypeMap[f.name]; ok {
			continue
		}

		if _, ok := g.data.DestToSrcFuncMap[f.name]; ok {
			continue
		}

		src = append(src, f.name)
	}
	if len(src) > 0 {
		names := strings.Join(src, ", ")
		logx.Warnf("%s: these fields are never assigned:\n\t%s", g.data.TypeName, names)
	}

	var dest []string
outer:
	for _, f := range g.destExpList {
		if g.assignedDestSet[f.name] {
			continue
		}

		for _, d := range g.data.ExactMatchMap {
			if f.name == d {
				continue outer
			}
		}
		for s := range g.data.SrcToDestTypeMap {
			if f.name == g.data.ConvMatchMap[s] {
				continue outer
			}
		}
		for s := range g.data.SrcToDestFuncMap {
			if f.name == g.data.MismatchMap[s] {
				continue outer
			}
		}

		dest = append(dest, f.name)
	}
	if len(dest) > 0 {
		names := strings.Join(dest, ", ")
		logx.Warnf("%s: these fields are never assigned:\n\t%s", g.data.QualifiedDestTypeName, names)
	}
}
