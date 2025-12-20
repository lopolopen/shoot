package mapper

import (
	"sort"
	"strings"
)

const dot = "."

func (g *Generator) makePtrPath() {
	g.data.SrcPtrTypeMap = make(map[string]string)
	// g.data.DestPtrTypeMap = make(map[string]string)

	var ptrPathList []string
	for _, f := range g.exportedFields {
		if !strings.Contains(f.path, dot) {
			continue
		}
		for p, t := range g.ptrTypeMap {
			if _, ok := g.data.SrcPtrTypeMap[p]; ok {
				continue
			}

			if strings.HasPrefix(f.path, p) {
				ptrPathList = append(ptrPathList, p)
				g.data.SrcPtrTypeMap[p] = t
			}
		}
	}
	sort.Strings(ptrPathList)
	g.data.SrcPtrPathList = ptrPathList

}
