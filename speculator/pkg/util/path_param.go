package util

import (
	"strings"
)

const (
	ParamPrefix = "{"
	ParamSuffix = "}"
)

func IsPathParam(segment string) bool {
	return strings.HasPrefix(segment, ParamPrefix) &&
		strings.HasSuffix(segment, ParamSuffix)
}
