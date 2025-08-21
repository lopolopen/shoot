package main

import (
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"

	"golang.org/x/tools/imports"

	"github.com/lopolopen/shoot/internal/constructor"
	"github.com/lopolopen/shoot/shoot"
)

func main() {
	log.SetFlags(0)
	flag.Usage = func() {
		log.Println(`Usage: shoot <subcommand> [options]`)
		log.Println()
		log.Println(`These are all the sub commands supported as of now:`)
		log.Println()
		// log.Printf("%s [-bit] [-json]\n", enumer.SubCmd)
		// log.Println()
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

	var g shoot.Generator
	switch subCmd {
	case constructor.SubCmd:
		g = constructor.New()
	default:
		flag.Usage()
		os.Exit(2)
	}

	g.ParseFlags()
	srcMap := g.Generate()
	var fileNames []string
	for typName, src := range srcMap {
		fileName := g.FileName(typName)
		fileNames = append(fileNames, fileName)
		notedownSrc(fileName, src)
	}

	log.Println("go generate successfully:")
	for _, fn := range fileNames {
		log.Printf("\t%s\n", fn)
	}
}

func formatSrc(src []byte) []byte {
	// format imports
	src, err := imports.Process("./_.go", src, nil)
	if err != nil {
		log.Fatalf("format imports: %s", err)
	}

	// format source code
	src, err = format.Source(src)
	if err != nil {
		log.Fatalf("format source: %s", err)
	}
	return src
}

func notedownSrc(fileName string, src []byte) {
	src = formatSrc(src)

	// write to tmpfile first
	tmpFile, err := os.CreateTemp(".", fmt.Sprintf(".%s_", fileName))
	defer func() {
		if tmpFile != nil {
			_ = tmpFile.Close()
		}
	}()
	if err != nil {
		log.Fatalf("creating temporary file for output: %s", err)
	}
	_, err = tmpFile.Write(src)
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		log.Fatalf("writing output: %s", err)
	}
	tmpFile.Close()

	// rename tmpfile to output file
	err = os.Rename(tmpFile.Name(), fileName)
	if err != nil {
		log.Fatalf("moving tempfile to output file: %s", err)
	}
}
