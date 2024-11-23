package apispec

import (
	"encoding/json"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yudai/gojsondiff"
)

var agentStatusBody = `{"active":true,
"certificateVersion":"86eb5278-676a-3b7c-b29d-4a57007dc7be",
"controllerInstanceInfo":{"replicaId":"portshift-agent-66fc77c848-tmmk8"},
"policyAndAppVersion":1621477900361,
"statusCodes":["NO_METRICS_SERVER"],
"version":"1.147.1"}`

var cvssBody = `{"cvss":[{"score":7.8,"vector":"AV:L/AC:L/PR:N/UI:R/S:U/C:H/I:H/A:H","version":"3"}]}`

func generateDefaultOAuthToken(scopes []string) (string, string) {
	mySigningKey := []byte("AllYourBase")

	var defaultOAuth2Claims jwt.Claims = OAuth2Claims{
		strings.Join(scopes, " "),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "test",
			Subject:   "somebody",
			Audience:  []string{"somebody_else"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, defaultOAuth2Claims)
	bearerToken, err := token.SignedString(mySigningKey)
	if err != nil {
		logger.Errorf("Failed to create default OAuth2 Bearer Token: %v", err)
		return bearerToken, ""
	}

	oAuth2JSON := ""
	encoded, err := json.Marshal(scopes)
	if err != nil {
		logger.Errorf("Cannot encode token scopes: %v", scopes)
	} else {
		oAuth2JSON = string(encoded)
	}

	return bearerToken, oAuth2JSON
}

func generateQueryParams(t *testing.T, query string) url.Values {
	t.Helper()
	parseQuery, err := url.ParseQuery(query)
	if err != nil {
		t.Fatal(err)
	}
	return parseQuery
}

func TestGenerateSpecOperation(t *testing.T) {
	sd := openapi3.SecuritySchemes{}
	opGen := CreateTestNewOperationGenerator()
	operation, err := opGen.GenerateSpecOperation(&HTTPInteractionData{
		ReqBody:  agentStatusBody,
		RespBody: cvssBody,
		ReqHeaders: map[string]string{
			"X-Request-ID":        "77e1c83b-7bb0-437b-bc50-a7a58e5660ac",
			"X-Float-Test":        "12.2",
			"X-Collection-Test":   "a,b,c,d",
			contentTypeHeaderName: mediaTypeApplicationJSON,
		},
		RespHeaders: map[string]string{
			"X-RateLimit-Limit":   "12",
			"X-RateLimit-Reset":   "2016-10-12T11:00:00Z",
			contentTypeHeaderName: mediaTypeApplicationJSON,
		},
		QueryParams: generateQueryParams(t, "offset=30&limit=10"),
		statusCode:  200,
	}, sd)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(marshal(operation))
	t.Log(marshal(sd))
}

func validateOperation(t *testing.T, got *openapi3.Operation, want string) bool {
	t.Helper()
	templateB, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}

	differ := gojsondiff.New()
	diff, err := differ.Compare(templateB, []byte(want))
	if err != nil {
		t.Fatal(err)
	}
	return diff.Modified() == false
}

func TestGenerateSpecOperation1(t *testing.T) {
	defaultOAuth2Scopes := []string{"admin", "write:pets"}
	defaultOAuth2BearerToken, defaultOAuth2JSON := generateDefaultOAuthToken(defaultOAuth2Scopes)
	defaultOAuthSecurityScheme := NewOAuth2SecurityScheme(defaultOAuth2Scopes)
	defaultAPIKeyHeaderName := ""
	for key := range APIKeyNames {
		defaultAPIKeyHeaderName = key
		break
	}
	type args struct {
		data *HTTPInteractionData
	}
	opGen := CreateTestNewOperationGenerator()
	tests := []struct {
		name       string
		args       args
		want       string
		wantErr    bool
		expectedSd openapi3.SecuritySchemes
	}{
		{
			name: "Basic authorization req header",
			args: args{
				data: &HTTPInteractionData{
					ReqBody:  agentStatusBody,
					RespBody: cvssBody,
					ReqHeaders: map[string]string{
						contentTypeHeaderName:       mediaTypeApplicationHalJSON,
						authorizationTypeHeaderName: BasicAuthPrefix + "=token",
					},
					RespHeaders: map[string]string{
						contentTypeHeaderName: mediaTypeApplicationHalJSON,
					},
					statusCode: 200,
				},
			},
			want: "{\"requestBody\":{\"content\":{\"application/json\":{\"schema\":{\"properties\":{\"active\":{\"type\":\"boolean\"},\"certificateVersion\":{\"format\":\"uuid\",\"type\":\"string\"},\"controllerInstanceInfo\":{\"properties\":{\"replicaId\":{\"type\":\"string\"}},\"type\":\"object\"},\"policyAndAppVersion\":{\"format\":\"int64\",\"type\":\"integer\"},\"statusCodes\":{\"items\":{\"type\":\"string\"},\"type\":\"array\"},\"version\":{\"type\":\"string\"}},\"type\":\"object\"}}}},\"responses\":{\"200\":{\"content\":{\"application/json\":{\"schema\":{\"properties\":{\"cvss\":{\"items\":{\"properties\":{\"score\":{\"type\":\"number\"},\"vector\":{\"type\":\"string\"},\"version\":{\"type\":\"string\"}},\"type\":\"object\"},\"type\":\"array\"}},\"type\":\"object\"}}},\"description\":\"response\"},\"default\":{\"description\":\"default\"}},\"security\":[{\"BasicAuth\":[]}]}",
			expectedSd: openapi3.SecuritySchemes{
				BasicAuthSecuritySchemeKey: &openapi3.SecuritySchemeRef{Value: NewBasicAuthSecurityScheme()},
			},
			wantErr: false,
		},
		{
			name: "OAuth 2.0 authorization req header",
			args: args{
				data: &HTTPInteractionData{
					ReqBody:  agentStatusBody,
					RespBody: cvssBody,
					ReqHeaders: map[string]string{
						contentTypeHeaderName:       mediaTypeApplicationJSON,
						authorizationTypeHeaderName: BearerAuthPrefix + defaultOAuth2BearerToken,
					},
					RespHeaders: map[string]string{
						contentTypeHeaderName: mediaTypeApplicationJSON,
					},
					statusCode: 200,
				},
			},
			want: "{\"requestBody\":{\"content\":{\"application/json\":{\"schema\":{\"properties\":{\"active\":{\"type\":\"boolean\"},\"certificateVersion\":{\"format\":\"uuid\",\"type\":\"string\"},\"controllerInstanceInfo\":{\"properties\":{\"replicaId\":{\"type\":\"string\"}},\"type\":\"object\"},\"policyAndAppVersion\":{\"format\":\"int64\",\"type\":\"integer\"},\"statusCodes\":{\"items\":{\"type\":\"string\"},\"type\":\"array\"},\"version\":{\"type\":\"string\"}},\"type\":\"object\"}}}},\"responses\":{\"200\":{\"content\":{\"application/json\":{\"schema\":{\"properties\":{\"cvss\":{\"items\":{\"properties\":{\"score\":{\"type\":\"number\"},\"vector\":{\"type\":\"string\"},\"version\":{\"type\":\"string\"}},\"type\":\"object\"},\"type\":\"array\"}},\"type\":\"object\"}}},\"description\":\"response\"},\"default\":{\"description\":\"default\"}},\"security\":[{\"OAuth2\":[\"admin\",\"write:pets\"]}]}",
			expectedSd: openapi3.SecuritySchemes{
				OAuth2SecuritySchemeKey: &openapi3.SecuritySchemeRef{Value: defaultOAuthSecurityScheme},
			},
			wantErr: false,
		},
		{
			name: "OAuth 2.0 URI Query Parameter",
			args: args{
				data: &HTTPInteractionData{
					ReqBody:  agentStatusBody,
					RespBody: cvssBody,
					ReqHeaders: map[string]string{
						contentTypeHeaderName: mediaTypeApplicationJSON,
					},
					RespHeaders: map[string]string{
						contentTypeHeaderName: mediaTypeApplicationJSON,
					},
					QueryParams: generateQueryParams(t, AccessTokenParamKey+"="+defaultOAuth2BearerToken),
					statusCode:  200,
				},
			},
			want: "{\"requestBody\":{\"content\":{\"application/json\":{\"schema\":{\"properties\":{\"active\":{\"type\":\"boolean\"},\"certificateVersion\":{\"format\":\"uuid\",\"type\":\"string\"},\"controllerInstanceInfo\":{\"properties\":{\"replicaId\":{\"type\":\"string\"}},\"type\":\"object\"},\"policyAndAppVersion\":{\"format\":\"int64\",\"type\":\"integer\"},\"statusCodes\":{\"items\":{\"type\":\"string\"},\"type\":\"array\"},\"version\":{\"type\":\"string\"}},\"type\":\"object\"}}}},\"responses\":{\"200\":{\"content\":{\"application/json\":{\"schema\":{\"properties\":{\"cvss\":{\"items\":{\"properties\":{\"score\":{\"type\":\"number\"},\"vector\":{\"type\":\"string\"},\"version\":{\"type\":\"string\"}},\"type\":\"object\"},\"type\":\"array\"}},\"type\":\"object\"}}},\"description\":\"response\"},\"default\":{\"description\":\"default\"}},\"security\":[{\"OAuth2\":[\"admin\",\"write:pets\"]}]}",
			expectedSd: openapi3.SecuritySchemes{
				OAuth2SecuritySchemeKey: &openapi3.SecuritySchemeRef{Value: defaultOAuthSecurityScheme},
			},
			wantErr: false,
		},
		{
			name: "OAuth 2.0 Form-Encoded Body Parameter",
			args: args{
				data: &HTTPInteractionData{
					ReqBody:  AccessTokenParamKey + "=" + defaultOAuth2BearerToken + "&key=val",
					RespBody: cvssBody,
					ReqHeaders: map[string]string{
						contentTypeHeaderName: mediaTypeApplicationForm,
					},
					RespHeaders: map[string]string{
						contentTypeHeaderName: mediaTypeApplicationJSON,
					},
					statusCode: 200,
				},
			},
			want: "{\"requestBody\":{\"content\":{\"application/x-www-form-urlencoded\":{\"schema\":{\"properties\":{\"key\":{\"type\":\"string\"}},\"type\":\"object\"}}}},\"responses\":{\"200\":{\"content\":{\"application/json\":{\"schema\":{\"properties\":{\"cvss\":{\"items\":{\"properties\":{\"score\":{\"type\":\"number\"},\"vector\":{\"type\":\"string\"},\"version\":{\"type\":\"string\"}},\"type\":\"object\"},\"type\":\"array\"}},\"type\":\"object\"}}},\"description\":\"response\"},\"default\":{\"description\":\"default\"}},\"security\":[{\"OAuth2\":[\"admin\",\"write:pets\"]}]}",
			expectedSd: openapi3.SecuritySchemes{
				OAuth2SecuritySchemeKey: &openapi3.SecuritySchemeRef{Value: defaultOAuthSecurityScheme},
			},
			wantErr: false,
		},
		{
			name: "OAuth 2.0 Multiple parameters: Authorization Req Header and URI Query Parameter",
			args: args{
				data: &HTTPInteractionData{
					ReqBody:  agentStatusBody,
					RespBody: cvssBody,
					ReqHeaders: map[string]string{
						contentTypeHeaderName:       mediaTypeApplicationJSON,
						authorizationTypeHeaderName: BearerAuthPrefix + defaultOAuth2BearerToken,
					},
					RespHeaders: map[string]string{
						contentTypeHeaderName: mediaTypeApplicationJSON,
					},
					QueryParams: generateQueryParams(t, AccessTokenParamKey+"=bogus.key.material"),
					statusCode:  200,
				},
			},
			want: "{\"requestBody\":{\"content\":{\"application/json\":{\"schema\":{\"properties\":{\"active\":{\"type\":\"boolean\"},\"certificateVersion\":{\"format\":\"uuid\",\"type\":\"string\"},\"controllerInstanceInfo\":{\"properties\":{\"replicaId\":{\"type\":\"string\"}},\"type\":\"object\"},\"policyAndAppVersion\":{\"format\":\"int64\",\"type\":\"integer\"},\"statusCodes\":{\"items\":{\"type\":\"string\"},\"type\":\"array\"},\"version\":{\"type\":\"string\"}},\"type\":\"object\"}}}},\"responses\":{\"200\":{\"content\":{\"application/json\":{\"schema\":{\"properties\":{\"cvss\":{\"items\":{\"properties\":{\"score\":{\"type\":\"number\"},\"vector\":{\"type\":\"string\"},\"version\":{\"type\":\"string\"}},\"type\":\"object\"},\"type\":\"array\"}},\"type\":\"object\"}}},\"description\":\"response\"},\"default\":{\"description\":\"default\"}},\"security\":[{\"OAuth2\":" + defaultOAuth2JSON + "}]}",
			expectedSd: openapi3.SecuritySchemes{
				// Note: Auth Header will be used before Query Parameter is ignored.
				OAuth2SecuritySchemeKey: &openapi3.SecuritySchemeRef{Value: defaultOAuthSecurityScheme},
			},
			wantErr: false,
		},
		{
			name: "API Key in header",
			args: args{
				data: &HTTPInteractionData{
					ReqBody:  agentStatusBody,
					RespBody: cvssBody,
					ReqHeaders: map[string]string{
						contentTypeHeaderName:   mediaTypeApplicationJSON,
						defaultAPIKeyHeaderName: "mybogusapikey",
					},
					RespHeaders: map[string]string{
						contentTypeHeaderName: mediaTypeApplicationJSON,
					},
					statusCode: 200,
				},
			},
			want: "{\"requestBody\":{\"content\":{\"application/json\":{\"schema\":{\"properties\":{\"active\":{\"type\":\"boolean\"},\"certificateVersion\":{\"format\":\"uuid\",\"type\":\"string\"},\"controllerInstanceInfo\":{\"properties\":{\"replicaId\":{\"type\":\"string\"}},\"type\":\"object\"},\"policyAndAppVersion\":{\"format\":\"int64\",\"type\":\"integer\"},\"statusCodes\":{\"items\":{\"type\":\"string\"},\"type\":\"array\"},\"version\":{\"type\":\"string\"}},\"type\":\"object\"}}}},\"responses\":{\"200\":{\"content\":{\"application/json\":{\"schema\":{\"properties\":{\"cvss\":{\"items\":{\"properties\":{\"score\":{\"type\":\"number\"},\"vector\":{\"type\":\"string\"},\"version\":{\"type\":\"string\"}},\"type\":\"object\"},\"type\":\"array\"}},\"type\":\"object\"}}},\"description\":\"response\"},\"default\":{\"description\":\"default\"}},\"security\":[{\"ApiKeyAuth\":[]}]}",
			expectedSd: openapi3.SecuritySchemes{
				APIKeyAuthSecuritySchemeKey: &openapi3.SecuritySchemeRef{Value: NewAPIKeySecuritySchemeInHeader(defaultAPIKeyHeaderName)},
			},
			wantErr: false,
		},
		{
			name: "API Key URI Query Parameter",
			args: args{
				data: &HTTPInteractionData{
					ReqBody:  agentStatusBody,
					RespBody: cvssBody,
					ReqHeaders: map[string]string{
						contentTypeHeaderName: mediaTypeApplicationJSON,
					},
					RespHeaders: map[string]string{
						contentTypeHeaderName: mediaTypeApplicationJSON,
					},
					QueryParams: generateQueryParams(t, defaultAPIKeyHeaderName+"=mybogusapikey"),
					statusCode:  200,
				},
			},
			want: "{\"requestBody\":{\"content\":{\"application/json\":{\"schema\":{\"properties\":{\"active\":{\"type\":\"boolean\"},\"certificateVersion\":{\"format\":\"uuid\",\"type\":\"string\"},\"controllerInstanceInfo\":{\"properties\":{\"replicaId\":{\"type\":\"string\"}},\"type\":\"object\"},\"policyAndAppVersion\":{\"format\":\"int64\",\"type\":\"integer\"},\"statusCodes\":{\"items\":{\"type\":\"string\"},\"type\":\"array\"},\"version\":{\"type\":\"string\"}},\"type\":\"object\"}}}},\"responses\":{\"200\":{\"content\":{\"application/json\":{\"schema\":{\"properties\":{\"cvss\":{\"items\":{\"properties\":{\"score\":{\"type\":\"number\"},\"vector\":{\"type\":\"string\"},\"version\":{\"type\":\"string\"}},\"type\":\"object\"},\"type\":\"array\"}},\"type\":\"object\"}}},\"description\":\"response\"},\"default\":{\"description\":\"default\"}},\"security\":[{\"ApiKeyAuth\":[]}]}",
			expectedSd: openapi3.SecuritySchemes{
				APIKeyAuthSecuritySchemeKey: &openapi3.SecuritySchemeRef{Value: NewAPIKeySecuritySchemeInQuery(defaultAPIKeyHeaderName)},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sd := openapi3.SecuritySchemes{}
			got, err := opGen.GenerateSpecOperation(tt.args.data, sd)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateSpecOperation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !validateOperation(t, got, tt.want) {
				t.Errorf("GenerateSpecOperation() got = %v, want %v", marshal(got), marshal(tt.want))
			}

			assertEqual(t, sd, tt.expectedSd)
		})
	}
}

func Test_getStringSchema(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantSchema *openapi3.Schema
	}{
		{
			name: "date",
			args: args{
				value: "2017-07-21",
			},
			wantSchema: openapi3.NewStringSchema().WithFormat("date"),
		},
		{
			name: "time",
			args: args{
				value: "17:32:28",
			},
			wantSchema: openapi3.NewStringSchema().WithFormat("time"),
		},
		{
			name: "date-time",
			args: args{
				value: "2017-07-21T17:32:28Z",
			},
			wantSchema: openapi3.NewDateTimeSchema(),
		},
		{
			name: "email",
			args: args{
				value: "test@securecn.com",
			},
			wantSchema: openapi3.NewStringSchema().WithFormat("email"),
		},
		{
			name: "ipv4",
			args: args{
				value: "1.1.1.1",
			},
			wantSchema: openapi3.NewStringSchema().WithFormat("ipv4"),
		},
		{
			name: "ipv6",
			args: args{
				value: "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			},
			wantSchema: openapi3.NewStringSchema().WithFormat("ipv6"),
		},
		{
			name: "uuid",
			args: args{
				value: "123e4567-e89b-12d3-a456-426614174000",
			},
			wantSchema: openapi3.NewStringSchema().WithFormat("uuid"),
		},
		{
			name: "json-pointer",
			args: args{
				value: "/k%22l",
			},
			wantSchema: openapi3.NewStringSchema().WithFormat("json-pointer"),
		},
		{
			name: "string",
			args: args{
				value: "it is very hard to get a simple string",
			},
			wantSchema: openapi3.NewStringSchema(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSchema := getStringSchema(tt.args.value); !reflect.DeepEqual(gotSchema, tt.wantSchema) {
				t.Errorf("getStringSchema() = %v, want %v", gotSchema, tt.wantSchema)
			}
		})
	}
}

func Test_getNumberSchema(t *testing.T) {
	type args struct {
		value interface{}
	}
	tests := []struct {
		name       string
		args       args
		wantSchema *openapi3.Schema
	}{
		{
			name: "int",
			args: args{
				value: json.Number("85"),
			},
			wantSchema: openapi3.NewInt64Schema(),
		},
		{
			name: "float",
			args: args{
				value: json.Number("85.1"),
			},
			wantSchema: openapi3.NewFloat64Schema(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSchema := getNumberSchema(tt.args.value); !reflect.DeepEqual(gotSchema, tt.wantSchema) {
				t.Errorf("getNumberSchema() = %v, want %v", gotSchema, tt.wantSchema)
			}
		})
	}
}

func Test_escapeString(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "nothing to strip",
			args: args{
				key: "key",
			},
			want: "key",
		},
		{
			name: "escape double quotes",
			args: args{
				key: "{\"key1\":\"value1\", \"key2\":\"value2\"}",
			},
			want: "{\\\"key1\\\":\\\"value1\\\", \\\"key2\\\":\\\"value2\\\"}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := escapeString(tt.args.key); got != tt.want {
				t.Errorf("stripKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCloneOperation(t *testing.T) {
	type args struct {
		op *openapi3.Operation
	}
	tests := []struct {
		name    string
		args    args
		want    *openapi3.Operation
		wantErr bool
	}{
		{
			name: "sanity",
			args: args{
				op: createTestOperation().
					WithParameter(openapi3.NewHeaderParameter("header")).
					WithResponse(200, openapi3.NewResponse().WithDescription("keep").
						WithJSONSchemaRef(openapi3.NewSchemaRef("",
							openapi3.NewObjectSchema().WithProperty("test", openapi3.NewStringSchema())))).Op,
			},
			want: createTestOperation().
				WithParameter(openapi3.NewHeaderParameter("header")).
				WithResponse(200, openapi3.NewResponse().WithDescription("keep").
					WithJSONSchemaRef(openapi3.NewSchemaRef("",
						openapi3.NewObjectSchema().WithProperty("test", openapi3.NewStringSchema())))).Op,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CloneOperation(tt.args.op)
			if (err != nil) != tt.wantErr {
				t.Errorf("CloneOperation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assertEqual(t, got, tt.want)
			if got != nil {
				got.Responses = nil
				if tt.args.op.Responses == nil {
					t.Errorf("CloneOperation() original object should not have been changed")
					return
				}
			}
		})
	}
}

func Test_handleAuthReqHeader(t *testing.T) {
	type args struct {
		operation       *openapi3.Operation
		securitySchemes openapi3.SecuritySchemes
		value           string
	}
	defaultOAuth2Scopes := []string{"superman", "write:novel"}
	defaultOAuth2BearerToken, _ := generateDefaultOAuthToken(defaultOAuth2Scopes)
	defaultOAuthSecurityScheme := NewOAuth2SecurityScheme(defaultOAuth2Scopes)
	tests := []struct {
		name   string
		args   args
		wantOp *openapi3.Operation
		wantSd openapi3.SecuritySchemes
	}{
		{
			name: "BearerAuthPrefix",
			args: args{
				operation:       openapi3.NewOperation(),
				securitySchemes: openapi3.SecuritySchemes{},
				value:           BearerAuthPrefix + defaultOAuth2BearerToken,
			},
			wantOp: createTestOperation().WithSecurityRequirement(map[string][]string{OAuth2SecuritySchemeKey: defaultOAuth2Scopes}).Op,
			wantSd: openapi3.SecuritySchemes{
				OAuth2SecuritySchemeKey: &openapi3.SecuritySchemeRef{Value: defaultOAuthSecurityScheme},
			},
		},
		{
			name: "BasicAuthPrefix",
			args: args{
				operation:       openapi3.NewOperation(),
				securitySchemes: openapi3.SecuritySchemes{},
				value:           BasicAuthPrefix + "token",
			},
			wantOp: createTestOperation().WithSecurityRequirement(map[string][]string{BasicAuthSecuritySchemeKey: {}}).Op,
			wantSd: openapi3.SecuritySchemes{
				BasicAuthSecuritySchemeKey: &openapi3.SecuritySchemeRef{Value: NewBasicAuthSecurityScheme()},
			},
		},
		{
			name: "ignoring unknown authorization header value",
			args: args{
				operation:       openapi3.NewOperation(),
				securitySchemes: openapi3.SecuritySchemes{},
				value:           "invalid token",
			},
			wantOp: openapi3.NewOperation(),
			wantSd: openapi3.SecuritySchemes{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := handleAuthReqHeader(tt.args.operation, tt.args.securitySchemes, tt.args.value)
			if !reflect.DeepEqual(got, tt.wantOp) {
				t.Errorf("handleAuthReqHeader() got = %v, want %v", got, tt.wantOp)
			}
			if !reflect.DeepEqual(got1, tt.wantSd) {
				t.Errorf("handleAuthReqHeader() got1 = %v, want %v", got1, tt.wantSd)
			}
		})
	}
}

func Test_getScopesFromJWTClaims(t *testing.T) {
	type args struct {
		claims jwt.MapClaims
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "nil claims - expected nil scopes",
			args: args{
				claims: nil,
			},
			want: nil,
		},
		{
			name: "no scopes defined - expected nil scopes",
			args: args{
				claims: jwt.MapClaims{
					"no-scopes": "123",
				},
			},
			want: nil,
		},
		{
			name: "no scopes defined - expected nil scopes",
			args: args{
				claims: jwt.MapClaims{
					"no-scope": "123",
					"scope":    "scope1 scope2 scope3",
				},
			},
			want: []string{"scope1", "scope2", "scope3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getScopesFromJWTClaims(tt.args.claims); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getScopesFromJWTClaims() = %v, want %v", got, tt.want)
			}
		})
	}
}
