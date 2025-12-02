package shoot

type GenWorker interface {
	SubCmd() string

	TmplTxt() string

	Data() Data

	TypeNames() []string

	Do(typeName string) bool
}
