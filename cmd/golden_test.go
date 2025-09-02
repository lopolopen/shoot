package main

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/lopolopen/shoot/internal/constructor"
	"github.com/lopolopen/shoot/internal/enumer"
	"github.com/lopolopen/shoot/internal/restclient"
	"github.com/lopolopen/shoot/internal/shoot"
	"github.com/sebdah/goldie/v2"
)

type Golden struct {
	cmd   string
	names []string
}

var goldens_new = []Golden{
	// {
	// 	cmd: "shoot new -getset -type=User ./notexists -v",
	// 	names: []string{
	// 		"new_getset.go_user.go",
	// 	},
	// },
	{
		cmd: "shoot new -getset -type=User ./testdata/constructor",
		names: []string{
			"new_getset.go_user.go",
		},
	},
	{
		cmd: "shoot new -opt -type=Conf ./testdata/constructor",
		names: []string{
			"new_opt.go_conf.go",
		},
	},
	{
		cmd: "shoot new -file=new_file.go ./testdata/constructor",
		names: []string{
			"new_file.go_shootnew.go",
		},
	},
}

var goldens_enum = []Golden{
	{
		cmd:   "shoot enum -bit -file=nothing.go ./testdata/enumer",
		names: []string{},
	},
	{
		cmd: "shoot enum -bit -type=FormatStyle ./testdata/enumer",
		names: []string{
			"enum_bit.go_formatstyle.go",
		},
	},
	{
		cmd: "shoot enum -json -type=Color ./testdata/enumer",
		names: []string{
			"enum_json.go_color.go",
		},
	},
}

var goldens_rest = []Golden{
	{
		cmd: "shoot rest -type=Client ./testdata/restclient",
		names: []string{
			"rest.go_client.go",
		},
	},
}

func TestShootNew_Golden(t *testing.T) {
	g := goldie.New(t,
		goldie.WithFixtureDir("testdata/constructor"),
	)

	for _, test := range goldens_new {
		srcMap := generate(test, constructor.New())

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

func TestShootEnum_Golden(t *testing.T) {
	g := goldie.New(t,
		goldie.WithFixtureDir("testdata/enumer"),
	)

	for _, test := range goldens_enum {
		srcMap := generate(test, enumer.New())

		if len(srcMap) != len(test.names) {
			t.Errorf("expected count: %d, got: %d", len(test.names), len(srcMap))
		}

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

func TestShootRest_Golden(t *testing.T) {
	g := goldie.New(t,
		goldie.WithFixtureDir("testdata/restclient"),
	)

	for _, test := range goldens_rest {
		srcMap := generate(test, restclient.New())

		if len(srcMap) != len(test.names) {
			t.Errorf("expected count: %d, got: %d", len(test.names), len(srcMap))
		}

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

func generate(test Golden, gen shoot.Generator) map[string][]byte {
	os.Args = strings.Split(test.cmd, " ")
	flag.Parse()

	gen.ParseFlags()

	srcMap := gen.Generate()
	return srcMap
}
