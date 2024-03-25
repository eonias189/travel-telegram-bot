package utils

import (
	"errors"
	"net/url"
	"strconv"
)

var (
	ErrConvertaion = errors.New("ConvertaionError")
)

func Filter[T any](s []T, check func(T) bool) []T {
	res := []T{}
	for _, i := range s {
		if check(i) {
			res = append(res, i)
		}
	}
	return res
}

func GetInt(query url.Values, key string) (int, error) {
	val := query.Get(key)
	if val == "" {
		return 0, ErrConvertaion
	}

	res, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}

	return res, nil
}

func GetInt64(query url.Values, key string) (int64, error) {
	res, err := GetInt(query, key)
	return int64(res), err
}

func Map[T any, R any](s []T, mod func(T) R) []R {
	res := make([]R, len(s))
	for i, item := range s {
		res[i] = mod(item)
	}
	return res
}

type num interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

func Min[T num](s []T) T {
	if len(s) == 0 {
		var t T
		return t
	}
	res := s[0]
	for _, i := range s {
		if i < res {
			res = i
		}
	}
	return res
}
func Max[T num](s []T) T {
	if len(s) == 0 {
		var t T
		return t
	}
	res := s[0]
	for _, i := range s {
		if i > res {
			res = i
		}
	}
	return res
}
