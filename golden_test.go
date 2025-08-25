package main

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/lopolopen/shoot/internal/constructor"
	"github.com/sebdah/goldie/v2"
)

type Golden struct {
	cmd   string
	names []string
}

var goldens = []Golden{
	{
		cmd: "shoot new -getset -type=User ./testdata",
		names: []string{
			"new_getset.go_user.go",
		},
	},
	{
		cmd: "shoot new -opt -type=Conf ./testdata",
		names: []string{
			"new_opt.go_conf.go",
			"shootnew_opt.go",
		},
	},
	{
		cmd: "shoot new -file=new_file.go ./testdata",
		names: []string{
			"new_file.go_shootnew.go",
		},
	},
}

func TestGenerate_Golden(t *testing.T) {
	g := goldie.New(t,
		goldie.WithFixtureDir("testdata"),
	)

	for _, test := range goldens {
		srcMap := generate(t, test)

		for _, name := range test.names {
			got, ok := srcMap[name]
			if !ok {
				var keys []string
				for key := range srcMap {
					keys = append(keys, key)
				}
				t.Errorf("expected file: %s, got: %v", name, keys)
			}
			g.Assert(t, name, got)
		}
	}
}

func generate(t *testing.T, test Golden) map[string][]byte {
	os.Args = strings.Split(test.cmd, " ")
	flag.Parse()

	gen := constructor.New()
	gen.ParseFlags()

	srcMap := gen.Generate()
	if len(srcMap) != len(test.names) {
		t.Errorf("expected count: %d, got: %d", len(test.names), len(srcMap))
	}
	return srcMap
}
