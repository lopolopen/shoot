package shoot

type Generator interface {
	ParseFlags()

	Generate() map[string][]byte

	FileName(typeName string) string
}
