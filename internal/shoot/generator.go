package shoot

type Generator interface {
	DataMaker
	TypesFilter

	ParseFlags()

	ParsePackage(typesParser TypesFilter)

	Generate(dataMaker DataMaker) map[string][]byte
}

type DataMaker interface {
	MakeData(typName string) any
}

type TypesFilter interface {
	FilterTypes() []string
}
