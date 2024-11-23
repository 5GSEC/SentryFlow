package apispec

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"

	"github.com/5gsec/sentryflow/speculator/pkg/util"
)

var logger = util.GetLogger()

// Note: securityDefinitions might be updated.
func (s *Spec) telemetryToOperation(telemetry *Telemetry, securitySchemes openapi3.SecuritySchemes) (*openapi3.Operation, error) {
	statusCode, err := strconv.Atoi(telemetry.Response.StatusCode)
	if err != nil {
		return nil, fmt.Errorf("failed to convert status code: %v. %v", statusCode, err)
	}

	queryParams, err := extractQueryParams(telemetry.Request.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to convert query params: %v", err)
	}

	if s.OpGenerator == nil {
		return nil, fmt.Errorf("operation generator was not set")
	}

	// Generate operation from telemetry
	telemetryOp, err := s.OpGenerator.GenerateSpecOperation(&HTTPInteractionData{
		ReqBody:     string(telemetry.Request.Common.Body),
		RespBody:    string(telemetry.Response.Common.Body),
		ReqHeaders:  ConvertHeadersToMap(telemetry.Request.Common.Headers),
		RespHeaders: ConvertHeadersToMap(telemetry.Response.Common.Headers),
		QueryParams: queryParams,
		statusCode:  statusCode,
	}, securitySchemes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate spec operation. %v", err)
	}
	return telemetryOp, nil
}

// example: for "/example-path?param=value" returns "/example-path", "param=value"
func GetPathAndQuery(fullPath string) (path, query string) {
	index := strings.IndexByte(fullPath, '?')
	if index == -1 {
		return fullPath, ""
	}

	// /path?
	if index == (len(fullPath) - 1) {
		return fullPath, ""
	}

	path = fullPath[:index]
	query = fullPath[index+1:]
	return
}
