package apispec

import (
	"fmt"
	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
)

func addQueryParam(operation *openapi3.Operation, key string, values []string) *openapi3.Operation {
	operation.AddParameter(openapi3.NewQueryParameter(key).WithSchema(getSchemaFromQueryValues(values)))
	return operation
}

func getSchemaFromQueryValues(values []string) *openapi3.Schema {
	var schema *openapi3.Schema
	if len(values) == 0 || values[0] == "" {
		schema = openapi3.NewBoolSchema()
		schema.AllowEmptyValue = true
	} else {
		schema = getSchemaFromValues(values, true, openapi3.ParameterInQuery)
	}
	return schema
}

func extractQueryParams(path string) (url.Values, error) {
	_, query := GetPathAndQuery(path)
	if query == "" {
		return nil, nil
	}

	values, err := url.ParseQuery(query)
	if err != nil {
		return nil, fmt.Errorf("failed to parse query: %v", err)
	}

	return values, nil
}
