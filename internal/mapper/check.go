package mapper

import (
	"fmt"
	"sort"
	"strings"

	"github.com/lopolopen/shoot/internal/tools/logx"
	"github.com/lopolopen/shoot/internal/transfer"
)

//tips: write src means read rest, and vice versa

func (g *Generator) makeReadWriteCheck() {
	g.makeReadCond()
	g.makeToSrcFromDest()
	g.makeFromSrcToDest()
	logx.DebugJSONln(g.data.ExactMatchMap)
}

func (g *Generator) makeToSrcFromDest() {
	g.data.DestNeedReadCheckMap = make(map[string]string)

	readDest := make(map[string]string)
	var xWriteSrc []string

	for _, f := range g.exportedFields {
		if g.assignedSrcSet[f.name] { //without writing, no reading
			continue
		}
		c := false
		for path := range g.assignedSrcSet { //Model covers Model.ID
			if f.CoveredBy(path) {
				c = true
				break
			}
		}
		if c {
			continue
		}

		src := f.name
		if d, ok := g.data.ExactMatchMap[src]; ok { //d, ok
			readDest[src] = d
			continue
		}
		if _, ok := g.data.DestToSrcTypeMap[src]; ok { //_, ok
			readDest[src] = g.data.ConvMatchMap[src]
			continue
		}
		if _, ok := g.data.DestToSrcFuncMap[src]; ok {
			readDest[src] = g.data.MismatchFuncMap[src]
			continue
		}
		if d, ok := g.data.MismatchSubMap[src]; ok {
			readDest[src] = d
			continue
		}
		if d, ok := g.data.MismatchSubListMap[src]; ok {
			readDest[src] = d
			continue
		}

		xWriteSrc = append(xWriteSrc, f.path)
	}

	for s, d := range readDest {
		if _, ok := g.destReadPathsMap[d]; ok {
			g.data.DestNeedReadCheckMap[s] = d
		}
	}

	if len(xWriteSrc) > 0 {
		names := strings.Join(xWriteSrc, ", ")
		reportWarn(g.data.PackageName+"."+g.data.TypeName, names)
	}
}

func (g *Generator) makeFromSrcToDest() {
	g.data.SrcNeedReadCheckMap = make(map[string]string)

	readSrc := make(map[string]string)
	var xWriteDest []string
outer:
	for _, f := range g.destExportedFields {
		if g.assignedDestSet[f.name] { //without writing, no reading
			continue
		}
		for path := range g.assignedDestSet { //Model covers Model.ID
			if f.CoveredBy(path) {
				continue outer
			}
		}

		dest := f.name
		for s, d := range g.data.ExactMatchMap { //for s, d
			if dest == d {
				readSrc[s] = dest
				continue outer
			}
		}
		for s := range g.data.SrcToDestTypeMap { //for s
			if dest == g.data.ConvMatchMap[s] {
				readSrc[s] = dest
				continue outer
			}
		}
		for s := range g.data.SrcToDestFuncMap {
			if dest == g.data.MismatchFuncMap[s] {
				readSrc[s] = dest
				continue outer
			}
		}
		for s, d := range g.data.DestMismatchSubMap {
			if dest == d {
				readSrc[s] = dest
				continue outer
			}
		}
		for s, d := range g.data.MismatchSubListMap {
			if dest == d {
				readSrc[s] = dest
				continue outer
			}
		}

		xWriteDest = append(xWriteDest, f.name)
	}

	for s, d := range readSrc {
		if _, ok := g.readPathsMap[s]; ok {
			g.data.SrcNeedReadCheckMap[s] = d
		}
	}

	if len(xWriteDest) > 0 {
		names := strings.Join(xWriteDest, ", ")
		reportWarn(g.data.QualifiedDestTypeName, names)
	}
}

func reportWarn(typename, fields string) {
	logx.Warnf("%s: these fields are never assigned:\n\t%s", typename, fields)
}

func (g *Generator) makeReadCond() {
	g.RegisterTransfer("condofread", transfer.ID)

	g.readPathsMap = make(map[string][]string)
	g.destReadPathsMap = make(map[string][]string)
	prepareReadPaths(g.exportedFields, g.ptrTypeMap, g.readPathsMap)
	prepareReadPaths(g.destExportedFields, g.destPtrTypeMap, g.destReadPathsMap)

	g.RegisterTransfer("condofread", func(v, name string, isSrc bool) string {
		var readPathsMap map[string][]string
		if isSrc {
			readPathsMap = g.readPathsMap
		} else {
			readPathsMap = g.destReadPathsMap
		}
		paths := readPathsMap[name]
		sort.Strings(paths)
		for i, p := range paths {
			paths[i] = fmt.Sprintf(" %s.%s != nil ", v, p)
		}
		return strings.Join(paths, "&&")
	})
}

func prepareReadPaths(fields []Field, ptrTypeMap map[string]string, readPathsMap map[string][]string) {
	for _, f := range fields {
		if !f.IsEmbeded() {
			continue
		}

		var readPaths []string
		ps := strings.Split(f.path, dot)
		for i := 1; i < len(ps); i++ {
			path := strings.Join(ps[:i], dot)
			_, ok := ptrTypeMap[path]
			if !ok {
				continue
			}
			readPaths = append(readPaths, path)
		}
		readPathsMap[f.name] = readPaths
	}
}
