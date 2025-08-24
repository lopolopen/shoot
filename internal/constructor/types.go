package constructor

import (
	"go/ast"
	"go/types"

	"github.com/lopolopen/shoot/shoot"
	"golang.org/x/tools/go/packages"
)

type Data struct {
	shoot.BaseData
	// GoFile  string
	Imports string
	//All = Exported + Unexported
	AllList     []string
	NewList     []string
	GetSet      bool
	GetterList  []string
	SetterList  []string
	Option      bool
	DefaultList []string
	Json        bool
	//Marshal: Getteer + Exported
	//Unmarshal: Setter + Exported
	ExportedList []string
	EmbedList    []string
}

type Flags struct {
	verbose   bool
	typeNames []string
	fileName  string
	//if true:
	//[ ] = [get;set] => get+set
	//[get] => get-only
	//[set] => set-only
	//if false:
	//[ ] => neither
	//[get] => get-only
	//[set] => set-only
	//[get;set] => get+set
	getset   bool
	json     bool
	opt      bool
	separate bool
}

// Value represents a declared constant.
type Value struct {
	originalName string // The name of the constant before transformation
	name         string // The name of the constant after transformation (i.e. camel case => snake case)
	// The value is stored as a bit pattern alone. The boolean tells us
	// whether to interpret it as an int64 or a uint64; the only place
	// this matters is when sorting.
	// Much of the time the str field is all we need; it is printed
	// by Value.String.
	value  uint64 // Will be converted to int64 when needed.
	signed bool   // Whether the constant is a signed type.
	str    string // The string representation given by the "go/exact" package.
}

// File holds a single parsed file and associated data.
type File struct {
	pkg  *Package  // Package to which this file belongs.
	file *ast.File // Parsed AST.
	// These fields are reset for each type being generated.
	typeName    string  // Name of the constant type.
	values      []Value // Accumulator for constant values of that type.
	trimPrefix  string
	lineComment bool
}

// Package holds information about a Go package
type Package struct {
	pkg      *packages.Package
	dir      string
	name     string
	defs     map[*ast.Ident]types.Object
	files    []*File
	typesPkg *types.Package
}
