package shoot

type NewConf struct {
}

type MapConf struct {
}

type Conf struct {
	Suites []Suite
}

type Suite struct {
	Name  string
	Cmd   string
	Args  []string
	Dir   string
	Type  string
	Types []string
	New   NewConf
	Map   MapConf
}
