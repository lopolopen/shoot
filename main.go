package main

import (
	"flag"
	"log"
	"os"

	"github.com/lopolopen/shoot/internal/constructor"
	"github.com/lopolopen/shoot/internal/enumer"
)

func main() {
	log.SetFlags(0)
	flag.Usage = func() {
		log.Println(`Usage: shoot <subcommand> [options]`)
		log.Println()
		log.Println(`These are all the sub commands supported as of now:`)
		log.Println()
		log.Printf("%s [-bit] [-json]\n", enumer.SubCmd)
		log.Println()
		log.Printf("%s [-new] [-opt] [-getset]  [-json]\n", constructor.SubCmd)
		log.Println()
		flag.PrintDefaults()
	}
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	subCmd := flag.Args()[0]
	switch subCmd {
	case enumer.SubCmd:
		enumer.Gen()
	case constructor.SubCmd:
		constructor.Gen()
	default:
		flag.Usage()
		os.Exit(2)
	}

}
