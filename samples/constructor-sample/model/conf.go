package model

//go:generate go tool shoot new -opt -file=$GOFILE

type Conf struct {
	//error: //shoot: def="1"
	name string
	host string
	port int
}

// shoot: ignore
type Conf2 struct {
	//shoot: new
	name string
	//shoot: def="localhost"
	host string
	//shoot: def=80
	port int
}
