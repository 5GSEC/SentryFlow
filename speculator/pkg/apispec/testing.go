package apispec

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

var req1 = `{"active":true,
"certificateVersion":"86eb5278-676a-3b7c-b29d-4a57007dc7be",
"controllerInstanceInfo":{"replicaId":"portshift-agent-66fc77c848-tmmk8"},
"policyAndAppVersion":1621477900361,
"version":"1.147.1"}`

var res1 = `{"cvss":[{"score":7.8,"vector":"AV:L/AC:L/PR:N/UI:R/S:U/C:H/I:H/A:H"}]}`

var req2 = `{"active":true,"statusCodes":["NO_METRICS_SERVER"],"version":"1.147.1"}`

var res2 = `{"cvss":[{"version":"3"}]}`

var combinedReq = `{"active":true,"statusCodes":["NO_METRICS_SERVER"],
"certificateVersion":"86eb5278-676a-3b7c-b29d-4a57007dc7be",
"controllerInstanceInfo":{"replicaId":"portshift-agent-66fc77c848-tmmk8"},
"policyAndAppVersion":1621477900361,
"version":"1.147.1"}`

var combinedRes = `{"cvss":[{"score":7.8,"vector":"AV:L/AC:L/PR:N/UI:R/S:U/C:H/I:H/A:H","version":"3"}]}`

type TestSpec struct {
	Doc *openapi3.T
}

func (t *TestSpec) WithPathItem(path string, pathItem *openapi3.PathItem) *TestSpec {
	t.Doc.Paths.Set(path, pathItem)
	return t
}

type TestPathItem struct {
	PathItem openapi3.PathItem
}

func NewTestPathItem() *TestPathItem {
	return &TestPathItem{
		PathItem: openapi3.PathItem{},
	}
}

func (t *TestPathItem) WithPathParams(name string, schema *openapi3.Schema) *TestPathItem {
	pathParam := createPathParam(name, schema)
	t.PathItem.Parameters = append(t.PathItem.Parameters, &openapi3.ParameterRef{Value: pathParam.Parameter})
	return t
}

func (t *TestPathItem) WithOperation(method string, op *openapi3.Operation) *TestPathItem {
	switch method {
	case http.MethodGet:
		t.PathItem.Get = op
	case http.MethodDelete:
		t.PathItem.Delete = op
	case http.MethodOptions:
		t.PathItem.Options = op
	case http.MethodPatch:
		t.PathItem.Patch = op
	case http.MethodHead:
		t.PathItem.Head = op
	case http.MethodPost:
		t.PathItem.Post = op
	case http.MethodPut:
		t.PathItem.Put = op
	}
	return t
}

type TestOperation struct {
	Op *openapi3.Operation
}

func NewOperation(t *testing.T, data *HTTPInteractionData) *TestOperation {
	t.Helper()
	securitySchemes := openapi3.SecuritySchemes{}
	operation, err := CreateTestNewOperationGenerator().GenerateSpecOperation(data, securitySchemes)
	if err != nil {
		t.Fatal(err)
	}
	return &TestOperation{
		Op: operation,
	}
}

func CreateTestNewOperationGenerator() *OperationGenerator {
	return NewOperationGenerator(testOperationGeneratorConfig)
}

var testOperationGeneratorConfig = OperationGeneratorConfig{
	ResponseHeadersToIgnore: []string{contentTypeHeaderName},
	RequestHeadersToIgnore:  []string{acceptTypeHeaderName, authorizationTypeHeaderName, contentTypeHeaderName},
}

func (op *TestOperation) Deprecated() *TestOperation {
	op.Op.Deprecated = true
	return op
}

func (op *TestOperation) WithResponse(status int, response *openapi3.Response) *TestOperation {
	op.Op.AddResponse(status, response)
	if status != 0 {
		// we don't need it to create default response in tests unless we explicitly asked for (status == 0)
		delete(op.Op.Responses.Map(), "default")
	}
	return op
}

func (op *TestOperation) WithParameter(param *openapi3.Parameter) *TestOperation {
	op.Op.AddParameter(param)
	return op
}

func (op *TestOperation) WithRequestBody(requestBody *openapi3.RequestBody) *TestOperation {
	operationSetRequestBody(op.Op, requestBody)
	return op
}

func (op *TestOperation) WithSecurityRequirement(securityRequirement openapi3.SecurityRequirement) *TestOperation {
	if op.Op.Security == nil {
		op.Op.Security = openapi3.NewSecurityRequirements()
	}
	op.Op.Security.With(securityRequirement)
	return op
}

func createTestOperation() *TestOperation {
	return &TestOperation{Op: openapi3.NewOperation()}
}

type TestResponse struct {
	*openapi3.Response
}

func createTestResponse() *TestResponse {
	return &TestResponse{
		Response: openapi3.NewResponse(),
	}
}

func (r *TestResponse) WithHeader(name string, schema *openapi3.Schema) *TestResponse {
	if r.Response.Headers == nil {
		r.Response.Headers = make(openapi3.Headers)
	}
	r.Response.Headers[name] = &openapi3.HeaderRef{
		Value: &openapi3.Header{
			Parameter: openapi3.Parameter{
				Schema: &openapi3.SchemaRef{
					Value: schema,
				},
			},
		},
	}
	return r
}

func (r *TestResponse) WithJSONSchema(schema *openapi3.Schema) *TestResponse {
	r.Response.WithJSONSchema(schema)
	return r
}

type TestResponses struct {
	openapi3.Responses
}

func createTestResponses() *TestResponses {
	return &TestResponses{
		Responses: *openapi3.NewResponses(),
	}
}

func (r *TestResponses) WithResponse(code string, response *openapi3.Response) *TestResponses {
	r.Responses.Set(code, &openapi3.ResponseRef{
		Value: response,
	})

	return r
}

func assertEqual(t *testing.T, got any, want any) {
	t.Helper()

	gotBytes, err := json.Marshal(got)
	if err != nil {
		t.Errorf("failed to marshal got: %v", err)
	}

	wantBytes, err := json.Marshal(want)
	if err != nil {
		t.Errorf("failed to marshal want: %v", err)
	}

	if !bytes.Equal(gotBytes, wantBytes) {
		t.Errorf("%v()\ngot = %v\nwant %v", t.Name(), string(gotBytes), string(wantBytes))
	}
}
