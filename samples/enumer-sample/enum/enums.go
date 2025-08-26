package enum

//go:generate go tool shoot enum -file=$GOFILE

type Light int32

const (
	LightRed Light = iota
	LightYello
	LightGreen
)
