package apispec

import (
	"encoding/json"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

type ApprovedSpec struct {
	PathItems       map[string]*openapi3.PathItem
	SecuritySchemes openapi3.SecuritySchemes
	SpecVersion     OASVersion
}

func (a *ApprovedSpec) GetPathItem(path string) *openapi3.PathItem {
	if pi, exists := a.PathItems[path]; exists {
		return pi
	}
	return nil
}

func (a *ApprovedSpec) GetSpecVersion() OASVersion {
	return a.SpecVersion
}

func (a *ApprovedSpec) Clone() (*ApprovedSpec, error) {
	clonedApprovedSpec := new(ApprovedSpec)

	approvedSpecB, err := json.Marshal(a)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal approved spec: %w", err)
	}

	if err := json.Unmarshal(approvedSpecB, &clonedApprovedSpec); err != nil {
		return nil, fmt.Errorf("failed to unmarshal approved spec: %w", err)
	}

	return clonedApprovedSpec, nil
}
