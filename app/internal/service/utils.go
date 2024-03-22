package service

import (
	"errors"
	"strings"
)

var (
	ErrNotFound = errors.New("not found")
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

func parseCommand(cmd string) []any {
	res := []any{}

	cmd = strings.ReplaceAll(cmd, "\n", " ")
	cmd = strings.ReplaceAll(cmd, "\t", " ")
	splitCmd := strings.Split(cmd, " ")
	splitCmd = Filter(splitCmd, func(i string) bool {
		return i != ""
	})
	for _, i := range splitCmd {
		res = append(res, i)
	}

	return res
}
