package shoot

var restFactories = map[any]func(RestConf) RestClient{}

type restClient interface {
	RestClient
}

func Register[T restClient](ctor func(RestConf) RestClient) {
	var t T
	restFactories[t] = ctor
}

type RestClient interface {
	// SetCont(conf *RestConf)
}

type ClientPtr[T any] interface {
	~*T
	RestClient
}

//go:generate go run github.com/lopolopen/shoot/cmd/shoot new -getset -opt -type=RestConf

type RestConf struct {
	baseURL string
	timeout int
}

func NewRest[T RestClient](opts ...Option[RestConf, *RestConf]) T {
	var t T
	conf := NewWith(opts...)
	ctor, ok := restFactories[t]
	if !ok {
		panic("!!!")
	}
	t = ctor(*conf).(T)
	return t

}
