package shoot

type Generator interface {
	DataMaker
	TypeLister

	ParseFlags()

	LoadPackage()

	ParsePackage(typeLister TypeLister)

	Generate(dataMaker DataMaker) map[string][]byte
}

type DataMaker interface {
	MakeData(typName string) any
}

type TypeLister interface {
	ListTypes() []string
}
