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
