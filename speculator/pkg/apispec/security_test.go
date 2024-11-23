package apispec

import (
	"reflect"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func createOperationWithSecurity(sec *openapi3.SecurityRequirements) *openapi3.Operation {
	operation := openapi3.NewOperation()
	operation.Security = sec
	return operation
}

func Test_updateSecurityDefinitionsFromOperation(t *testing.T) {
	type args struct {
		securitySchemes openapi3.SecuritySchemes
		op              *openapi3.Operation
	}
	tests := []struct {
		name string
		args args
		want openapi3.SecuritySchemes
	}{
		{
			name: "OAuth2 OR BasicAuth",
			args: args{
				securitySchemes: openapi3.SecuritySchemes{},
				op: createOperationWithSecurity(&openapi3.SecurityRequirements{
					{
						OAuth2SecuritySchemeKey: {"admin"},
					},
					{
						BasicAuthSecuritySchemeKey: {},
					},
				}),
			},
			want: openapi3.SecuritySchemes{
				OAuth2SecuritySchemeKey:    &openapi3.SecuritySchemeRef{Value: NewOAuth2SecurityScheme(nil)},
				BasicAuthSecuritySchemeKey: &openapi3.SecuritySchemeRef{Value: NewBasicAuthSecurityScheme()},
			},
		},
		{
			name: "OAuth2 AND BasicAuth",
			args: args{
				securitySchemes: openapi3.SecuritySchemes{},
				op: createOperationWithSecurity(&openapi3.SecurityRequirements{
					{
						OAuth2SecuritySchemeKey:    {"admin"},
						BasicAuthSecuritySchemeKey: {},
					},
				}),
			},
			want: openapi3.SecuritySchemes{
				OAuth2SecuritySchemeKey:    &openapi3.SecuritySchemeRef{Value: NewOAuth2SecurityScheme(nil)},
				BasicAuthSecuritySchemeKey: &openapi3.SecuritySchemeRef{Value: NewBasicAuthSecurityScheme()},
			},
		},
		{
			name: "OAuth2 AND BasicAuth OR BasicAuth",
			args: args{
				securitySchemes: openapi3.SecuritySchemes{},
				op: createOperationWithSecurity(&openapi3.SecurityRequirements{
					{
						OAuth2SecuritySchemeKey:    {"admin"},
						BasicAuthSecuritySchemeKey: {},
					},
					{
						BasicAuthSecuritySchemeKey: {},
					},
				}),
			},
			want: openapi3.SecuritySchemes{
				OAuth2SecuritySchemeKey:    &openapi3.SecuritySchemeRef{Value: NewOAuth2SecurityScheme(nil)},
				BasicAuthSecuritySchemeKey: &openapi3.SecuritySchemeRef{Value: NewBasicAuthSecurityScheme()},
			},
		},
		{
			name: "Unsupported SecurityDefinition key - no change to securitySchemes",
			args: args{
				securitySchemes: openapi3.SecuritySchemes{
					OAuth2SecuritySchemeKey: &openapi3.SecuritySchemeRef{Value: NewOAuth2SecurityScheme(nil)},
				},
				op: createOperationWithSecurity(&openapi3.SecurityRequirements{
					{
						"unsupported": {"admin"},
					},
				}),
			},
			want: openapi3.SecuritySchemes{
				OAuth2SecuritySchemeKey: &openapi3.SecuritySchemeRef{Value: NewOAuth2SecurityScheme(nil)},
			},
		},
		{
			name: "nil operation - no change to securitySchemes",
			args: args{
				securitySchemes: openapi3.SecuritySchemes{
					OAuth2SecuritySchemeKey: &openapi3.SecuritySchemeRef{Value: NewOAuth2SecurityScheme(nil)},
				},
				op: nil,
			},
			want: openapi3.SecuritySchemes{
				OAuth2SecuritySchemeKey: &openapi3.SecuritySchemeRef{Value: NewOAuth2SecurityScheme(nil)},
			},
		},
		{
			name: "operation without security - no change to securitySchemes",
			args: args{
				securitySchemes: openapi3.SecuritySchemes{
					OAuth2SecuritySchemeKey: &openapi3.SecuritySchemeRef{Value: NewOAuth2SecurityScheme(nil)},
				},
				op: createOperationWithSecurity(nil),
			},
			want: openapi3.SecuritySchemes{
				OAuth2SecuritySchemeKey: &openapi3.SecuritySchemeRef{Value: NewOAuth2SecurityScheme(nil)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := updateSecuritySchemesFromOperation(tt.args.securitySchemes, tt.args.op); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("updateSecuritySchemesFromOperation() = %v, want %v", got, tt.want)
			}
		})
	}
}
