package utils

func Filter[T any](s []T, check func(T) bool) []T {
	res := []T{}
	for _, i := range s {
		if check(i) {
			res = append(res, i)
		}
	}
	return res
}
