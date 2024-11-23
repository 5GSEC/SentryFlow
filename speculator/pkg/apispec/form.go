package apispec

import (
	"fmt"
	"mime/multipart"
	"net/url"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

const (
	// taken from net/http/request.go.
	defaultMaxMemory = 32 << 20 // 32 MB
)

func handleApplicationFormURLEncodedBody(operation *openapi3.Operation, securitySchemes openapi3.SecuritySchemes, body string) (*openapi3.Operation, openapi3.SecuritySchemes, error) {
	parseQuery, err := url.ParseQuery(body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse query. body=%v: %v", body, err)
	}

	objSchema := openapi3.NewObjectSchema()

	for key, values := range parseQuery {
		if key == AccessTokenParamKey {
			// https://datatracker.ietf.org/doc/html/rfc6750#section-2.2
			operation, securitySchemes = handleAuthQueryParam(operation, securitySchemes, values)
		} else {
			objSchema.WithProperty(key, getSchemaFromQueryValues(values))
		}
	}

	if len(objSchema.Properties) != 0 {
		operationSetRequestBody(operation, openapi3.NewRequestBody().WithContent(openapi3.NewContentWithSchema(objSchema, []string{mediaTypeApplicationForm})))
		// TODO: handle encoding
		// https://swagger.io/docs/specification/describing-request-body/
		// operation.RequestBody.Value.GetMediaType(mediaTypeApplicationForm).Encoding
	}

	return operation, securitySchemes, nil
}

func getMultipartFormDataSchema(body string, mediaTypeParams map[string]string) (*openapi3.Schema, error) {
	boundary, ok := mediaTypeParams["boundary"]
	if !ok {
		return nil, fmt.Errorf("no multipart boundary param in Content-Type")
	}

	form, err := multipart.NewReader(strings.NewReader(body), boundary).ReadForm(defaultMaxMemory)
	if err != nil {
		return nil, fmt.Errorf("failed to read form: %w", err)
	}

	schema := openapi3.NewObjectSchema()

	// https://swagger.io/docs/specification/describing-request-body/file-upload/
	for key, fileHeaders := range form.File {
		fileSchema := openapi3.NewStringSchema().WithFormat("binary")
		switch len(fileHeaders) {
		case 0:
			// do nothing
		case 1:
			// single file
			schema.WithProperty(key, fileSchema)
		default:
			// array of files
			schema.WithProperty(key, openapi3.NewArraySchema().WithItems(fileSchema))
		}
	}

	// add values formData
	for key, values := range form.Value {
		schema.WithProperty(key, getSchemaFromValues(values, false, ""))
	}

	return schema, nil
}
