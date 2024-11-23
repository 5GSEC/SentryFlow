package apispec

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"k8s.io/utils/field"
)

type conflict struct {
	path *field.Path
	obj1 any
	obj2 any
	msg  string
}

func createConflictMsg(path *field.Path, t1, t2 any) string {
	return fmt.Sprintf("%s: type mismatch: %+v != %+v", path, t1, t2)
}

func createHeaderInConflictMsg(path *field.Path, in, in2 any) string {
	return fmt.Sprintf("%s: header in mismatch: %+v != %+v", path, in, in2)
}

func (c conflict) String() string {
	return c.msg
}

const (
	NoConflict = iota
	PreferType1
	PreferType2
	ConflictUnresolved
)

// conflictSolver will get 2 types and returns:
//
//	NoConflict - type1 and type2 are equal
//	PreferType1 - type1 should be used
//	PreferType2 - type2 should be used
//	ConflictUnresolved - types conflict can't be resolved
func conflictSolver(type1, type2 *openapi3.Types) int {
	if type1.Is(type2.Slice()[0]) {
		return NoConflict
	}

	if shouldPreferType(type1, type2) {
		return PreferType1
	}

	if shouldPreferType(type2, type1) {
		return PreferType2
	}

	return ConflictUnresolved
}

// shouldPreferType return true if type1 should be preferred over type2.
func shouldPreferType(type1, type2 *openapi3.Types) bool {
	if type1.Includes(openapi3.TypeBoolean) ||
		type1.Includes(openapi3.TypeObject) ||
		type1.Includes(openapi3.TypeArray) {
		// Should not prefer boolean, object and array type over any other type.
		return false
	}

	if type1.Includes(openapi3.TypeNumber) {
		// Preferring number to integer type.
		return type2.Includes(openapi3.TypeInteger)
	}

	if type1.Includes(openapi3.TypeString) {
		// Preferring string to any type.
		return true
	}

	return false
}
