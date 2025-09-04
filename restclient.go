package shoot

import "github.com/lopolopen/shoot/constraints"

var restFactories = map[any]func(RestConf) any{}

func Register[T constraints.ShootRest[T]](ctor func(RestConf) any) {
	var t T
	restFactories[t] = ctor
}

//go:generate go run github.com/lopolopen/shoot/cmd/shoot new -getset -opt -type=RestConf

type RestConf struct {
	baseURL string
	timeout int
}

func NewRest[T constraints.ShootRest[T]](opts ...Option[RestConf, *RestConf]) T {
	var t T
	conf := NewWith(opts...)
	ctor, ok := restFactories[t]
	if !ok {
		panic("non rest constructor registered")
	}
	t = ctor(*conf).(T)
	return t

}
