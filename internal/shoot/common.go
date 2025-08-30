package shoot

func Contains[T comparable](slice []T, val T) bool {
	for _, x := range slice {
		if x == val {
			return true
		}
	}
	return false
}
