package apispec

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func newBoolSchemaWithAllowEmptyValue() *openapi3.Schema {
	schema := openapi3.NewBoolSchema()
	schema.AllowEmptyValue = true
	return schema
}

func Test_handleApplicationFormURLEncodedBody(t *testing.T) {
	type args struct {
		operation       *openapi3.Operation
		securitySchemes openapi3.SecuritySchemes
		body            string
	}
	tests := []struct {
		name    string
		args    args
		want    *openapi3.Operation
		want1   openapi3.SecuritySchemes
		wantErr bool
	}{
		{
			name: "sanity",
			args: args{
				operation: openapi3.NewOperation(),
				body:      "name=Amy&fav_number=321.1",
			},
			want: createTestOperation().WithRequestBody(openapi3.NewRequestBody().WithSchema(
				openapi3.NewObjectSchema().WithProperties(map[string]*openapi3.Schema{
					"name":       openapi3.NewStringSchema(),
					"fav_number": openapi3.NewFloat64Schema(),
				}), []string{mediaTypeApplicationForm})).Op,
		},
		{
			name: "parameters without a value",
			args: args{
				operation: openapi3.NewOperation(),
				body:      "foo&bar&baz",
			},
			want: createTestOperation().WithRequestBody(openapi3.NewRequestBody().WithSchema(
				openapi3.NewObjectSchema().WithProperties(map[string]*openapi3.Schema{
					"foo": newBoolSchemaWithAllowEmptyValue(),
					"bar": newBoolSchemaWithAllowEmptyValue(),
					"baz": newBoolSchemaWithAllowEmptyValue(),
				}), []string{mediaTypeApplicationForm})).Op,
		},
		{
			name: "multiple parameter instances",
			args: args{
				operation: openapi3.NewOperation(),
				body:      "param=value1&param=value2&param=value3",
			},
			want: createTestOperation().WithRequestBody(openapi3.NewRequestBody().WithSchema(
				openapi3.NewObjectSchema().WithProperties(map[string]*openapi3.Schema{
					"param": openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema()),
				}), []string{mediaTypeApplicationForm})).Op,
		},
		{
			name: "bad query",
			args: args{
				operation: openapi3.NewOperation(),
				body:      "name%2",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "OAuth2 security",
			args: args{
				operation:       openapi3.NewOperation(),
				body:            AccessTokenParamKey + "=token",
				securitySchemes: openapi3.SecuritySchemes{},
			},
			want: createTestOperation().
				WithSecurityRequirement(map[string][]string{OAuth2SecuritySchemeKey: {}}).Op,
			want1: openapi3.SecuritySchemes{
				OAuth2SecuritySchemeKey: {Value: NewOAuth2SecurityScheme([]string{})},
			},
		},
		{
			name: "OAuth2 security + some params",
			args: args{
				operation:       openapi3.NewOperation(),
				body:            AccessTokenParamKey + "=token&name=Amy",
				securitySchemes: openapi3.SecuritySchemes{},
			},
			want: createTestOperation().
				WithSecurityRequirement(map[string][]string{OAuth2SecuritySchemeKey: {}}).
				WithRequestBody(openapi3.NewRequestBody().WithSchema(
					openapi3.NewObjectSchema().WithProperties(map[string]*openapi3.Schema{
						"name": openapi3.NewStringSchema(),
					}), []string{mediaTypeApplicationForm})).Op,
			want1: openapi3.SecuritySchemes{
				OAuth2SecuritySchemeKey: {Value: NewOAuth2SecurityScheme([]string{})},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op, securitySchemes, err := handleApplicationFormURLEncodedBody(tt.args.operation, tt.args.securitySchemes, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("handleApplicationFormURLEncodedBody() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			op = sortParameters(op)
			tt.want = sortParameters(tt.want)

			assertEqual(t, op, tt.want)
			assertEqual(t, securitySchemes, tt.want1)
		})
	}
}

var formDataBodyMultipleFileUpload = "--cdce6441022a3dcf\r\n" +
	"Content-Disposition: form-data; name=\"fileName\"; filename=\"file1.txt\"\r\n\r\n" +
	"Content-Type: text/plain\r\n\r\n" +
	"File contents go here.\r\n" +
	"--cdce6441022a3dcf\r\n" +
	"Content-Disposition: form-data; name=\"fileName\"; filename=\"file2.png\"\r\n\r\n" +
	"Content-Type: image/png\r\n\r\n" +
	"File contents go here.\r\n" +
	"--cdce6441022a3dcf\r\n" +
	"Content-Disposition: form-data; name=\"fileName\"; filename=\"file3.jpg\"\r\n\r\n" +
	"Content-Type: image/jpeg\r\n\r\n" +
	"File contents go here.\r\n" +
	"--cdce6441022a3dcf--\r\n"

var formDataBody = "--cdce6441022a3dcf\r\n" +
	"Content-Disposition: form-data; name=\"upfile\"; filename=\"example.txt\"\r\n" +
	"Content-Type: text/plain\r\n\r\n" +
	"File contents go here.\r\n" +
	"--cdce6441022a3dcf\r\n" +
	"Content-Disposition: form-data; name=\"array-to-ignore-expected-string\"\r\n\r\n" +
	"1,2\r\n" +
	"--cdce6441022a3dcf\r\n" +
	"Content-Disposition: form-data; name=\"string\"\r\n\r\n" +
	"str\r\n" +
	"--cdce6441022a3dcf\r\n" +
	"Content-Disposition: form-data; name=\"integer\"\r\n\r\n" +
	"12\r\n" +
	"--cdce6441022a3dcf\r\n" +
	"Content-Disposition: form-data; name=\"id\"\r\n" +
	"Content-Type: text/plain\r\n\r\n" +
	"123e4567-e89b-12d3-a456-426655440000\r\n" +
	"--cdce6441022a3dcf\r\n" +
	"Content-Disposition: form-data; name=\"boolean\"\r\n\r\n" +
	"false\r\n" +
	"--cdce6441022a3dcf--\r\n"

func Test_addMultipartFormDataParams(t *testing.T) {
	type args struct {
		body   string
		params map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    *openapi3.Schema
		wantErr bool
	}{
		{
			name: "sanity",
			args: args{
				body:   formDataBody,
				params: map[string]string{"boundary": "cdce6441022a3dcf"},
			},
			want: openapi3.NewObjectSchema().WithProperties(map[string]*openapi3.Schema{
				"upfile":                          openapi3.NewStringSchema().WithFormat("binary"),
				"integer":                         openapi3.NewInt64Schema(),
				"boolean":                         openapi3.NewBoolSchema(),
				"string":                          openapi3.NewStringSchema(),
				"array-to-ignore-expected-string": openapi3.NewStringSchema(),
				"id":                              openapi3.NewUUIDSchema(),
			}),
			wantErr: false,
		},
		{
			name: "Multiple File Upload",
			args: args{
				body:   formDataBodyMultipleFileUpload,
				params: map[string]string{"boundary": "cdce6441022a3dcf"},
			},
			want: openapi3.NewObjectSchema().WithProperties(map[string]*openapi3.Schema{
				"fileName": openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema().WithFormat("binary")),
			}),
			wantErr: false,
		},
		{
			name: "missing boundary param",
			args: args{
				body:   formDataBody,
				params: map[string]string{},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getMultipartFormDataSchema(tt.args.body, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("getMultipartFormDataSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assertEqual(t, got, tt.want)
		})
	}
}
