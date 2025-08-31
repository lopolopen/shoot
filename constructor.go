package shoot

// defaultSetter defines a type that can initialize itself with default values.
// This is useful for ensuring consistent baseline state before applying options.
type defaultSetter interface {
	SetDefault()
}

// DefaultSetter is a generic constraint for pointer types that implement SetDefault.
// It ensures that *T can be cast to PT and that PT supports default initialization.
type DefaultSetter[T any] interface {
	~*T
	defaultSetter
}

// Option defines a functional option for configuring an instance of type T.
// PT is a pointer type that satisfies DefaultSetter[T], allowing default setup.
type Option[T any, PT DefaultSetter[T]] func(*T)

// NewWith creates a new instance of type T using the functional options pattern.
// It first initializes the instance with default values via SetDefault,
// then applies each provided Option in order.
// This pattern promotes clean, declarative construction of configurable types.
func NewWith[T any, PT DefaultSetter[T]](opts ...Option[T, PT]) *T {
	t := new(T)
	(PT(t)).SetDefault()
	for _, opt := range opts {
		opt(t)
	}
	return t
}
