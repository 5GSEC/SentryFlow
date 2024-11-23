package util

import (
	"reflect"
)

func IsNil(a any) bool {
	return a == nil || (reflect.ValueOf(a).Kind() == reflect.Ptr && reflect.ValueOf(a).IsNil())
}
