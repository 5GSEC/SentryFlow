package apispec

import (
	"reflect"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func Test_shouldIgnoreHeader(t *testing.T) {
	ignoredHeaders := map[string]struct{}{
		contentTypeHeaderName:       {},
		acceptTypeHeaderName:        {},
		authorizationTypeHeaderName: {},
	}
	type args struct {
		headerKey string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should ignore",
			args: args{
				headerKey: "Accept",
			},
			want: true,
		},
		{
			name: "should ignore",
			args: args{
				headerKey: "Content-Type",
			},
			want: true,
		},
		{
			name: "should ignore",
			args: args{
				headerKey: "Authorization",
			},
			want: true,
		},
		{
			name: "should not ignore",
			args: args{
				headerKey: "X-Test",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldIgnoreHeader(ignoredHeaders, tt.args.headerKey); got != tt.want {
				t.Errorf("shouldIgnoreHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addResponseHeader(t *testing.T) {
	op := NewOperationGenerator(OperationGeneratorConfig{})
	type args struct {
		response    *openapi3.Response
		headerKey   string
		headerValue string
	}
	tests := []struct {
		name string
		args args
		want *openapi3.Response
	}{
		{
			name: "primitive",
			args: args{
				response:    openapi3.NewResponse(),
				headerKey:   "X-Test-Uuid",
				headerValue: "77e1c83b-7bb0-437b-bc50-a7a58e5660ac",
			},
			want: createTestResponse().
				WithHeader("X-Test-Uuid", openapi3.NewUUIDSchema()).Response,
		},
		{
			name: "array",
			args: args{
				response:    openapi3.NewResponse(),
				headerKey:   "X-Test-Array",
				headerValue: "1,2,3,4",
			},
			want: createTestResponse().
				WithHeader("X-Test-Array", openapi3.NewArraySchema().WithItems(openapi3.NewInt64Schema())).Response,
		},
		{
			name: "date",
			args: args{
				response:    openapi3.NewResponse(),
				headerKey:   "date",
				headerValue: "Mon, 23 Aug 2021 06:52:48 GMT",
			},
			want: createTestResponse().
				WithHeader("date", openapi3.NewStringSchema()).Response,
		},
		{
			name: "ignore header",
			args: args{
				response:    openapi3.NewResponse(),
				headerKey:   "Accept",
				headerValue: "",
			},
			want: openapi3.NewResponse(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := op.addResponseHeader(tt.args.response, tt.args.headerKey, tt.args.headerValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addResponseHeader() = %v, want %v", marshal(got), marshal(tt.want))
			}
		})
	}
}

func Test_addHeaderParam(t *testing.T) {
	op := NewOperationGenerator(OperationGeneratorConfig{})
	type args struct {
		operation   *openapi3.Operation
		headerKey   string
		headerValue string
	}
	tests := []struct {
		name string
		args args
		want *openapi3.Operation
	}{
		{
			name: "primitive",
			args: args{
				operation:   openapi3.NewOperation(),
				headerKey:   "X-Test-Uuid",
				headerValue: "77e1c83b-7bb0-437b-bc50-a7a58e5660ac",
			},
			want: createTestOperation().WithParameter(openapi3.NewHeaderParameter("X-Test-Uuid").
				WithSchema(openapi3.NewUUIDSchema())).Op,
		},
		{
			name: "array",
			args: args{
				operation:   openapi3.NewOperation(),
				headerKey:   "X-Test-Array",
				headerValue: "1,2,3,4",
			},
			want: createTestOperation().WithParameter(openapi3.NewHeaderParameter("X-Test-Array").
				WithSchema(openapi3.NewArraySchema().WithItems(openapi3.NewInt64Schema()))).Op,
		},
		{
			name: "ignore header",
			args: args{
				operation:   openapi3.NewOperation(),
				headerKey:   "Accept",
				headerValue: "",
			},
			want: openapi3.NewOperation(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := op.addHeaderParam(tt.args.operation, tt.args.headerKey, tt.args.headerValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addHeaderParam() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createHeadersToIgnore(t *testing.T) {
	type args struct {
		headers []string
	}
	tests := []struct {
		name string
		args args
		want map[string]struct{}
	}{
		{
			name: "only default headers",
			args: args{
				headers: nil,
			},
			want: map[string]struct{}{
				acceptTypeHeaderName:        {},
				contentTypeHeaderName:       {},
				authorizationTypeHeaderName: {},
			},
		},
		{
			name: "with custom headers",
			args: args{
				headers: []string{
					"X-H1",
					"X-H2",
				},
			},
			want: map[string]struct{}{
				acceptTypeHeaderName:        {},
				contentTypeHeaderName:       {},
				authorizationTypeHeaderName: {},
				"x-h1":                      {},
				"x-h2":                      {},
			},
		},
		{
			name: "custom headers are sub list of the default headers",
			args: args{
				headers: []string{
					acceptTypeHeaderName,
					contentTypeHeaderName,
				},
			},
			want: map[string]struct{}{
				acceptTypeHeaderName:        {},
				contentTypeHeaderName:       {},
				authorizationTypeHeaderName: {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createHeadersToIgnore(tt.args.headers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createHeadersToIgnore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperationGenerator_addCookieParam(t *testing.T) {
	op := NewOperationGenerator(OperationGeneratorConfig{})
	type args struct {
		operation   *openapi3.Operation
		headerValue string
	}
	tests := []struct {
		name string
		args args
		want *openapi3.Operation
	}{
		{
			name: "sanity",
			args: args{
				operation:   openapi3.NewOperation(),
				headerValue: "debug=0; csrftoken=BUSe35dohU3O1MZvDCUOJ",
			},
			want: createTestOperation().
				WithParameter(openapi3.NewCookieParameter("debug").WithSchema(openapi3.NewInt64Schema())).
				WithParameter(openapi3.NewCookieParameter("csrftoken").WithSchema(openapi3.NewStringSchema())).
				Op,
		},
		{
			name: "array",
			args: args{
				operation:   openapi3.NewOperation(),
				headerValue: "array=1,2,3",
			},
			want: createTestOperation().
				WithParameter(openapi3.NewCookieParameter("array").WithSchema(openapi3.NewArraySchema().WithItems(openapi3.NewInt64Schema()))).
				Op,
		},
		{
			name: "unsupported cookie param",
			args: args{
				operation:   openapi3.NewOperation(),
				headerValue: "unsupported=unsupported=unsupported; csrftoken=BUSe35dohU3O1MZvDCUOJ",
			},
			want: createTestOperation().
				WithParameter(openapi3.NewCookieParameter("csrftoken").WithSchema(openapi3.NewStringSchema())).
				Op,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := op.addCookieParam(tt.args.operation, tt.args.headerValue); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addCookieParam() = %v, want %v", got, tt.want)
			}
		})
	}
}
