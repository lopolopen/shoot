package shoot

import (
	"flag"
	"log"
	"os"

	"github.com/lopolopen/shoot/internal/constructor"
	"github.com/lopolopen/shoot/internal/enumer"
	"github.com/lopolopen/shoot/internal/shoot"
	"github.com/lopolopen/shoot/internal/tools/logx"
)

func Gen(conf *Conf) {
	log.SetFlags(0)

	var g shoot.Generator
	for _, s := range conf.Suites {
		var args []string
		args = append(args, "shoot", s.Cmd)
		args = append(args, s.Args...)
		args = append(args, "-type="+s.Type)
		args = append(args, s.Dir)

		os.Args = args
		flag.Parse()

		switch s.Cmd {
		case constructor.SubCmd:
			g = constructor.New()
		case enumer.SubCmd:
		}

		g.ParseFlags()
		g.LoadPackage()
		srcMap := g.Generate(g)
		// logx.DebugJSON(srcMap)
		for _, src := range srcMap {
			logx.Pin(string(src))
		}
	}
}
