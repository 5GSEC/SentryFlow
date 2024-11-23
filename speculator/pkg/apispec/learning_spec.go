package apispec

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type LearningSpec struct {
	// map parameterized path into path item
	PathItems       map[string]*openapi3.PathItem
	SecuritySchemes openapi3.SecuritySchemes
}

func (l *LearningSpec) AddPathItem(path string, pathItem *openapi3.PathItem) {
	l.PathItems[path] = pathItem
}

func (l *LearningSpec) GetPathItem(path string) *openapi3.PathItem {
	pi, ok := l.PathItems[path]
	if !ok {
		return nil
	}

	return pi
}
