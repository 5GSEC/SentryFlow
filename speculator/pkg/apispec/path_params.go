package apispec

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofrs/uuid"
)

type PathParam struct {
	*openapi3.Parameter
}

func generateParamName(i int) string {
	return fmt.Sprintf("param%v", i)
}

var digitCheck = regexp.MustCompile(`^[0-9]+$`)

func createParameterizedPath(path string) string {
	var ParameterizedPathParts []string
	paramCount := 0
	pathParts := strings.Split(path, "/")

	for _, part := range pathParts {
		// if part is a suspect param, replace it with a param name, otherwise do nothing
		if isSuspectPathParam(part) {
			paramCount++
			paramName := generateParamName(paramCount)
			ParameterizedPathParts = append(ParameterizedPathParts, "{"+paramName+"}")
		} else {
			ParameterizedPathParts = append(ParameterizedPathParts, part)
		}
	}

	parameterizedPath := strings.Join(ParameterizedPathParts, "/")

	return parameterizedPath
}

type paramFormat string

const (
	paramFormatUnset  paramFormat = "paramFormatUnset"
	paramFormatNumber paramFormat = "paramFormatNumber"
	paramFormatUUID   paramFormat = "paramFormatUUID"
	paramFormatMixed  paramFormat = "paramFormatMixed"
)

// /api/1/foo, api/2/foo and index 1 will return:
// []string{1, 2}.
func getOnlyIndexedPartFromPaths(paths map[string]bool, i int) []string {
	var ret []string
	for path := range paths {
		path = strings.TrimPrefix(path, "/")
		splt := strings.Split(path, "/")
		if len(splt) <= i {
			continue
		}
		ret = append(ret, splt[i])
	}
	return ret
}

// If all params in paramList can be guessed as same schema, this schema will be returned, otherwise,
// if there is a couple of formats, string schema with no format will be returned.
func getParamSchema(paramsList []string) *openapi3.Schema {
	parameterFormat := paramFormatUnset

	for _, pathPart := range paramsList {
		if isNumber(pathPart) {
			// in case there is a conflict, we will return string as the type and empty format
			if parameterFormat != paramFormatNumber && parameterFormat != paramFormatUnset {
				return openapi3.NewStringSchema()
			}
			parameterFormat = paramFormatNumber
			continue
		}
		if isUUID(pathPart) {
			if parameterFormat != paramFormatUUID && parameterFormat != paramFormatUnset {
				return openapi3.NewStringSchema()
			}
			parameterFormat = paramFormatUUID
			continue
		}
		if isMixed(pathPart) {
			if parameterFormat != paramFormatMixed && parameterFormat != paramFormatUnset {
				return openapi3.NewStringSchema()
			}
			parameterFormat = paramFormatMixed
		}
	}

	switch parameterFormat {
	case paramFormatMixed:
		return openapi3.NewStringSchema()
	case paramFormatUUID:
		return openapi3.NewUUIDSchema()
	case paramFormatNumber:
		return openapi3.NewInt64Schema()
	case paramFormatUnset:
		return openapi3.NewStringSchema()
	}

	return openapi3.NewStringSchema()
}

func isSuspectPathParam(pathPart string) bool {
	if isNumber(pathPart) {
		return true
	}
	if isUUID(pathPart) {
		return true
	}
	if isMixed(pathPart) {
		return true
	}
	return false
}

func isNumber(pathPart string) bool {
	return digitCheck.MatchString(pathPart)
}

func isUUID(pathPart string) bool {
	_, err := uuid.FromString(pathPart)
	return err == nil
}

// Check if a path part that is mixed from digits and chars can be considered as parameter following hard-coded heuristics.
// Temporary, we'll consider strings as parameters that are at least 8 chars longs and has at least 3 digits.
func isMixed(pathPart string) bool {
	const maxLen = 8
	const minDigitsLen = 2

	if len(pathPart) < maxLen {
		return false
	}

	return countDigitsInString(pathPart) > minDigitsLen
}

func countDigitsInString(s string) int {
	count := 0
	for _, c := range s {
		if unicode.IsNumber(c) {
			count++
		}
	}
	return count
}

func createPathParam(name string, schema *openapi3.Schema) *PathParam {
	return &PathParam{
		Parameter: openapi3.NewPathParameter(name).WithSchema(schema),
	}
}
