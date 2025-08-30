package shoot

import (
	"fmt"
)

type Enumer[T any] interface {
	~int | ~int32
	Values() []T
	ValueMap() map[string]T
}

// ParseEnum returns the enum value represented by the string.
func ParseEnum[T Enumer[T]](str string) (T, error) {
	var t T
	m := t.ValueMap()
	t, ok := m[str]
	if ok {
		return t, nil
	}
	return t, fmt.Errorf("requested value '%s' was not found", str)
}

func IsEnum[T Enumer[T], TV ~int | ~int32](value TV) bool {
	var t T
	for _, v := range t.Values() {
		if v == T(value) {
			return true
		}
	}
	return false
}
