package model

//go:generate go tool shoot new -getset -opt -file=$GOFILE

type Conf struct {
	//shoot: new
	name string
	host []string //todo: shoot: def=host1,host2
	//shoot: def=80
	port int
}
