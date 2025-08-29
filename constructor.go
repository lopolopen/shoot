package shoot

type DefaultSetter interface {
	SetDefault()
}

type InstancePtr[T any] interface {
	~*T
	SetDefault()
}

type Option[T any, PT InstancePtr[T]] func(*T)

// NewWith constructs a new instance of type T using the functional options pattern
func NewWith[T any, PT InstancePtr[T]](opts ...Option[T, PT]) *T {
	t := new(T)
	(PT(t)).SetDefault()
	for _, opt := range opts {
		opt(t)
	}
	return t
}
