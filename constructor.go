package shoot

// DefaultSetter defines a type that can initialize itself with default values.
// This is useful for ensuring consistent baseline state before applying options.
type DefaultSetter interface {
	SetDefault()
}

// InstancePtr is a generic constraint for pointer types that implement SetDefault.
// It ensures that *T can be cast to PT and that PT supports default initialization.
type InstancePtr[T any] interface {
	~*T
	SetDefault()
}

// Option defines a functional option for configuring an instance of type T.
// PT is a pointer type that satisfies InstancePtr[T], allowing default setup.
type Option[T any, PT InstancePtr[T]] func(*T)

// NewWith creates a new instance of type T using the functional options pattern.
// It first initializes the instance with default values via SetDefault,
// then applies each provided Option in order.
// This pattern promotes clean, declarative construction of configurable types.
func NewWith[T any, PT InstancePtr[T]](opts ...Option[T, PT]) *T {
	t := new(T)
	(PT(t)).SetDefault()
	for _, opt := range opts {
		opt(t)
	}
	return t
}
