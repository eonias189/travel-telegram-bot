package service

import (
	"errors"
	"strings"

	"github.com/Central-University-IT-prod/backend-eonias189/internal/lib/utils"
)

var (
	ErrNotFound = errors.New("not found")
)

func parseCommand(cmd string) []any {
	res := []any{}
	cmd = strings.ReplaceAll(cmd, "\n", " ")
	cmd = strings.ReplaceAll(cmd, "\t", " ")
	splitCmd := strings.Split(cmd, " ")
	splitCmd = utils.Filter(splitCmd, func(i string) bool {
		return i != ""
	})
	for _, i := range splitCmd {
		res = append(res, i)
	}

	return res
}
