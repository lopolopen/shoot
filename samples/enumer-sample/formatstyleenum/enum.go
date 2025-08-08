package formatstyleenum

//go:generate go tool shoot enum -bit -type=Enum

type Enum int

const (
	None Enum = 0
	Bold Enum = 1 << iota
	Italic
	Underline
	Strikethrough
)
