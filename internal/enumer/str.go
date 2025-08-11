package enumer

import (
	"bytes"
	"fmt"
	"log"

	"github.com/lopolopen/shoot/shoot"
)

// usize returns the number of bits of the smallest unsigned integer
// type that will hold n. Used to create the smallest possible slice of
// integers to use as indexes into the concatenated strings.
func usize(n int) int {
	switch {
	case n < 1<<8:
		return 8
	case n < 1<<16:
		return 16
	default:
		// 2^32 is enough constants for anyone.
		return 32
	}
}

// createIndexAndNameDecl returns the pair of declarations for the run.
func createIndexAndNameDecl(run []shoot.Value, typeName string, suffix string) (string, string) {
	if len(run) == 0 {
		log.Fatalln("empty run")
	}
	b := new(bytes.Buffer)
	indexes := make([]int, len(run))
	for i := range run {
		b.WriteString(run[i].Name())
		indexes[i] = b.Len()
	}
	nameConst := fmt.Sprintf("_%sName%s = %q", typeName, suffix, b.String())
	nameLen := b.Len()
	b.Reset()
	_, _ = fmt.Fprintf(b, "_%sIndex%s = [...]uint%d{0, ", typeName, suffix, usize(nameLen))
	for i, v := range indexes {
		if i > 0 {
			_, _ = fmt.Fprintf(b, ", ")
		}
		_, _ = fmt.Fprintf(b, "%d", v)
	}
	_, _ = fmt.Fprintf(b, "}")
	return b.String(), nameConst
}
