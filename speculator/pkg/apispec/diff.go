package apispec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofrs/uuid"
)

type DiffType string

const (
	DiffTypeNoDiff      DiffType = "NO_DIFF"
	DiffTypeZombieDiff  DiffType = "ZOMBIE_DIFF"
	DiffTypeShadowDiff  DiffType = "SHADOW_DIFF"
	DiffTypeGeneralDiff DiffType = "GENERAL_DIFF"
)

type APIDiff struct {
	Type             DiffType
	Path             string
	OriginalPathItem *openapi3.PathItem
	ModifiedPathItem *openapi3.PathItem
	InteractionID    uuid.UUID
	SpecID           uuid.UUID
}

type operationDiff struct {
	OriginalOperation *openapi3.Operation
	ModifiedOperation *openapi3.Operation
}

type DiffParams struct {
	operation *openapi3.Operation
	method    string
	path      string
	requestID string
	response  *Response
}

func (s *Spec) createDiffParamsFromTelemetry(telemetry *Telemetry) (*DiffParams, error) {
	securitySchemes := openapi3.SecuritySchemes{}

	path, _ := GetPathAndQuery(telemetry.Request.Path)
	telemetryOp, err := s.telemetryToOperation(telemetry, securitySchemes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert telemetry to operation: %w", err)
	}
	return &DiffParams{
		operation: telemetryOp,
		method:    telemetry.Request.Method,
		path:      path,
		requestID: telemetry.RequestID,
		response:  telemetry.Response,
	}, nil
}

func (s *Spec) DiffTelemetry(telemetry *Telemetry, specSource SpecSource) (*APIDiff, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	var apiDiff *APIDiff
	var err error
	diffParams, err := s.createDiffParamsFromTelemetry(telemetry)
	if err != nil {
		return nil, fmt.Errorf("failed to create diff params from telemetry. %w", err)
	}

	switch specSource {
	case SpecSourceProvided:
		if !s.HasProvidedSpec() {
			logger.Infof("No provided spec to diff")
			return nil, nil
		}
		apiDiff, err = s.diffProvidedSpec(diffParams)
		if err != nil {
			return nil, fmt.Errorf("failed to diff provided spec. %w", err)
		}
	case SpecSourceReconstructed:
		if !s.HasApprovedSpec() {
			logger.Infof("No approved spec to diff")
			return nil, nil
		}
		apiDiff, err = s.diffApprovedSpec(diffParams)
		if err != nil {
			return nil, fmt.Errorf("failed to diff approved spec. %w", err)
		}
	default:
		return nil, fmt.Errorf("spec source: %v is not valid", specSource)
	}

	return apiDiff, nil
}

func (s *Spec) diffApprovedSpec(diffParams *DiffParams) (*APIDiff, error) {
	var pathItem *openapi3.PathItem
	pathFromTrie, _, found := s.ApprovedPathTrie.GetPathAndValue(diffParams.path)
	if found {
		diffParams.path = pathFromTrie // The diff will show the parametrized path if matched and not the telemetry path
		pathItem = s.ApprovedSpec.GetPathItem(pathFromTrie)
	}
	return s.diffPathItem(pathItem, diffParams)
}

func (s *Spec) diffProvidedSpec(diffParams *DiffParams) (*APIDiff, error) {
	var pathItem *openapi3.PathItem

	basePath := s.ProvidedSpec.GetBasePath()

	pathNoBase := trimBasePathIfNeeded(basePath, diffParams.path)

	pathFromTrie, _, found := s.ProvidedPathTrie.GetPathAndValue(pathNoBase)
	if found {
		// The diff will show the parametrized path if matched and not the telemetry path
		diffParams.path = addBasePathIfNeeded(basePath, pathFromTrie)
		pathItem = s.ProvidedSpec.GetPathItem(pathFromTrie)
	}

	return s.diffPathItem(pathItem, diffParams)
}

// For path /api/foo/bar and base path of /api, the path that will be saved in paths map will be /foo/bar
// All paths must start with a slash. We can't trim a leading slash.
func trimBasePathIfNeeded(basePath, path string) string {
	if hasBasePath(basePath) {
		return strings.TrimPrefix(path, basePath)
	}

	return path
}

func addBasePathIfNeeded(basePath, path string) string {
	if hasBasePath(basePath) {
		return basePath + path
	}

	return path
}

func hasBasePath(basePath string) bool {
	return basePath != "" && basePath != "/"
}

func (s *Spec) diffPathItem(pathItem *openapi3.PathItem, diffParams *DiffParams) (*APIDiff, error) {
	var apiDiff *APIDiff
	method := diffParams.method
	telemetryOp := diffParams.operation
	path := diffParams.path
	requestID := diffParams.requestID
	reqUUID := uuid.NewV5(uuid.Nil, requestID)

	if pathItem == nil {
		apiDiff = s.createAPIDiffEvent(DiffTypeShadowDiff, nil, createPathItemFromOperation(method, telemetryOp),
			reqUUID, path)
		return apiDiff, nil
	}

	specOp := GetOperationFromPathItem(pathItem, method)
	if specOp == nil {
		// new operation
		apiDiff := s.createAPIDiffEvent(DiffTypeShadowDiff, pathItem, CopyPathItemWithNewOperation(pathItem, method, telemetryOp),
			reqUUID, path)
		return apiDiff, nil
	}

	diff, err := calculateOperationDiff(specOp, telemetryOp, diffParams.response)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate operation diff: %w", err)
	}
	if diff != nil {
		diffType := DiffTypeGeneralDiff
		if specOp.Deprecated {
			diffType = DiffTypeZombieDiff
		}
		apiDiff := s.createAPIDiffEvent(diffType, createPathItemFromOperation(method, diff.OriginalOperation),
			createPathItemFromOperation(method, diff.ModifiedOperation), reqUUID, path)
		return apiDiff, nil
	}

	// no diff
	return s.createAPIDiffEvent(DiffTypeNoDiff, nil, nil, reqUUID, path), nil
}

func (s *Spec) createAPIDiffEvent(diffType DiffType, original, modified *openapi3.PathItem, interactionID uuid.UUID, path string) *APIDiff {
	return &APIDiff{
		Type:             diffType,
		Path:             path,
		OriginalPathItem: original,
		ModifiedPathItem: modified,
		InteractionID:    interactionID,
		SpecID:           s.ID,
	}
}

func createPathItemFromOperation(method string, operation *openapi3.Operation) *openapi3.PathItem {
	pathItem := openapi3.PathItem{}
	AddOperationToPathItem(&pathItem, method, operation)
	return &pathItem
}

func calculateOperationDiff(specOp, telemetryOp *openapi3.Operation, telemetryResponse *Response) (*operationDiff, error) {
	clonedTelemetryOp, err := CloneOperation(telemetryOp)
	if err != nil {
		return nil, fmt.Errorf("failed to clone telemetry operation: %w", err)
	}

	clonedSpecOp, err := CloneOperation(specOp)
	if err != nil {
		return nil, fmt.Errorf("failed to clone spec operation: %w", err)
	}

	clonedTelemetryOp = sortParameters(clonedTelemetryOp)
	clonedSpecOp = sortParameters(clonedSpecOp)

	// Keep only telemetry status code
	clonedSpecOp = keepResponseStatusCode(clonedSpecOp, telemetryResponse.StatusCode)

	hasDiff, err := compareObjects(clonedSpecOp, clonedTelemetryOp)
	if err != nil {
		return nil, fmt.Errorf("failed to compare operations: %w", err)
	}

	if hasDiff {
		return &operationDiff{
			OriginalOperation: clonedSpecOp,
			ModifiedOperation: clonedTelemetryOp,
		}, nil
	}

	// no diff
	return nil, nil
}

func compareObjects(obj1, obj2 any) (hasDiff bool, err error) {
	obj1B, err := json.Marshal(obj1)
	if err != nil {
		return false, fmt.Errorf("failed to marshal obj1: %w", err)
	}

	obj2B, err := json.Marshal(obj2)
	if err != nil {
		return false, fmt.Errorf("failed to marshal obj2: %w", err)
	}

	return !bytes.Equal(obj1B, obj2B), nil
}

// keepResponseStatusCode will remove all status codes from StatusCodeResponses map except the `statusCodeToKeep`.
func keepResponseStatusCode(op *openapi3.Operation, statusCodeToKeep string) *openapi3.Operation {
	// keep only the provided status code
	if op.Responses != nil {
		filteredResponses := &openapi3.Responses{}
		if responseRef := op.Responses.Value(statusCodeToKeep); responseRef != nil {
			filteredResponses.Set(statusCodeToKeep, responseRef)
		}
		// keep default if exists
		if responseRef := op.Responses.Value("default"); responseRef != nil {
			filteredResponses.Set("default", responseRef)
		}

		if filteredResponses.Len() == 0 {
			op.Responses = nil
		} else {
			op.Responses = filteredResponses
		}
	}

	return op
}

func sortParameters(operation *openapi3.Operation) *openapi3.Operation {
	if operation == nil {
		return operation
	}
	sort.Slice(operation.Parameters, func(i, j int) bool {
		right := operation.Parameters[i].Value
		left := operation.Parameters[j].Value
		// Sibling parameters must have unique name + in values
		return right.Name+right.In < left.Name+left.In
	})

	return operation
}
