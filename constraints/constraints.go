package constraints

type Integer interface {
	~int8 | ~int16 | ~int32 | int | ~int64 |
		~uint8 | ~uint16 | ~uint32 | ~uint | ~uint64
}

// NewShooter is a generic constraint for pointer types that implement SetDefault.
// It ensures that *T can be cast to PT and that PT supports default initialization.
type NewShooter[T any] interface {
	~*T
	ShootNew()
}

type EnumShooter[T any] interface {
	Integer
	ShootEnum()
	Values() []T
	ValueMap() map[string]T
}

type RestShooter interface {
	ShootRest()
}

type MapShooter interface {
	ShootMap()
}
