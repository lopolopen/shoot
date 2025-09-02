package shoot

import (
	"fmt"
	"net/http"
	"reflect"
	"time"
)

//go:generate go run github.com/lopolopen/shoot/cmd/shoot new -getset -opt -type=RestConf

type Middleware func(next http.RoundTripper) http.RoundTripper

// RestConf holds configuration parameters for initializing a RestClient.
type RestConf struct {
	baseURL        string
	timeout        time.Duration
	defaultHeaders map[string]string
	Middlewares    []Middleware
}

func (r *RestConf) use(m Middleware) *RestConf {
	r.Middlewares = append(r.Middlewares, m)
	return r
}

func Use(middleware Middleware) Option[RestConf, *RestConf] {
	return func(r *RestConf) {
		r.use(middleware)
	}
}

func (r *RestConf) BuildMiddleware() http.RoundTripper {
	t := http.DefaultTransport
	for i := len(r.Middlewares) - 1; i >= 0; i-- {
		t = r.Middlewares[i](t)
	}
	return t
}

type RoundTripper func(*http.Request) (*http.Response, error)

func (f RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// RestClient defines the interface for REST-capable clients.
type RestClient[T any] interface {
	ConfigHTTPClient(config func(client *http.Client)) T
	ShootRest()
}

// Register associates a constructor function with a specific RestClient implementation.
// The constructor must accept a RestConf and return a concrete type that satisfies RestClient.
func Register[T RestClient[T]](ctor func(RestConf) T) {
	typ := reflect.TypeOf((*T)(nil))
	_, ok := ctorRegistry[typ]
	if ok {
		panic(fmt.Errorf("ctor of interface %s should not be registered multiple times", typ.Elem()))
	}
	ctorRegistry[typ] = ctor
}

// NewRest creates a new instance of type T that implements the RestClient interface.
func NewRest[T RestClient[T]](opts ...Option[RestConf, *RestConf]) T {
	typ := reflect.TypeOf((*T)(nil))
	conf := NewWith(opts...)
	ctor, ok := ctorRegistry[typ]
	if !ok {
		panic(fmt.Errorf("ctor of interface %s is not regstered", typ.Elem()))
	}
	typedCtor, ok := ctor.(func(RestConf) T)
	if !ok {
		panic(fmt.Errorf("registered ctor of interface %s has invalid type %T", typ.Elem(), ctor))
	}
	t := typedCtor(*conf)
	return t
}

var ctorRegistry = map[reflect.Type]any{}
