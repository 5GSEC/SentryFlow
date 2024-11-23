package apispec

//import (
//	"reflect"
//	"sort"
//	"testing"
//
//	"github.com/getkin/kin-openapi/openapi3"
//	"k8s.io/utils/field"
//)
//
//func Test_merge(t *testing.T) {
//	securitySchemes := openapi3.SecuritySchemes{}
//	op := CreateTestNewOperationGenerator()
//	op1, err := op.GenerateSpecOperation(&HTTPInteractionData{
//		ReqBody:     req1,
//		RespBody:    res1,
//		ReqHeaders:  map[string]string{"X-Test-Req-1": "1", contentTypeHeaderName: mediaTypeApplicationJSON},
//		RespHeaders: map[string]string{"X-Test-Res-1": "1", contentTypeHeaderName: mediaTypeApplicationJSON},
//		statusCode:  200,
//	}, securitySchemes)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	op2, err := op.GenerateSpecOperation(&HTTPInteractionData{
//		ReqBody:     req2,
//		RespBody:    res2,
//		ReqHeaders:  map[string]string{"X-Test-Req-2": "2", contentTypeHeaderName: mediaTypeApplicationJSON},
//		RespHeaders: map[string]string{"X-Test-Res-2": "2", contentTypeHeaderName: mediaTypeApplicationJSON},
//		statusCode:  200,
//	}, securitySchemes)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	combinedOp, err := op.GenerateSpecOperation(&HTTPInteractionData{
//		ReqBody:     combinedReq,
//		RespBody:    combinedRes,
//		ReqHeaders:  map[string]string{"X-Test-Req-1": "1", "X-Test-Req-2": "2", contentTypeHeaderName: mediaTypeApplicationJSON},
//		RespHeaders: map[string]string{"X-Test-Res-1": "1", "X-Test-Res-2": "2", contentTypeHeaderName: mediaTypeApplicationJSON},
//		statusCode:  200,
//	}, securitySchemes)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	type args struct {
//		operation1 *openapi3.Operation
//		operation2 *openapi3.Operation
//	}
//	tests := []struct {
//		name          string
//		args          args
//		want          *openapi3.Operation
//		wantConflicts bool
//	}{
//		{
//			name: "sanity",
//			args: args{
//				operation1: op1,
//				operation2: op2,
//			},
//			want: combinedOp,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, conflicts := mergeOperation(tt.args.operation1, tt.args.operation2)
//			if (len(conflicts) > 0) != tt.wantConflicts {
//				t.Errorf("merge() conflicts = %v, wantConflicts %v", conflicts, tt.wantConflicts)
//				return
//			}
//			got = sortParameters(got)
//			tt.want = sortParameters(tt.want)
//			assertEqual(t, got, tt.want)
//			//assert.DeepEqual(t, got, tt.want, cmpopts.IgnoreUnexported(openapi3.Schema{}))
//		})
//	}
//}
//
//func Test_shouldReturnIfEmpty(t *testing.T) {
//	type args struct {
//		a interface{}
//		b interface{}
//	}
//	tests := []struct {
//		name  string
//		args  args
//		want  interface{}
//		want1 bool
//	}{
//		{
//			name: "second nil",
//			args: args{
//				a: openapi3.NewOperation(),
//				b: nil,
//			},
//			want:  openapi3.NewOperation(),
//			want1: true,
//		},
//		{
//			name: "first nil",
//			args: args{
//				a: nil,
//				b: openapi3.NewOperation(),
//			},
//			want:  openapi3.NewOperation(),
//			want1: true,
//		},
//		{
//			name: "both nil",
//			args: args{
//				a: nil,
//				b: nil,
//			},
//			want:  nil,
//			want1: true,
//		},
//		{
//			name: "not nil",
//			args: args{
//				a: openapi3.NewOperation(),
//				b: openapi3.NewOperation(),
//			},
//			want:  nil,
//			want1: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, got1 := shouldReturnIfNil(tt.args.a, tt.args.b)
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("shouldReturnIfNil() got = %v, want %v", got, tt.want)
//			}
//			if got1 != tt.want1 {
//				t.Errorf("shouldReturnIfNil() got1 = %v, want %v", got1, tt.want1)
//			}
//		})
//	}
//}
//
//func Test_shouldReturnIfNil(t *testing.T) {
//	var nilSchema *openapi3.Schema
//	schema := openapi3.Schema{Type: &openapi3.Types{"test"}}
//	type args struct {
//		a interface{}
//		b interface{}
//	}
//	tests := []struct {
//		name  string
//		args  args
//		want  interface{}
//		want1 bool
//	}{
//		{
//			name: "a is nil b is not",
//			args: args{
//				a: nilSchema,
//				b: schema,
//			},
//			want:  schema,
//			want1: true,
//		},
//		{
//			name: "b is nil a is not",
//			args: args{
//				a: schema,
//				b: nilSchema,
//			},
//			want:  schema,
//			want1: true,
//		},
//		{
//			name: "both nil",
//			args: args{
//				a: nilSchema,
//				b: nilSchema,
//			},
//			want:  nilSchema,
//			want1: true,
//		},
//		{
//			name: "both not nil",
//			args: args{
//				a: schema,
//				b: schema,
//			},
//			want:  nil,
//			want1: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, got1 := shouldReturnIfNil(tt.args.a, tt.args.b)
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("shouldReturnIfNil() got = %v, want %v", got, tt.want)
//			}
//			if got1 != tt.want1 {
//				t.Errorf("shouldReturnIfNil() got1 = %v, want %v", got1, tt.want1)
//			}
//		})
//	}
//}
//
//func Test_shouldReturnIfEmptySchemaType(t *testing.T) {
//	emptySchemaType := &openapi3.Schema{}
//	schema := &openapi3.Schema{Type: &openapi3.Types{"test"}}
//	type args struct {
//		s  *openapi3.Schema
//		s2 *openapi3.Schema
//	}
//	tests := []struct {
//		name  string
//		args  args
//		want  *openapi3.Schema
//		want1 bool
//	}{
//		{
//			name: "first is empty second is not",
//			args: args{
//				s:  emptySchemaType,
//				s2: schema,
//			},
//			want:  schema,
//			want1: true,
//		},
//		{
//			name: "second is empty first is not",
//			args: args{
//				s:  schema,
//				s2: emptySchemaType,
//			},
//			want:  schema,
//			want1: true,
//		},
//		{
//			name: "both empty",
//			args: args{
//				s:  emptySchemaType,
//				s2: emptySchemaType,
//			},
//			want:  emptySchemaType,
//			want1: true,
//		},
//		{
//			name: "both not empty",
//			args: args{
//				s:  schema,
//				s2: schema,
//			},
//			want:  nil,
//			want1: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, got1 := shouldReturnIfEmptySchemaType(tt.args.s, tt.args.s2)
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("shouldReturnIfEmptySchemaType() got = %v, want %v", got, tt.want)
//			}
//			if got1 != tt.want1 {
//				t.Errorf("shouldReturnIfEmptySchemaType() got1 = %v, want %v", got1, tt.want1)
//			}
//		})
//	}
//}
//
//func Test_shouldReturnIfEmptyParameters(t *testing.T) {
//	var emptyParameters openapi3.Parameters
//	parameters := openapi3.Parameters{
//		&openapi3.ParameterRef{Value: openapi3.NewHeaderParameter("test")},
//	}
//	type args struct {
//		parameters  openapi3.Parameters
//		parameters2 openapi3.Parameters
//	}
//	tests := []struct {
//		name  string
//		args  args
//		want  openapi3.Parameters
//		want1 bool
//	}{
//		{
//			name: "first is empty second is not",
//			args: args{
//				parameters:  emptyParameters,
//				parameters2: parameters,
//			},
//			want:  parameters,
//			want1: true,
//		},
//		{
//			name: "second is empty first is not",
//			args: args{
//				parameters:  parameters,
//				parameters2: emptyParameters,
//			},
//			want:  parameters,
//			want1: true,
//		},
//		{
//			name: "both empty",
//			args: args{
//				parameters:  emptyParameters,
//				parameters2: emptyParameters,
//			},
//			want:  emptyParameters,
//			want1: true,
//		},
//		{
//			name: "both not empty",
//			args: args{
//				parameters:  parameters,
//				parameters2: parameters,
//			},
//			want:  nil,
//			want1: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, got1 := shouldReturnIfEmptyParameters(tt.args.parameters, tt.args.parameters2)
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("shouldReturnIfEmptyParameters() got = %v, want %v", got, tt.want)
//			}
//			if got1 != tt.want1 {
//				t.Errorf("shouldReturnIfEmptyParameters() got1 = %v, want %v", got1, tt.want1)
//			}
//		})
//	}
//}
//
//func Test_mergeHeader(t *testing.T) {
//	type args struct {
//		header  *openapi3.Header
//		header2 *openapi3.Header
//		child   *field.Path
//	}
//	tests := []struct {
//		name  string
//		args  args
//		want  *openapi3.Header
//		want1 []conflict
//	}{
//		{
//			name: "nothing to merge",
//			args: args{
//				header: &openapi3.Header{
//					Parameter: *openapi3.NewHeaderParameter("test").WithSchema(openapi3.NewStringSchema()),
//				},
//				header2: &openapi3.Header{
//					Parameter: *openapi3.NewHeaderParameter("test").WithSchema(openapi3.NewStringSchema()),
//				},
//				child: nil,
//			},
//			want: &openapi3.Header{
//				Parameter: *openapi3.NewHeaderParameter("test").WithSchema(openapi3.NewStringSchema()),
//			},
//			want1: nil,
//		},
//		{
//			name: "merge string type removal",
//			args: args{
//				header: &openapi3.Header{
//					Parameter: *openapi3.NewHeaderParameter("test").WithSchema(openapi3.NewStringSchema()),
//				},
//				header2: &openapi3.Header{
//					Parameter: *openapi3.NewHeaderParameter("test").WithSchema(openapi3.NewUUIDSchema()),
//				},
//				child: nil,
//			},
//			want: &openapi3.Header{
//				Parameter: *openapi3.NewHeaderParameter("test").WithSchema(openapi3.NewStringSchema()),
//			},
//			want1: nil,
//		},
//		{
//			name: "header in conflicts",
//			args: args{
//				header: &openapi3.Header{
//					Parameter: *openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewStringSchema()),
//				},
//				header2: &openapi3.Header{
//					Parameter: *openapi3.NewCookieParameter("cookie").WithSchema(openapi3.NewArraySchema()),
//				},
//				child: field.NewPath("test"),
//			},
//			want: &openapi3.Header{
//				Parameter: *openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewStringSchema()),
//			},
//			want1: []conflict{
//				{
//					path: field.NewPath("test"),
//					obj1: &openapi3.Header{
//						Parameter: *openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewStringSchema()),
//					},
//					obj2: &openapi3.Header{
//						Parameter: *openapi3.NewCookieParameter("cookie").WithSchema(openapi3.NewArraySchema()),
//					},
//					msg: createHeaderInConflictMsg(field.NewPath("test"), openapi3.ParameterInHeader, openapi3.ParameterInCookie),
//				},
//			},
//		},
//		{
//			name: "type conflicts prefer string",
//			args: args{
//				header: &openapi3.Header{
//					Parameter: *openapi3.NewHeaderParameter("test").WithSchema(openapi3.NewStringSchema()),
//				},
//				header2: &openapi3.Header{
//					Parameter: *openapi3.NewHeaderParameter("test").WithSchema(openapi3.NewArraySchema()),
//				},
//				child: field.NewPath("test"),
//			},
//			want: &openapi3.Header{
//				Parameter: *openapi3.NewHeaderParameter("test").WithSchema(openapi3.NewStringSchema()),
//			},
//			want1: nil,
//		},
//		{
//			name: "type conflicts",
//			args: args{
//				header: &openapi3.Header{
//					Parameter: *openapi3.NewHeaderParameter("test").WithSchema(openapi3.NewInt64Schema()),
//				},
//				header2: &openapi3.Header{
//					Parameter: *openapi3.NewHeaderParameter("test").WithSchema(openapi3.NewArraySchema()),
//				},
//				child: field.NewPath("test"),
//			},
//			want: &openapi3.Header{
//				Parameter: *openapi3.NewHeaderParameter("test").WithSchema(openapi3.NewInt64Schema()),
//			},
//			want1: []conflict{
//				{
//					path: field.NewPath("test"),
//					obj1: openapi3.NewInt64Schema(),
//					obj2: openapi3.NewArraySchema(),
//					msg:  createConflictMsg(field.NewPath("test"), openapi3.TypeInteger, openapi3.TypeArray),
//				},
//			},
//		},
//		{
//			name: "empty header",
//			args: args{
//				header: &openapi3.Header{
//					Parameter: *openapi3.NewHeaderParameter("empty"),
//				},
//				header2: &openapi3.Header{
//					Parameter: *openapi3.NewHeaderParameter("test").WithSchema(openapi3.NewArraySchema()),
//				},
//				child: field.NewPath("test"),
//			},
//			want: &openapi3.Header{
//				Parameter: *openapi3.NewHeaderParameter("test").WithSchema(openapi3.NewArraySchema()),
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, got1 := mergeHeader(tt.args.header, tt.args.header2, tt.args.child)
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("mergeHeader() got = %v, want %v", got, tt.want)
//			}
//			if !reflect.DeepEqual(got1, tt.want1) {
//				t.Errorf("mergeHeader() got1 = %v, want %v", got1, tt.want1)
//			}
//		})
//	}
//}
//
//func Test_mergeResponseHeader(t *testing.T) {
//	type args struct {
//		headers  openapi3.Headers
//		headers2 openapi3.Headers
//		path     *field.Path
//	}
//	tests := []struct {
//		name  string
//		args  args
//		want  openapi3.Headers
//		want1 []conflict
//	}{
//		{
//			name: "first headers list empty",
//			args: args{
//				headers: openapi3.Headers{},
//				headers2: openapi3.Headers{
//					"test": createHeaderRef(openapi3.NewStringSchema()),
//				},
//				path: nil,
//			},
//			want: openapi3.Headers{
//				"test": createHeaderRef(openapi3.NewStringSchema()),
//			},
//			want1: nil,
//		},
//		{
//			name: "second headers list empty",
//			args: args{
//				headers: openapi3.Headers{
//					"test": createHeaderRef(openapi3.NewStringSchema()),
//				},
//				headers2: openapi3.Headers{},
//				path:     nil,
//			},
//			want: openapi3.Headers{
//				"test": createHeaderRef(openapi3.NewStringSchema()),
//			},
//			want1: nil,
//		},
//		{
//			name: "no common headers",
//			args: args{
//				headers: openapi3.Headers{
//					"test": createHeaderRef(openapi3.NewStringSchema()),
//				},
//				headers2: openapi3.Headers{
//					"test2": createHeaderRef(openapi3.NewStringSchema()),
//				},
//				path: nil,
//			},
//			want: openapi3.Headers{
//				"test":  createHeaderRef(openapi3.NewStringSchema()),
//				"test2": createHeaderRef(openapi3.NewStringSchema()),
//			},
//			want1: nil,
//		},
//		{
//			name: "merge mutual headers",
//			args: args{
//				headers: openapi3.Headers{
//					"test": createHeaderRef(openapi3.NewStringSchema()),
//				},
//				headers2: openapi3.Headers{
//					"test": createHeaderRef(openapi3.NewUUIDSchema()),
//				},
//				path: nil,
//			},
//			want: openapi3.Headers{
//				"test": createHeaderRef(openapi3.NewStringSchema()),
//			},
//			want1: nil,
//		},
//		{
//			name: "merge mutual headers and keep non mutual",
//			args: args{
//				headers: openapi3.Headers{
//					"mutual":     createHeaderRef(openapi3.NewStringSchema()),
//					"nonmutual1": createHeaderRef(openapi3.NewInt64Schema()),
//				},
//				headers2: openapi3.Headers{
//					"mutual":     createHeaderRef(openapi3.NewUUIDSchema()),
//					"nonmutual2": createHeaderRef(openapi3.NewBoolSchema()),
//				},
//				path: nil,
//			},
//			want: openapi3.Headers{
//				"mutual":     createHeaderRef(openapi3.NewStringSchema()),
//				"nonmutual1": createHeaderRef(openapi3.NewInt64Schema()),
//				"nonmutual2": createHeaderRef(openapi3.NewBoolSchema()),
//			},
//			want1: nil,
//		},
//		{
//			name: "merge mutual headers with conflicts",
//			args: args{
//				headers: openapi3.Headers{
//					"test": createHeaderRef(openapi3.NewInt64Schema()),
//				},
//				headers2: openapi3.Headers{
//					"test": createHeaderRef(openapi3.NewBoolSchema()),
//				},
//				path: field.NewPath("headers"),
//			},
//			want: openapi3.Headers{
//				"test": createHeaderRef(openapi3.NewInt64Schema()),
//			},
//			want1: []conflict{
//				{
//					path: field.NewPath("headers").Child("test"),
//					obj1: openapi3.NewInt64Schema(),
//					obj2: openapi3.NewBoolSchema(),
//					msg:  createConflictMsg(field.NewPath("headers").Child("test"), openapi3.TypeInteger, openapi3.TypeBoolean),
//				},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, got1 := mergeResponseHeader(tt.args.headers, tt.args.headers2, tt.args.path)
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("mergeResponseHeader() got = %+v, want %+v", got, tt.want)
//			}
//			if !reflect.DeepEqual(got1, tt.want1) {
//				t.Errorf("mergeResponseHeader() got1 = %v, want %v", got1, tt.want1)
//			}
//		})
//	}
//}
//
//func createHeaderRef(schema *openapi3.Schema) *openapi3.HeaderRef {
//	return &openapi3.HeaderRef{
//		Value: &openapi3.Header{
//			Parameter: openapi3.Parameter{
//				Schema: openapi3.NewSchemaRef("", schema),
//			},
//		},
//	}
//}
//
//func Test_mergeResponse(t *testing.T) {
//	type args struct {
//		response  *openapi3.Response
//		response2 *openapi3.Response
//		path      *field.Path
//	}
//	tests := []struct {
//		name  string
//		args  args
//		want  *openapi3.Response
//		want1 []conflict
//	}{
//		{
//			name: "first response is empty",
//			args: args{
//				response:  openapi3.NewResponse(),
//				response2: createTestResponse().WithHeader("X-Header", openapi3.NewStringSchema()).Response,
//				path:      nil,
//			},
//			want:  createTestResponse().WithHeader("X-Header", openapi3.NewStringSchema()).Response,
//			want1: nil,
//		},
//		{
//			name: "second response is empty",
//			args: args{
//				response:  createTestResponse().WithHeader("X-Header", openapi3.NewStringSchema()).Response,
//				response2: openapi3.NewResponse(),
//				path:      nil,
//			},
//			want:  createTestResponse().WithHeader("X-Header", openapi3.NewStringSchema()).Response,
//			want1: nil,
//		},
//		{
//			name: "merge response schema",
//			args: args{
//				response: createTestResponse().
//					WithJSONSchema(openapi3.NewDateTimeSchema()).
//					WithHeader("X-Header", openapi3.NewStringSchema()).Response,
//				response2: createTestResponse().
//					WithJSONSchema(openapi3.NewStringSchema()).
//					WithHeader("X-Header", openapi3.NewStringSchema()).Response,
//				path: nil,
//			},
//			want: createTestResponse().
//				WithJSONSchema(openapi3.NewStringSchema()).
//				WithHeader("X-Header", openapi3.NewStringSchema()).Response,
//			want1: nil,
//		},
//		{
//			name: "merge response header",
//			args: args{
//				response: createTestResponse().
//					WithJSONSchema(openapi3.NewStringSchema()).
//					WithHeader("X-Header", openapi3.NewUUIDSchema()).Response,
//				response2: createTestResponse().
//					WithJSONSchema(openapi3.NewStringSchema()).
//					WithHeader("X-Header", openapi3.NewStringSchema()).Response,
//				path: nil,
//			},
//			want: createTestResponse().
//				WithJSONSchema(openapi3.NewStringSchema()).
//				WithHeader("X-Header", openapi3.NewStringSchema()).Response,
//			want1: nil,
//		},
//		{
//			name: "merge response header and schema",
//			args: args{
//				response: createTestResponse().
//					WithJSONSchema(openapi3.NewDateTimeSchema()).
//					WithHeader("X-Header", openapi3.NewUUIDSchema()).Response,
//				response2: createTestResponse().
//					WithJSONSchema(openapi3.NewStringSchema()).
//					WithHeader("X-Header", openapi3.NewStringSchema()).Response,
//				path: nil,
//			},
//			want: createTestResponse().
//				WithJSONSchema(openapi3.NewStringSchema()).
//				WithHeader("X-Header", openapi3.NewStringSchema()).Response,
//			want1: nil,
//		},
//		{
//			name: "merge response header and schema prefer number",
//			args: args{
//				response: createTestResponse().
//					WithJSONSchema(openapi3.NewFloat64Schema()).
//					WithHeader("X-Header", openapi3.NewUUIDSchema()).Response,
//				response2: createTestResponse().
//					WithJSONSchema(openapi3.NewInt64Schema()).
//					WithHeader("X-Header", openapi3.NewBoolSchema()).Response,
//				path: field.NewPath("200"),
//			},
//			want: createTestResponse().
//				WithJSONSchema(openapi3.NewFloat64Schema()).
//				WithHeader("X-Header", openapi3.NewUUIDSchema()).Response,
//			want1: nil,
//		},
//		{
//			name: "merge response header and schema with conflicts",
//			args: args{
//				response: createTestResponse().
//					WithJSONSchema(openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())).
//					WithHeader("X-Header", openapi3.NewFloat64Schema()).Response,
//				response2: createTestResponse().
//					WithJSONSchema(openapi3.NewInt64Schema()).
//					WithHeader("X-Header", openapi3.NewBoolSchema()).Response,
//				path: field.NewPath("200"),
//			},
//			want: createTestResponse().
//				WithJSONSchema(openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())).
//				WithHeader("X-Header", openapi3.NewFloat64Schema()).Response,
//			want1: []conflict{
//				{
//					path: field.NewPath("200").Child("content").Child("application/json"),
//					obj1: openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema()),
//					obj2: openapi3.NewInt64Schema(),
//					msg: createConflictMsg(field.NewPath("200").Child("content").Child("application/json"),
//						openapi3.TypeArray, openapi3.TypeInteger),
//				},
//				{
//					path: field.NewPath("200").Child("headers").Child("X-Header"),
//					obj1: openapi3.NewFloat64Schema(),
//					obj2: openapi3.NewBoolSchema(),
//					msg: createConflictMsg(field.NewPath("200").Child("headers").Child("X-Header"),
//						openapi3.TypeNumber, openapi3.TypeBoolean),
//				},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, got1 := mergeResponse(tt.args.response, tt.args.response2, tt.args.path)
//			assertEqual(t, got, tt.want)
//			//assert.DeepEqual(t, got, tt.want, cmpopts.IgnoreUnexported(openapi3.Schema{}))
//			if !reflect.DeepEqual(got1, tt.want1) {
//				t.Errorf("mergeResponse() got1 = %v, want %v", got1, tt.want1)
//			}
//		})
//	}
//}
//
//func Test_mergeResponses(t *testing.T) {
//	type args struct {
//		responses  openapi3.Responses
//		responses2 openapi3.Responses
//		path       *field.Path
//	}
//	tests := []struct {
//		name  string
//		args  args
//		want  openapi3.Responses
//		want1 []conflict
//	}{
//		{
//			name: "first is nil",
//			args: args{
//				responses: openapi3.Responses{},
//				responses2: createTestResponses().
//					WithResponse("200", createTestResponse().
//						WithJSONSchema(openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())).
//						WithHeader("X-Header", openapi3.NewUUIDSchema()).Response).Responses,
//				path: nil,
//			},
//			want: createTestResponses().
//				WithResponse("200", createTestResponse().
//					WithJSONSchema(openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())).
//					WithHeader("X-Header", openapi3.NewUUIDSchema()).Response).Responses,
//			want1: nil,
//		},
//		{
//			name: "second is nil",
//			args: args{
//				responses: createTestResponses().
//					WithResponse("200", createTestResponse().
//						WithJSONSchema(openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())).
//						WithHeader("X-Header", openapi3.NewUUIDSchema()).Response).Responses,
//				responses2: openapi3.Responses{},
//				path:       nil,
//			},
//			want: createTestResponses().
//				WithResponse("200", createTestResponse().
//					WithJSONSchema(openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())).
//					WithHeader("X-Header", openapi3.NewUUIDSchema()).Response).Responses,
//			want1: nil,
//		},
//		{
//			name: "both are nil",
//			args: args{
//				responses:  openapi3.Responses{},
//				responses2: openapi3.Responses{},
//				path:       nil,
//			},
//			want:  openapi3.Responses{},
//			want1: nil,
//		},
//		{
//			name: "non mutual response code responses",
//			args: args{
//				responses: createTestResponses().
//					WithResponse("200", createTestResponse().
//						WithJSONSchema(openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())).
//						WithHeader("X-Header", openapi3.NewUUIDSchema()).Response).Responses,
//				responses2: createTestResponses().
//					WithResponse("201", createTestResponse().
//						WithJSONSchema(openapi3.NewStringSchema()).
//						WithHeader("X-Header2", openapi3.NewUUIDSchema()).Response).Responses,
//				path: nil,
//			},
//			want: createTestResponses().
//				WithResponse("200", createTestResponse().
//					WithJSONSchema(openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())).
//					WithHeader("X-Header", openapi3.NewUUIDSchema()).Response).
//				WithResponse("201", createTestResponse().
//					WithJSONSchema(openapi3.NewStringSchema()).
//					WithHeader("X-Header2", openapi3.NewUUIDSchema()).Response).Responses,
//			want1: nil,
//		},
//		{
//			name: "mutual response code responses",
//			args: args{
//				responses: createTestResponses().
//					WithResponse("200", createTestResponse().
//						WithJSONSchema(openapi3.NewArraySchema().WithItems(openapi3.NewDateTimeSchema())).
//						WithHeader("X-Header", openapi3.NewUUIDSchema()).Response).Responses,
//				responses2: createTestResponses().
//					WithResponse("200", createTestResponse().
//						WithJSONSchema(openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())).
//						WithHeader("X-Header", openapi3.NewUUIDSchema()).Response).Responses,
//				path: nil,
//			},
//			want: createTestResponses().
//				WithResponse("200", createTestResponse().
//					WithJSONSchema(openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())).
//					WithHeader("X-Header", openapi3.NewUUIDSchema()).Response).Responses,
//			want1: nil,
//		},
//		{
//			name: "mutual and non mutual response code responses",
//			args: args{
//				responses: createTestResponses().
//					WithResponse("200", createTestResponse().
//						WithJSONSchema(openapi3.NewArraySchema().WithItems(openapi3.NewDateTimeSchema())).
//						WithHeader("X-Header", openapi3.NewUUIDSchema()).Response).
//					WithResponse("201", createTestResponse().
//						WithJSONSchema(openapi3.NewDateTimeSchema()).
//						WithHeader("X-Header1", openapi3.NewUUIDSchema()).Response).Responses,
//				responses2: createTestResponses().
//					WithResponse("200", createTestResponse().
//						WithJSONSchema(openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())).
//						WithHeader("X-Header", openapi3.NewUUIDSchema()).Response).
//					WithResponse("202", createTestResponse().
//						WithJSONSchema(openapi3.NewBoolSchema()).
//						WithHeader("X-Header3", openapi3.NewUUIDSchema()).Response).Responses,
//				path: nil,
//			},
//			want: createTestResponses().
//				WithResponse("200", createTestResponse().
//					WithJSONSchema(openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())).
//					WithHeader("X-Header", openapi3.NewUUIDSchema()).Response).
//				WithResponse("201", createTestResponse().
//					WithJSONSchema(openapi3.NewDateTimeSchema()).
//					WithHeader("X-Header1", openapi3.NewUUIDSchema()).Response).
//				WithResponse("202", createTestResponse().
//					WithJSONSchema(openapi3.NewBoolSchema()).
//					WithHeader("X-Header3", openapi3.NewUUIDSchema()).Response).Responses,
//			want1: nil,
//		},
//		{
//			name: "mutual and non mutual response code responses solve conflicts",
//			args: args{
//				responses: createTestResponses().
//					WithResponse("200", createTestResponse().
//						WithJSONSchema(openapi3.NewArraySchema().WithItems(openapi3.NewArraySchema().WithItems(openapi3.NewDateTimeSchema()))).
//						WithHeader("X-Header", openapi3.NewUUIDSchema()).Response).
//					WithResponse("201", createTestResponse().
//						WithJSONSchema(openapi3.NewDateTimeSchema()).
//						WithHeader("X-Header1", openapi3.NewUUIDSchema()).Response).Responses,
//				responses2: createTestResponses().
//					WithResponse("200", createTestResponse().
//						WithJSONSchema(openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())).
//						WithHeader("X-Header", openapi3.NewUUIDSchema()).Response).
//					WithResponse("202", createTestResponse().
//						WithJSONSchema(openapi3.NewBoolSchema()).
//						WithHeader("X-Header3", openapi3.NewUUIDSchema()).Response).Responses,
//				path: field.NewPath("responses"),
//			},
//			want: createTestResponses().
//				WithResponse("200", createTestResponse().
//					WithJSONSchema(openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())).
//					WithHeader("X-Header", openapi3.NewUUIDSchema()).Response).
//				WithResponse("201", createTestResponse().
//					WithJSONSchema(openapi3.NewDateTimeSchema()).
//					WithHeader("X-Header1", openapi3.NewUUIDSchema()).Response).
//				WithResponse("202", createTestResponse().
//					WithJSONSchema(openapi3.NewBoolSchema()).
//					WithHeader("X-Header3", openapi3.NewUUIDSchema()).Response).Responses,
//			want1: nil,
//		},
//		{
//			name: "mutual and non mutual response code responses with conflicts",
//			args: args{
//				responses: createTestResponses().
//					WithResponse("200", createTestResponse().
//						WithJSONSchema(openapi3.NewArraySchema().WithItems(openapi3.NewArraySchema().WithItems(openapi3.NewDateTimeSchema()))).
//						WithHeader("X-Header", openapi3.NewUUIDSchema()).Response).
//					WithResponse("201", createTestResponse().
//						WithJSONSchema(openapi3.NewDateTimeSchema()).
//						WithHeader("X-Header1", openapi3.NewUUIDSchema()).Response).Responses,
//				responses2: createTestResponses().
//					WithResponse("200", createTestResponse().
//						WithJSONSchema(openapi3.NewArraySchema().WithItems(openapi3.NewFloat64Schema())).
//						WithHeader("X-Header", openapi3.NewUUIDSchema()).Response).
//					WithResponse("202", createTestResponse().
//						WithJSONSchema(openapi3.NewBoolSchema()).
//						WithHeader("X-Header3", openapi3.NewUUIDSchema()).Response).Responses,
//				path: field.NewPath("responses"),
//			},
//			want: createTestResponses().
//				WithResponse("200", createTestResponse().
//					WithJSONSchema(openapi3.NewArraySchema().WithItems(openapi3.NewArraySchema().WithItems(openapi3.NewDateTimeSchema()))).
//					WithHeader("X-Header", openapi3.NewUUIDSchema()).Response).
//				WithResponse("201", createTestResponse().
//					WithJSONSchema(openapi3.NewDateTimeSchema()).
//					WithHeader("X-Header1", openapi3.NewUUIDSchema()).Response).
//				WithResponse("202", createTestResponse().
//					WithJSONSchema(openapi3.NewBoolSchema()).
//					WithHeader("X-Header3", openapi3.NewUUIDSchema()).Response).Responses,
//			want1: []conflict{
//				{
//					path: field.NewPath("responses").Child("200").Child("content").
//						Child("application/json").Child("items"),
//					obj1: openapi3.NewArraySchema().WithItems(openapi3.NewDateTimeSchema()),
//					obj2: openapi3.NewFloat64Schema(),
//					msg: createConflictMsg(field.NewPath("responses").Child("200").Child("content").
//						Child("application/json").Child("items"), openapi3.TypeArray, openapi3.TypeNumber),
//				},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, got1 := mergeResponses(tt.args.responses, tt.args.responses2, tt.args.path)
//			//assert.DeepEqual(t, got, tt.want, cmpopts.IgnoreUnexported(openapi3.Schema{}))
//			assertEqual(t, got, tt.want)
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("mergeResponses() got = %v, want %v", got, tt.want)
//			}
//			if !reflect.DeepEqual(got1, tt.want1) {
//				t.Errorf("mergeResponses() got1 = %v, want %v", got1, tt.want1)
//			}
//		})
//	}
//}
//
//func Test_mergeProperties(t *testing.T) {
//	type args struct {
//		properties  openapi3.Schemas
//		properties2 openapi3.Schemas
//		path        *field.Path
//	}
//	tests := []struct {
//		name  string
//		args  args
//		want  openapi3.Schemas
//		want1 []conflict
//	}{
//		{
//			name: "first is nil",
//			args: args{
//				properties: nil,
//				properties2: openapi3.Schemas{
//					"string-key": openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
//				},
//				path: nil,
//			},
//			want: openapi3.Schemas{
//				"string-key": openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
//			},
//			want1: nil,
//		},
//		{
//			name: "second is nil",
//			args: args{
//				properties: openapi3.Schemas{
//					"string-key": openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
//				},
//				properties2: nil,
//				path:        nil,
//			},
//			want: openapi3.Schemas{
//				"string-key": openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
//			},
//			want1: nil,
//		},
//		{
//			name: "both are nil",
//			args: args{
//				properties:  nil,
//				properties2: nil,
//				path:        nil,
//			},
//			want:  make(openapi3.Schemas),
//			want1: nil,
//		},
//		{
//			name: "non mutual properties",
//			args: args{
//				properties: openapi3.Schemas{
//					"string-key": openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
//				},
//				properties2: openapi3.Schemas{
//					"bool-key": openapi3.NewSchemaRef("", openapi3.NewBoolSchema()),
//				},
//				path: nil,
//			},
//			want: openapi3.Schemas{
//				"string-key": openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
//				"bool-key":   openapi3.NewSchemaRef("", openapi3.NewBoolSchema()),
//			},
//			want1: nil,
//		},
//		{
//			name: "mutual properties",
//			args: args{
//				properties: openapi3.Schemas{
//					"string-key": openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
//				},
//				properties2: openapi3.Schemas{
//					"string-key": openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
//				},
//				path: nil,
//			},
//			want: openapi3.Schemas{
//				"string-key": openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
//			},
//			want1: nil,
//		},
//		{
//			name: "mutual and non mutual properties",
//			args: args{
//				properties: openapi3.Schemas{
//					"string-key": openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
//					"bool-key":   openapi3.NewSchemaRef("", openapi3.NewBoolSchema()),
//				},
//				properties2: openapi3.Schemas{
//					"string-key": openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
//					"int-key":    openapi3.NewSchemaRef("", openapi3.NewInt64Schema()),
//				},
//				path: nil,
//			},
//			want: openapi3.Schemas{
//				"string-key": openapi3.NewSchemaRef("", openapi3.NewStringSchema()),
//				"int-key":    openapi3.NewSchemaRef("", openapi3.NewInt64Schema()),
//				"bool-key":   openapi3.NewSchemaRef("", openapi3.NewBoolSchema()),
//			},
//			want1: nil,
//		},
//		{
//			name: "mutual and non mutual response code responses with conflicts",
//			args: args{
//				properties: openapi3.Schemas{
//					"conflict": openapi3.NewSchemaRef("", openapi3.NewBoolSchema()),
//					"bool-key": openapi3.NewSchemaRef("", openapi3.NewBoolSchema()),
//				},
//				properties2: openapi3.Schemas{
//					"conflict": openapi3.NewSchemaRef("", openapi3.NewInt64Schema()),
//					"int-key":  openapi3.NewSchemaRef("", openapi3.NewInt64Schema()),
//				},
//				path: field.NewPath("properties"),
//			},
//			want: openapi3.Schemas{
//				"conflict": openapi3.NewSchemaRef("", openapi3.NewBoolSchema()),
//				"int-key":  openapi3.NewSchemaRef("", openapi3.NewInt64Schema()),
//				"bool-key": openapi3.NewSchemaRef("", openapi3.NewBoolSchema()),
//			},
//			want1: []conflict{
//				{
//					path: field.NewPath("properties").Child("conflict"),
//					obj1: openapi3.NewBoolSchema(),
//					obj2: openapi3.NewInt64Schema(),
//					msg: createConflictMsg(field.NewPath("properties").Child("conflict"),
//						openapi3.TypeBoolean, openapi3.TypeInteger),
//				},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, got1 := mergeProperties(tt.args.properties, tt.args.properties2, tt.args.path)
//			//assert.DeepEqual(t, got, tt.want, cmpopts.IgnoreUnexported(openapi3.Schema{}))
//			assertEqual(t, got, tt.want)
//			if !reflect.DeepEqual(got1, tt.want1) {
//				t.Errorf("mergeProperties() got1 = %v, want %v", got1, tt.want1)
//			}
//		})
//	}
//}
//
//func Test_mergeSchemaItems(t *testing.T) {
//	type args struct {
//		items  *openapi3.SchemaRef
//		items2 *openapi3.SchemaRef
//		path   *field.Path
//	}
//	tests := []struct {
//		name  string
//		args  args
//		want  *openapi3.SchemaRef
//		want1 []conflict
//	}{
//		{
//			name: "no merge needed",
//			args: args{
//				items:  openapi3.NewSchemaRef("", openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())),
//				items2: openapi3.NewSchemaRef("", openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())),
//				path:   field.NewPath("test"),
//			},
//			want:  openapi3.NewSchemaRef("", openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())),
//			want1: nil,
//		},
//		{
//			name: "items with string format - format should be removed",
//			args: args{
//				items:  openapi3.NewSchemaRef("", openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())),
//				items2: openapi3.NewSchemaRef("", openapi3.NewArraySchema().WithItems(openapi3.NewUUIDSchema())),
//				path:   field.NewPath("test"),
//			},
//			want:  openapi3.NewSchemaRef("", openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())),
//			want1: nil,
//		},
//		{
//			name: "different type of items",
//			args: args{
//				items:  openapi3.NewSchemaRef("", openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())),
//				items2: openapi3.NewSchemaRef("", openapi3.NewInt64Schema()),
//				path:   field.NewPath("test"),
//			},
//			want: openapi3.NewSchemaRef("", openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())),
//			want1: []conflict{
//				{
//					path: field.NewPath("test").Child("items"),
//					obj1: openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema()),
//					obj2: openapi3.NewInt64Schema(),
//					msg:  createConflictMsg(field.NewPath("test").Child("items"), openapi3.TypeArray, openapi3.TypeInteger),
//				},
//			},
//		},
//		{
//			name: "items2 nil items - expected to get items",
//			args: args{
//				items:  openapi3.NewSchemaRef("", openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())),
//				items2: nil,
//				path:   field.NewPath("test"),
//			},
//			want:  openapi3.NewSchemaRef("", openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())),
//			want1: nil,
//		},
//		{
//			name: "items2 nil schema - expected to get items",
//			args: args{
//				items:  openapi3.NewSchemaRef("", openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())),
//				items2: openapi3.NewSchemaRef("", openapi3.NewArraySchema().WithItems(nil)),
//				path:   field.NewPath("test"),
//			},
//			want:  openapi3.NewSchemaRef("", openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())),
//			want1: nil,
//		},
//		{
//			name: "items nil items - expected to get items2",
//			args: args{
//				items:  nil,
//				items2: openapi3.NewSchemaRef("", openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())),
//				path:   field.NewPath("test"),
//			},
//			want:  openapi3.NewSchemaRef("", openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())),
//			want1: nil,
//		},
//		{
//			name: "both schemas nil items - expected to get schema",
//			args: args{
//				items:  nil,
//				items2: nil,
//				path:   field.NewPath("test"),
//			},
//			want:  nil,
//			want1: nil,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, got1 := mergeSchemaItems(tt.args.items, tt.args.items2, tt.args.path)
//			//assert.DeepEqual(t, got, tt.want, cmpopts.IgnoreUnexported(openapi3.Schema{}))
//			assertEqual(t, got, tt.want)
//			if !reflect.DeepEqual(got1, tt.want1) {
//				t.Errorf("mergeSchemaItems() got1 = %v, want %v", got1, tt.want1)
//			}
//		})
//	}
//}
//
//func Test_mergeSchema(t *testing.T) {
//	emptySchemaType := openapi3.NewSchema()
//	type args struct {
//		schema  *openapi3.Schema
//		schema2 *openapi3.Schema
//		path    *field.Path
//	}
//	tests := []struct {
//		name  string
//		args  args
//		want  *openapi3.Schema
//		want1 []conflict
//	}{
//		{
//			name: "no merge needed",
//			args: args{
//				schema:  openapi3.NewInt64Schema(),
//				schema2: openapi3.NewInt64Schema(),
//				path:    nil,
//			},
//			want:  openapi3.NewInt64Schema(),
//			want1: nil,
//		},
//		{
//			name: "first is nil",
//			args: args{
//				schema:  nil,
//				schema2: openapi3.NewInt64Schema(),
//				path:    nil,
//			},
//			want:  openapi3.NewInt64Schema(),
//			want1: nil,
//		},
//		{
//			name: "second is nil",
//			args: args{
//				schema:  openapi3.NewInt64Schema(),
//				schema2: nil,
//				path:    nil,
//			},
//			want:  openapi3.NewInt64Schema(),
//			want1: nil,
//		},
//		{
//			name: "both are nil",
//			args: args{
//				schema:  nil,
//				schema2: nil,
//				path:    nil,
//			},
//			want:  nil,
//			want1: nil,
//		},
//		{
//			name: "first has empty schema type",
//			args: args{
//				schema:  emptySchemaType,
//				schema2: openapi3.NewInt64Schema(),
//				path:    nil,
//			},
//			want:  openapi3.NewInt64Schema(),
//			want1: nil,
//		},
//		{
//			name: "second has empty schema type",
//			args: args{
//				schema:  openapi3.NewInt64Schema(),
//				schema2: emptySchemaType,
//				path:    nil,
//			},
//			want:  openapi3.NewInt64Schema(),
//			want1: nil,
//		},
//		{
//			name: "both has empty schema type",
//			args: args{
//				schema:  emptySchemaType,
//				schema2: emptySchemaType,
//				path:    nil,
//			},
//			want:  emptySchemaType,
//			want1: nil,
//		},
//		{
//			name: "type conflict",
//			args: args{
//				schema:  openapi3.NewInt64Schema(),
//				schema2: openapi3.NewBoolSchema(),
//				path:    field.NewPath("schema"),
//			},
//			want: openapi3.NewInt64Schema(),
//			want1: []conflict{
//				{
//					path: field.NewPath("schema"),
//					obj1: openapi3.NewInt64Schema(),
//					obj2: openapi3.NewBoolSchema(),
//					msg:  createConflictMsg(field.NewPath("schema"), openapi3.TypeInteger, openapi3.TypeBoolean),
//				},
//			},
//		},
//		{
//			name: "string type with different format - dismiss the format",
//			args: args{
//				schema:  openapi3.NewDateTimeSchema(),
//				schema2: openapi3.NewUUIDSchema(),
//				path:    field.NewPath("schema"),
//			},
//			want:  openapi3.NewStringSchema(),
//			want1: nil,
//		},
//		{
//			name: "array conflict prefer number",
//			args: args{
//				schema:  openapi3.NewArraySchema().WithItems(openapi3.NewInt64Schema()),
//				schema2: openapi3.NewArraySchema().WithItems(openapi3.NewFloat64Schema()),
//				path:    field.NewPath("schema"),
//			},
//			want:  openapi3.NewArraySchema().WithItems(openapi3.NewFloat64Schema()),
//			want1: nil,
//		},
//		{
//			name: "array conflict",
//			args: args{
//				schema:  openapi3.NewArraySchema().WithItems(openapi3.NewInt64Schema()),
//				schema2: openapi3.NewArraySchema().WithItems(openapi3.NewObjectSchema()),
//				path:    field.NewPath("schema"),
//			},
//			want: openapi3.NewArraySchema().WithItems(openapi3.NewInt64Schema()),
//			want1: []conflict{
//				{
//					path: field.NewPath("schema").Child("items"),
//					obj1: openapi3.NewInt64Schema(),
//					obj2: openapi3.NewObjectSchema(),
//					msg:  createConflictMsg(field.NewPath("schema").Child("items"), openapi3.TypeInteger, openapi3.TypeObject),
//				},
//			},
//		},
//		{
//			name: "merge object with conflict",
//			args: args{
//				schema: openapi3.NewObjectSchema().
//					WithProperty("bool", openapi3.NewBoolSchema()).
//					WithProperty("conflict prefer string", openapi3.NewBoolSchema()).
//					WithProperty("conflict", openapi3.NewObjectSchema()),
//				schema2: openapi3.NewObjectSchema().
//					WithProperty("float", openapi3.NewFloat64Schema()).
//					WithProperty("conflict prefer string", openapi3.NewStringSchema()).
//					WithProperty("conflict", openapi3.NewInt64Schema()),
//				path: field.NewPath("schema"),
//			},
//			want: openapi3.NewObjectSchema().
//				WithProperty("bool", openapi3.NewBoolSchema()).
//				WithProperty("conflict prefer string", openapi3.NewStringSchema()).
//				WithProperty("conflict", openapi3.NewObjectSchema()).
//				WithProperty("float", openapi3.NewFloat64Schema()),
//			want1: []conflict{
//				{
//					path: field.NewPath("schema").Child("properties").Child("conflict"),
//					obj1: openapi3.NewObjectSchema(),
//					obj2: openapi3.NewInt64Schema(),
//					msg: createConflictMsg(field.NewPath("schema").Child("properties").Child("conflict"),
//						openapi3.TypeObject, openapi3.TypeInteger),
//				},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, got1 := mergeSchema(tt.args.schema, tt.args.schema2, tt.args.path)
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("mergeSchema() got = %v, want %v", got, tt.want)
//			}
//			if !reflect.DeepEqual(got1, tt.want1) {
//				t.Errorf("mergeSchema() got1 = %v, want %v", got1, tt.want1)
//			}
//		})
//	}
//}
//
//func Test_mergeParameter(t *testing.T) {
//	type args struct {
//		parameter  *openapi3.Parameter
//		parameter2 *openapi3.Parameter
//		path       *field.Path
//	}
//	tests := []struct {
//		name  string
//		args  args
//		want  *openapi3.Parameter
//		want1 []conflict
//	}{
//		{
//			name: "param type solve conflict",
//			args: args{
//				parameter:  openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewStringSchema()),
//				parameter2: openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewBoolSchema()),
//				path:       field.NewPath("param-name"),
//			},
//			want:  openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewStringSchema()),
//			want1: nil,
//		},
//		{
//			name: "param type conflict",
//			args: args{
//				parameter:  openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewInt64Schema()),
//				parameter2: openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewBoolSchema()),
//				path:       field.NewPath("param-name"),
//			},
//			want: openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewInt64Schema()),
//			want1: []conflict{
//				{
//					path: field.NewPath("param-name"),
//					obj1: openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewInt64Schema()),
//					obj2: openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewBoolSchema()),
//					msg:  createConflictMsg(field.NewPath("param-name"), openapi3.TypeInteger, openapi3.TypeBoolean),
//				},
//			},
//		},
//		{
//			name: "string merge",
//			args: args{
//				parameter:  openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewStringSchema()),
//				parameter2: openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewUUIDSchema()),
//				path:       field.NewPath("param-name"),
//			},
//			want:  openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewStringSchema()),
//			want1: nil,
//		},
//		{
//			name: "array merge with conflict",
//			args: args{
//				parameter:  openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())),
//				parameter2: openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewArraySchema().WithItems(openapi3.NewBoolSchema())),
//				path:       field.NewPath("param-name"),
//			},
//			want:  openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewArraySchema().WithItems(openapi3.NewStringSchema())),
//			want1: nil,
//		},
//		{
//			name: "object merge",
//			args: args{
//				parameter:  openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewObjectSchema().WithProperty("string", openapi3.NewStringSchema())),
//				parameter2: openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewObjectSchema().WithProperty("bool", openapi3.NewBoolSchema())),
//				path:       field.NewPath("param-name"),
//			},
//			want: openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewObjectSchema().
//				WithProperty("bool", openapi3.NewBoolSchema()).
//				WithProperty("string", openapi3.NewStringSchema())),
//			want1: nil,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, got1 := mergeParameter(tt.args.parameter, tt.args.parameter2, tt.args.path)
//			//assert.DeepEqual(t, got, tt.want, cmpopts.IgnoreUnexported(openapi3.Schema{}))
//			assertEqual(t, got, tt.want)
//			if !reflect.DeepEqual(got1, tt.want1) {
//				t.Errorf("mergeParameter() got1 = %v, want %v", got1, tt.want1)
//			}
//		})
//	}
//}
//
//func Test_makeParametersMapByName(t *testing.T) {
//	type args struct {
//		parameters openapi3.Parameters
//	}
//	tests := []struct {
//		name string
//		args args
//		want map[string]*openapi3.ParameterRef
//	}{
//		{
//			name: "sanity",
//			args: args{
//				parameters: openapi3.Parameters{
//					&openapi3.ParameterRef{Value: openapi3.NewHeaderParameter("header")},
//					&openapi3.ParameterRef{Value: openapi3.NewHeaderParameter("header2")},
//					&openapi3.ParameterRef{Value: openapi3.NewPathParameter("path")},
//					&openapi3.ParameterRef{Value: openapi3.NewPathParameter("path2")},
//				},
//			},
//			want: map[string]*openapi3.ParameterRef{
//				"header":  {Value: openapi3.NewHeaderParameter("header")},
//				"header2": {Value: openapi3.NewHeaderParameter("header2")},
//				"path":    {Value: openapi3.NewPathParameter("path")},
//				"path2":   {Value: openapi3.NewPathParameter("path2")},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := makeParametersMapByName(tt.args.parameters); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("makeParametersMapByName() = %v, want %v", marshal(got), marshal(tt.want))
//			}
//		})
//	}
//}
//
//func Test_mergeParametersByInType(t *testing.T) {
//	type args struct {
//		parameters  openapi3.Parameters
//		parameters2 openapi3.Parameters
//		path        *field.Path
//	}
//	tests := []struct {
//		name  string
//		args  args
//		want  openapi3.Parameters
//		want1 []conflict
//	}{
//		{
//			name: "first is nil",
//			args: args{
//				parameters:  nil,
//				parameters2: openapi3.Parameters{{Value: openapi3.NewHeaderParameter("h")}},
//				path:        nil,
//			},
//			want:  openapi3.Parameters{{Value: openapi3.NewHeaderParameter("h")}},
//			want1: nil,
//		},
//		{
//			name: "second is nil",
//			args: args{
//				parameters:  openapi3.Parameters{{Value: openapi3.NewHeaderParameter("h")}},
//				parameters2: nil,
//				path:        nil,
//			},
//			want:  openapi3.Parameters{{Value: openapi3.NewHeaderParameter("h")}},
//			want1: nil,
//		},
//		{
//			name: "both are nil",
//			args: args{
//				parameters:  nil,
//				parameters2: nil,
//				path:        nil,
//			},
//			want:  nil,
//			want1: nil,
//		},
//		{
//			name: "non mutual parameters",
//			args: args{
//				parameters:  openapi3.Parameters{{Value: openapi3.NewHeaderParameter("X-Header-1")}},
//				parameters2: openapi3.Parameters{{Value: openapi3.NewHeaderParameter("X-Header-2")}},
//				path:        nil,
//			},
//			want:  openapi3.Parameters{{Value: openapi3.NewHeaderParameter("X-Header-1")}, {Value: openapi3.NewHeaderParameter("X-Header-2")}},
//			want1: nil,
//		},
//		{
//			name: "mutual parameters",
//			args: args{
//				parameters:  openapi3.Parameters{{Value: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewUUIDSchema())}},
//				parameters2: openapi3.Parameters{{Value: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewStringSchema())}},
//				path:        nil,
//			},
//			want:  openapi3.Parameters{{Value: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewStringSchema())}},
//			want1: nil,
//		},
//		{
//			name: "mutual and non mutual parameters",
//			args: args{
//				parameters: openapi3.Parameters{
//					{Value: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewUUIDSchema())},
//					{Value: openapi3.NewHeaderParameter("X-Header-2").WithSchema(openapi3.NewBoolSchema())},
//				},
//				parameters2: openapi3.Parameters{
//					{Value: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewStringSchema())},
//					{Value: openapi3.NewHeaderParameter("X-Header-3").WithSchema(openapi3.NewInt64Schema())},
//				},
//				path: nil,
//			},
//			want: openapi3.Parameters{
//				{Value: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewStringSchema())},
//				{Value: openapi3.NewHeaderParameter("X-Header-2").WithSchema(openapi3.NewBoolSchema())},
//				{Value: openapi3.NewHeaderParameter("X-Header-3").WithSchema(openapi3.NewInt64Schema())},
//			},
//			want1: nil,
//		},
//		{
//			name: "mutual and non mutual parameters solve conflicts",
//			args: args{
//				parameters: openapi3.Parameters{
//					{Value: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewUUIDSchema())},
//					{Value: openapi3.NewHeaderParameter("X-Header-2").WithSchema(openapi3.NewBoolSchema())},
//				},
//				parameters2: openapi3.Parameters{
//					{Value: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewBoolSchema())},
//					{Value: openapi3.NewHeaderParameter("X-Header-3").WithSchema(openapi3.NewInt64Schema())},
//				},
//				path: field.NewPath("parameters"),
//			},
//			want: openapi3.Parameters{
//				{Value: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewUUIDSchema())},
//				{Value: openapi3.NewHeaderParameter("X-Header-2").WithSchema(openapi3.NewBoolSchema())},
//				{Value: openapi3.NewHeaderParameter("X-Header-3").WithSchema(openapi3.NewInt64Schema())},
//			},
//			want1: nil,
//		},
//		{
//			name: "mutual and non mutual parameters with conflicts",
//			args: args{
//				parameters: openapi3.Parameters{
//					{Value: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewInt64Schema())},
//					{Value: openapi3.NewHeaderParameter("X-Header-2").WithSchema(openapi3.NewBoolSchema())},
//				},
//				parameters2: openapi3.Parameters{
//					{Value: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewBoolSchema())},
//					{Value: openapi3.NewHeaderParameter("X-Header-3").WithSchema(openapi3.NewInt64Schema())},
//				},
//				path: field.NewPath("parameters"),
//			},
//			want: openapi3.Parameters{
//				{Value: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewInt64Schema())},
//				{Value: openapi3.NewHeaderParameter("X-Header-2").WithSchema(openapi3.NewBoolSchema())},
//				{Value: openapi3.NewHeaderParameter("X-Header-3").WithSchema(openapi3.NewInt64Schema())},
//			},
//			want1: []conflict{
//				{
//					path: field.NewPath("parameters").Child("X-Header-1"),
//					obj1: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewInt64Schema()),
//					obj2: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewBoolSchema()),
//					msg: createConflictMsg(field.NewPath("parameters").Child("X-Header-1"), openapi3.TypeInteger,
//						openapi3.TypeBoolean),
//				},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, got1 := mergeParametersByInType(tt.args.parameters, tt.args.parameters2, tt.args.path)
//			sortParam(got)
//			sortParam(tt.want)
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("mergeParametersByInType() got = %v, want %v", got, tt.want)
//			}
//			if !reflect.DeepEqual(got1, tt.want1) {
//				t.Errorf("mergeParametersByInType() got1 = %v, want %v", got1, tt.want1)
//			}
//		})
//	}
//}
//
//func Test_getParametersByIn(t *testing.T) {
//	type args struct {
//		parameters openapi3.Parameters
//	}
//	tests := []struct {
//		name string
//		args args
//		want map[string]openapi3.Parameters
//	}{
//		{
//			name: "sanity",
//			args: args{
//				parameters: openapi3.Parameters{
//					{Value: openapi3.NewHeaderParameter("h1")},
//					{Value: openapi3.NewHeaderParameter("h2")},
//					{Value: openapi3.NewPathParameter("p1")},
//					{Value: openapi3.NewPathParameter("p2")},
//					{Value: openapi3.NewQueryParameter("q1")},
//					{Value: openapi3.NewQueryParameter("q2")},
//					{Value: openapi3.NewCookieParameter("c1")},
//					{Value: openapi3.NewCookieParameter("c2")},
//					{Value: &openapi3.Parameter{In: "not-supported"}},
//				},
//			},
//			want: map[string]openapi3.Parameters{
//				openapi3.ParameterInCookie: {{Value: openapi3.NewCookieParameter("c1")}, {Value: openapi3.NewCookieParameter("c2")}},
//				openapi3.ParameterInHeader: {{Value: openapi3.NewHeaderParameter("h1")}, {Value: openapi3.NewHeaderParameter("h2")}},
//				openapi3.ParameterInQuery:  {{Value: openapi3.NewQueryParameter("q1")}, {Value: openapi3.NewQueryParameter("q2")}},
//				openapi3.ParameterInPath:   {{Value: openapi3.NewPathParameter("p1")}, {Value: openapi3.NewPathParameter("p2")}},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := getParametersByIn(tt.args.parameters); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("getParametersByIn() = %v, want %v", marshal(got), marshal(tt.want))
//			}
//		})
//	}
//}
//
//func Test_mergeParameters(t *testing.T) {
//	type args struct {
//		parameters  openapi3.Parameters
//		parameters2 openapi3.Parameters
//		path        *field.Path
//	}
//	tests := []struct {
//		name  string
//		args  args
//		want  openapi3.Parameters
//		want1 []conflict
//	}{
//		{
//			name: "first is nil",
//			args: args{
//				parameters:  nil,
//				parameters2: openapi3.Parameters{{Value: openapi3.NewHeaderParameter("h")}},
//				path:        nil,
//			},
//			want:  openapi3.Parameters{{Value: openapi3.NewHeaderParameter("h")}},
//			want1: nil,
//		},
//		{
//			name: "second is nil",
//			args: args{
//				parameters:  openapi3.Parameters{{Value: openapi3.NewHeaderParameter("h")}},
//				parameters2: nil,
//				path:        nil,
//			},
//			want:  openapi3.Parameters{{Value: openapi3.NewHeaderParameter("h")}},
//			want1: nil,
//		},
//		{
//			name: "both are nil",
//			args: args{
//				parameters:  nil,
//				parameters2: nil,
//				path:        nil,
//			},
//			want:  nil,
//			want1: nil,
//		},
//		{
//			name: "non mutual parameters",
//			args: args{
//				parameters: openapi3.Parameters{
//					{Value: openapi3.NewHeaderParameter("X-Header-1")},
//					{Value: openapi3.NewQueryParameter("query-1")},
//				},
//				parameters2: openapi3.Parameters{
//					{Value: openapi3.NewHeaderParameter("X-Header-2")},
//					{Value: openapi3.NewQueryParameter("query-2")},
//					{Value: openapi3.NewHeaderParameter("header")},
//				},
//				path: nil,
//			},
//			want: openapi3.Parameters{
//				{Value: openapi3.NewHeaderParameter("X-Header-1")},
//				{Value: openapi3.NewQueryParameter("query-1")},
//				{Value: openapi3.NewHeaderParameter("header")},
//				{Value: openapi3.NewHeaderParameter("X-Header-2")},
//				{Value: openapi3.NewQueryParameter("query-2")},
//			},
//			want1: nil,
//		},
//		{
//			name: "mutual parameters",
//			args: args{
//				parameters: openapi3.Parameters{
//					{Value: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewUUIDSchema())},
//					{Value: openapi3.NewQueryParameter("query-1").WithSchema(openapi3.NewUUIDSchema())},
//					{Value: openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewObjectSchema().WithProperty("str", openapi3.NewStringSchema()))},
//				},
//				parameters2: openapi3.Parameters{
//					{Value: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewUUIDSchema())},
//					{Value: openapi3.NewQueryParameter("query-1").WithSchema(openapi3.NewUUIDSchema())},
//					{Value: openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewObjectSchema().WithProperty("str", openapi3.NewDateTimeSchema()))},
//				},
//				path: nil,
//			},
//			want: openapi3.Parameters{
//				{Value: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewUUIDSchema())},
//				{Value: openapi3.NewQueryParameter("query-1").WithSchema(openapi3.NewUUIDSchema())},
//				{Value: openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewObjectSchema().WithProperty("str", openapi3.NewStringSchema()))},
//			},
//			want1: nil,
//		},
//		{
//			name: "mutual and non mutual parameters",
//			args: args{
//				parameters: openapi3.Parameters{
//					{Value: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewUUIDSchema())},
//					{Value: openapi3.NewQueryParameter("query-1").WithSchema(openapi3.NewUUIDSchema())},
//					{Value: openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewObjectSchema().WithProperty("str", openapi3.NewStringSchema()))},
//					{Value: openapi3.NewPathParameter("non-mutual-1").WithSchema(openapi3.NewStringSchema())},
//					{Value: openapi3.NewCookieParameter("non-mutual-2").WithSchema(openapi3.NewStringSchema())},
//				},
//				parameters2: openapi3.Parameters{
//					{Value: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewUUIDSchema())},
//					{Value: openapi3.NewQueryParameter("query-1").WithSchema(openapi3.NewUUIDSchema())},
//					{Value: openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewObjectSchema().WithProperty("str", openapi3.NewDateTimeSchema()))},
//					{Value: openapi3.NewPathParameter("non-mutual-3").WithSchema(openapi3.NewStringSchema())},
//					{Value: openapi3.NewCookieParameter("non-mutual-4").WithSchema(openapi3.NewStringSchema())},
//				},
//				path: nil,
//			},
//			want: openapi3.Parameters{
//				{Value: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewUUIDSchema())},
//				{Value: openapi3.NewQueryParameter("query-1").WithSchema(openapi3.NewUUIDSchema())},
//				{Value: openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewObjectSchema().WithProperty("str", openapi3.NewStringSchema()))},
//				{Value: openapi3.NewPathParameter("non-mutual-1").WithSchema(openapi3.NewStringSchema())},
//				{Value: openapi3.NewCookieParameter("non-mutual-2").WithSchema(openapi3.NewStringSchema())},
//				{Value: openapi3.NewPathParameter("non-mutual-3").WithSchema(openapi3.NewStringSchema())},
//				{Value: openapi3.NewCookieParameter("non-mutual-4").WithSchema(openapi3.NewStringSchema())},
//			},
//			want1: nil,
//		},
//		{
//			name: "mutual and non mutual parameters solve conflicts",
//			args: args{
//				parameters: openapi3.Parameters{
//					{Value: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewBoolSchema())},
//					{Value: openapi3.NewQueryParameter("query-1").WithSchema(openapi3.NewInt64Schema())},
//					{Value: openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewObjectSchema().
//						WithProperty("bool", openapi3.NewBoolSchema()))},
//					{Value: openapi3.NewPathParameter("non-mutual-1").WithSchema(openapi3.NewStringSchema())},
//					{Value: openapi3.NewCookieParameter("non-mutual-2").WithSchema(openapi3.NewStringSchema())},
//				},
//				parameters2: openapi3.Parameters{
//					{Value: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewUUIDSchema())},
//					{Value: openapi3.NewQueryParameter("query-1").WithSchema(openapi3.NewUUIDSchema())},
//					{Value: openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewObjectSchema().
//						WithProperty("str", openapi3.NewDateTimeSchema()))},
//					{Value: openapi3.NewPathParameter("non-mutual-3").WithSchema(openapi3.NewStringSchema())},
//					{Value: openapi3.NewCookieParameter("non-mutual-4").WithSchema(openapi3.NewStringSchema())},
//				},
//				path: field.NewPath("parameters"),
//			},
//			want: openapi3.Parameters{
//				{Value: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewUUIDSchema())},
//				{Value: openapi3.NewQueryParameter("query-1").WithSchema(openapi3.NewUUIDSchema())},
//				{Value: openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewObjectSchema().
//					WithProperty("str", openapi3.NewDateTimeSchema()).
//					WithProperty("bool", openapi3.NewBoolSchema()))},
//				{Value: openapi3.NewPathParameter("non-mutual-1").WithSchema(openapi3.NewStringSchema())},
//				{Value: openapi3.NewCookieParameter("non-mutual-2").WithSchema(openapi3.NewStringSchema())},
//				{Value: openapi3.NewPathParameter("non-mutual-3").WithSchema(openapi3.NewStringSchema())},
//				{Value: openapi3.NewCookieParameter("non-mutual-4").WithSchema(openapi3.NewStringSchema())},
//			},
//			want1: nil,
//		},
//		{
//			name: "mutual and non mutual parameters with conflicts",
//			args: args{
//				parameters: openapi3.Parameters{
//					{Value: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewBoolSchema())},
//					{Value: openapi3.NewQueryParameter("query-1").WithSchema(openapi3.NewInt64Schema())},
//					{Value: openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewObjectSchema().
//						WithProperty("bool", openapi3.NewBoolSchema()))},
//					{Value: openapi3.NewPathParameter("non-mutual-1").WithSchema(openapi3.NewStringSchema())},
//					{Value: openapi3.NewCookieParameter("non-mutual-2").WithSchema(openapi3.NewStringSchema())},
//				},
//				parameters2: openapi3.Parameters{
//					{Value: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewInt64Schema())},
//					{Value: openapi3.NewQueryParameter("query-1").WithSchema(openapi3.NewBoolSchema())},
//					{Value: openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewObjectSchema().
//						WithProperty("str", openapi3.NewDateTimeSchema()))},
//					{Value: openapi3.NewPathParameter("non-mutual-3").WithSchema(openapi3.NewStringSchema())},
//					{Value: openapi3.NewCookieParameter("non-mutual-4").WithSchema(openapi3.NewStringSchema())},
//				},
//				path: field.NewPath("parameters"),
//			},
//			want: openapi3.Parameters{
//				{Value: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewBoolSchema())},
//				{Value: openapi3.NewQueryParameter("query-1").WithSchema(openapi3.NewInt64Schema())},
//				{Value: openapi3.NewHeaderParameter("header").WithSchema(openapi3.NewObjectSchema().
//					WithProperty("str", openapi3.NewDateTimeSchema()).
//					WithProperty("bool", openapi3.NewBoolSchema()))},
//				{Value: openapi3.NewPathParameter("non-mutual-1").WithSchema(openapi3.NewStringSchema())},
//				{Value: openapi3.NewCookieParameter("non-mutual-2").WithSchema(openapi3.NewStringSchema())},
//				{Value: openapi3.NewPathParameter("non-mutual-3").WithSchema(openapi3.NewStringSchema())},
//				{Value: openapi3.NewCookieParameter("non-mutual-4").WithSchema(openapi3.NewStringSchema())},
//			},
//			want1: []conflict{
//				{
//					path: field.NewPath("parameters").Child("X-Header-1"),
//					obj1: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewBoolSchema()),
//					obj2: openapi3.NewHeaderParameter("X-Header-1").WithSchema(openapi3.NewInt64Schema()),
//					msg: createConflictMsg(field.NewPath("parameters").Child("X-Header-1"), openapi3.TypeBoolean,
//						openapi3.TypeInteger),
//				},
//				{
//					path: field.NewPath("parameters").Child("query-1"),
//					obj1: openapi3.NewQueryParameter("query-1").WithSchema(openapi3.NewInt64Schema()),
//					obj2: openapi3.NewQueryParameter("query-1").WithSchema(openapi3.NewBoolSchema()),
//					msg: createConflictMsg(field.NewPath("parameters").Child("query-1"), openapi3.TypeInteger,
//						openapi3.TypeBoolean),
//				},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, got1 := mergeParameters(tt.args.parameters, tt.args.parameters2, tt.args.path)
//			sortParam(got)
//			sortParam(tt.want)
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("mergeParameters() got = %v, want %v", marshal(got), marshal(tt.want))
//			}
//			if !reflect.DeepEqual(got1, tt.want1) {
//				t.Errorf("mergeParameters() got1 = %v, want %v", got1, tt.want1)
//			}
//		})
//	}
//}
//
//func sortParam(got openapi3.Parameters) {
//	sort.Slice(got, func(i, j int) bool {
//		right := got[i]
//		left := got[j]
//		// Sibling parameters must have unique name + in values
//		return right.Value.Name+right.Value.In < left.Value.Name+left.Value.In
//	})
//}
//
//func Test_appendSecurityIfNeeded(t *testing.T) {
//	type args struct {
//		securityMap          openapi3.SecurityRequirement
//		mergedSecurity       openapi3.SecurityRequirements
//		ignoreSecurityKeyMap map[string]bool
//	}
//	tests := []struct {
//		name                     string
//		args                     args
//		wantMergedSecurity       openapi3.SecurityRequirements
//		wantIgnoreSecurityKeyMap map[string]bool
//	}{
//		{
//			name: "sanity",
//			args: args{
//				securityMap:          openapi3.SecurityRequirement{"key": {"val1", "val2"}},
//				mergedSecurity:       nil,
//				ignoreSecurityKeyMap: map[string]bool{},
//			},
//			wantMergedSecurity:       openapi3.SecurityRequirements{{"key": {"val1", "val2"}}},
//			wantIgnoreSecurityKeyMap: map[string]bool{"key": true},
//		},
//		{
//			name: "key should be ignored",
//			args: args{
//				securityMap:          openapi3.SecurityRequirement{"key": {"val1", "val2"}},
//				mergedSecurity:       openapi3.SecurityRequirements{{"old-key": {}}},
//				ignoreSecurityKeyMap: map[string]bool{"key": true},
//			},
//			wantMergedSecurity:       openapi3.SecurityRequirements{{"old-key": {}}},
//			wantIgnoreSecurityKeyMap: map[string]bool{"key": true},
//		},
//		{
//			name: "new key should not be ignored, old key should be ignored",
//			args: args{
//				securityMap:          openapi3.SecurityRequirement{"old-key": {}, "new key": {"val1", "val2"}},
//				mergedSecurity:       openapi3.SecurityRequirements{{"old-key": {}}},
//				ignoreSecurityKeyMap: map[string]bool{"old-key": true, "key": true},
//			},
//			wantMergedSecurity:       openapi3.SecurityRequirements{{"old-key": {}}, {"new key": {"val1", "val2"}}},
//			wantIgnoreSecurityKeyMap: map[string]bool{"old-key": true, "key": true, "new key": true},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, got1 := appendSecurityIfNeeded(tt.args.securityMap, tt.args.mergedSecurity, tt.args.ignoreSecurityKeyMap)
//			if !reflect.DeepEqual(got, tt.wantMergedSecurity) {
//				t.Errorf("appendSecurityIfNeeded() got = %v, want %v", got, tt.wantMergedSecurity)
//			}
//			if !reflect.DeepEqual(got1, tt.wantIgnoreSecurityKeyMap) {
//				t.Errorf("appendSecurityIfNeeded() got1 = %v, want %v", got1, tt.wantIgnoreSecurityKeyMap)
//			}
//		})
//	}
//}
//
//func Test_mergeOperationSecurity(t *testing.T) {
//	type args struct {
//		security  *openapi3.SecurityRequirements
//		security2 *openapi3.SecurityRequirements
//	}
//	tests := []struct {
//		name string
//		args args
//		want *openapi3.SecurityRequirements
//	}{
//		{
//			name: "no merge is needed",
//			args: args{
//				security:  &openapi3.SecurityRequirements{{"key1": {}}, {"key2": {"val1", "val2"}}},
//				security2: &openapi3.SecurityRequirements{{"key1": {}}, {"key2": {"val1", "val2"}}},
//			},
//			want: &openapi3.SecurityRequirements{{"key1": {}}, {"key2": {"val1", "val2"}}},
//		},
//		{
//			name: "full merge",
//			args: args{
//				security:  &openapi3.SecurityRequirements{{"key1": {}}},
//				security2: &openapi3.SecurityRequirements{{"key2": {"val1", "val2"}}},
//			},
//			want: &openapi3.SecurityRequirements{{"key1": {}}, {"key2": {"val1", "val2"}}},
//		},
//		{
//			name: "second list is a sub list of the first - result should be the first list",
//			args: args{
//				security:  &openapi3.SecurityRequirements{{"key1": {}}, {"key2": {"val1", "val2"}}, {"key3": {}}},
//				security2: &openapi3.SecurityRequirements{{"key2": {"val1", "val2"}}},
//			},
//			want: &openapi3.SecurityRequirements{{"key1": {}}, {"key2": {"val1", "val2"}}, {"key3": {}}},
//		},
//		{
//			name: "first list is provided as an AND - output as OR",
//			args: args{
//				security: &openapi3.SecurityRequirements{
//					{"key1": {} /*AND*/, "key2": {"val1", "val2"}},
//				},
//				security2: &openapi3.SecurityRequirements{{"key2": {"val1", "val2"}}},
//			},
//			want: &openapi3.SecurityRequirements{
//				{"key1": {}},
//				/*OR*/
//				{"key2": {"val1", "val2"}},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got := mergeOperationSecurity(tt.args.security, tt.args.security2)
//			sort.Slice(*got, func(i, j int) bool {
//				_, ok := (*got)[i]["key1"]
//				return ok
//			})
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("mergeOperationSecurity() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func Test_isEmptyRequestBody(t *testing.T) {
//	nonEmptyContent := openapi3.NewContent()
//	nonEmptyContent["test"] = openapi3.NewMediaType()
//	type args struct {
//		body *openapi3.RequestBodyRef
//	}
//	tests := []struct {
//		name string
//		args args
//		want bool
//	}{
//		{
//			name: "body == nil",
//			args: args{
//				body: nil,
//			},
//			want: true,
//		},
//		{
//			name: "body.Value == nil",
//			args: args{
//				body: &openapi3.RequestBodyRef{Value: nil},
//			},
//			want: true,
//		},
//		{
//			name: "len(body.Value.Content) == 0",
//			args: args{
//				body: &openapi3.RequestBodyRef{Value: openapi3.NewRequestBody().WithContent(nil)},
//			},
//			want: true,
//		},
//		{
//			name: "len(body.Value.Content) == 0",
//			args: args{
//				body: &openapi3.RequestBodyRef{Value: openapi3.NewRequestBody().WithContent(openapi3.Content{})},
//			},
//			want: true,
//		},
//		{
//			name: "not empty",
//			args: args{
//				body: &openapi3.RequestBodyRef{Value: openapi3.NewRequestBody().WithContent(nonEmptyContent)},
//			},
//			want: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := isEmptyRequestBody(tt.args.body); got != tt.want {
//				t.Errorf("isEmptyRequestBody() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
//
//func Test_shouldReturnIfEmptyRequestBody(t *testing.T) {
//	nonEmptyContent := openapi3.NewContent()
//	nonEmptyContent["test"] = openapi3.NewMediaType()
//	reqBody := &openapi3.RequestBodyRef{Value: openapi3.NewRequestBody().WithContent(nonEmptyContent)}
//	type args struct {
//		body  *openapi3.RequestBodyRef
//		body2 *openapi3.RequestBodyRef
//	}
//	tests := []struct {
//		name  string
//		args  args
//		want  *openapi3.RequestBodyRef
//		want1 bool
//	}{
//		{
//			name: "first body is nil",
//			args: args{
//				body:  nil,
//				body2: reqBody,
//			},
//			want:  reqBody,
//			want1: true,
//		},
//		{
//			name: "second body is nil",
//			args: args{
//				body:  reqBody,
//				body2: nil,
//			},
//			want:  reqBody,
//			want1: true,
//		},
//		{
//			name: "both bodies non nil",
//			args: args{
//				body:  reqBody,
//				body2: reqBody,
//			},
//			want:  nil,
//			want1: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, got1 := shouldReturnIfEmptyRequestBody(tt.args.body, tt.args.body2)
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("shouldReturnIfEmptyRequestBody() got = %v, want %v", got, tt.want)
//			}
//			if got1 != tt.want1 {
//				t.Errorf("shouldReturnIfEmptyRequestBody() got1 = %v, want %v", got1, tt.want1)
//			}
//		})
//	}
//}
//
//func Test_mergeRequestBody(t *testing.T) {
//	requestBody := openapi3.NewRequestBody()
//	requestBody.Content = openapi3.NewContent()
//	requestBody.Content["application/json"] = openapi3.NewMediaType().WithSchema(openapi3.NewStringSchema())
//	requestBody.Content["application/xml"] = openapi3.NewMediaType().WithSchema(openapi3.NewStringSchema())
//
//	type args struct {
//		body  *openapi3.RequestBodyRef
//		body2 *openapi3.RequestBodyRef
//		path  *field.Path
//	}
//	tests := []struct {
//		name  string
//		args  args
//		want  *openapi3.RequestBodyRef
//		want1 []conflict
//	}{
//		{
//			name: "first is nil",
//			args: args{
//				body: nil,
//				body2: &openapi3.RequestBodyRef{
//					Value: openapi3.NewRequestBody().WithJSONSchema(openapi3.NewStringSchema()),
//				},
//				path: nil,
//			},
//			want: &openapi3.RequestBodyRef{
//				Value: openapi3.NewRequestBody().WithJSONSchema(openapi3.NewStringSchema()),
//			},
//			want1: nil,
//		},
//		{
//			name: "second is nil",
//			args: args{
//				body: &openapi3.RequestBodyRef{
//					Value: openapi3.NewRequestBody().WithJSONSchema(openapi3.NewStringSchema()),
//				},
//				body2: nil,
//				path:  nil,
//			},
//			want: &openapi3.RequestBodyRef{
//				Value: openapi3.NewRequestBody().WithJSONSchema(openapi3.NewStringSchema()),
//			},
//			want1: nil,
//		},
//		{
//			name: "both are nil",
//			args: args{
//				body:  nil,
//				body2: nil,
//				path:  nil,
//			},
//			want:  nil,
//			want1: nil,
//		},
//		{
//			name: "non mutual contents",
//			args: args{
//				body: &openapi3.RequestBodyRef{
//					Value: openapi3.NewRequestBody().WithJSONSchema(openapi3.NewStringSchema()),
//				},
//				body2: &openapi3.RequestBodyRef{
//					Value: openapi3.NewRequestBody().WithSchema(openapi3.NewStringSchema(), []string{"application/xml"}),
//				},
//				path: nil,
//			},
//			want: &openapi3.RequestBodyRef{
//				Value: openapi3.NewRequestBody().
//					WithSchema(openapi3.NewStringSchema(), []string{"application/json", "application/xml"}),
//			},
//			want1: nil,
//		},
//		{
//			name: "mutual contents",
//			args: args{
//				body: &openapi3.RequestBodyRef{
//					Value: openapi3.NewRequestBody().WithJSONSchema(openapi3.NewStringSchema()),
//				},
//				body2: &openapi3.RequestBodyRef{
//					Value: openapi3.NewRequestBody().WithJSONSchema(openapi3.NewStringSchema()),
//				},
//				path: nil,
//			},
//			want: &openapi3.RequestBodyRef{
//				Value: openapi3.NewRequestBody().
//					WithJSONSchema(openapi3.NewStringSchema()),
//			},
//			want1: nil,
//		},
//		{
//			name: "mutual and non mutual contents",
//			args: args{
//				body: &openapi3.RequestBodyRef{
//					Value: openapi3.NewRequestBody().WithJSONSchema(openapi3.NewStringSchema()),
//				},
//				body2: &openapi3.RequestBodyRef{
//					Value: openapi3.NewRequestBody().WithSchema(openapi3.NewStringSchema(), []string{"application/xml", "application/json"}),
//				},
//				path: nil,
//			},
//			want: &openapi3.RequestBodyRef{
//				Value: openapi3.NewRequestBody().
//					WithSchema(openapi3.NewStringSchema(), []string{"application/xml", "application/json"}),
//			},
//			want1: nil,
//		},
//		{
//			name: "non mutual contents solve conflicts",
//			args: args{
//				body: &openapi3.RequestBodyRef{
//					Value: openapi3.NewRequestBody().
//						WithSchema(openapi3.NewStringSchema(), []string{"application/xml"}),
//				},
//				body2: &openapi3.RequestBodyRef{
//					Value: openapi3.NewRequestBody().
//						WithSchema(openapi3.NewInt64Schema(), []string{"application/xml"}),
//				},
//				path: field.NewPath("requestBody"),
//			},
//			want: &openapi3.RequestBodyRef{
//				Value: openapi3.NewRequestBody().
//					WithSchema(openapi3.NewStringSchema(), []string{"application/xml"}),
//			},
//			want1: nil,
//		},
//		{
//			name: "non mutual contents with conflicts",
//			args: args{
//				body: &openapi3.RequestBodyRef{
//					Value: openapi3.NewRequestBody().
//						WithSchema(openapi3.NewBoolSchema(), []string{"application/xml"}),
//				},
//				body2: &openapi3.RequestBodyRef{
//					Value: openapi3.NewRequestBody().
//						WithSchema(openapi3.NewInt64Schema(), []string{"application/xml"}),
//				},
//				path: field.NewPath("requestBody"),
//			},
//			want: &openapi3.RequestBodyRef{
//				Value: openapi3.NewRequestBody().
//					WithSchema(openapi3.NewBoolSchema(), []string{"application/xml"}),
//			},
//			want1: []conflict{
//				{
//					path: field.NewPath("requestBody").Child("content").Child("application/xml"),
//					obj1: openapi3.NewBoolSchema(),
//					obj2: openapi3.NewInt64Schema(),
//					msg: createConflictMsg(field.NewPath("requestBody").Child("content").Child("application/xml"),
//						openapi3.TypeBoolean, openapi3.TypeInteger),
//				},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, got1 := mergeRequestBody(tt.args.body, tt.args.body2, tt.args.path)
//			//assert.DeepEqual(t, got, tt.want, cmpopts.IgnoreUnexported(openapi3.Schema{}))
//			assertEqual(t, got, tt.want)
//			if !reflect.DeepEqual(got1, tt.want1) {
//				t.Errorf("mergeRequestBody() got1 = %v, want %v", got1, tt.want1)
//			}
//		})
//	}
//}
//
//func Test_mergeContent(t *testing.T) {
//	type args struct {
//		content  openapi3.Content
//		content2 openapi3.Content
//		path     *field.Path
//	}
//	tests := []struct {
//		name  string
//		args  args
//		want  openapi3.Content
//		want1 []conflict
//	}{
//		{
//			name: "first is nil",
//			args: args{
//				content: nil,
//				content2: openapi3.Content{
//					"json": openapi3.NewMediaType().WithSchema(openapi3.NewStringSchema()),
//				},
//				path: nil,
//			},
//			want: openapi3.Content{
//				"json": openapi3.NewMediaType().WithSchema(openapi3.NewStringSchema()),
//			},
//			want1: nil,
//		},
//		{
//			name: "second is nil",
//			args: args{
//				content: openapi3.Content{
//					"json": openapi3.NewMediaType().WithSchema(openapi3.NewStringSchema()),
//				},
//				content2: nil,
//				path:     nil,
//			},
//			want: openapi3.Content{
//				"json": openapi3.NewMediaType().WithSchema(openapi3.NewStringSchema()),
//			},
//			want1: nil,
//		},
//		{
//			name: "both are nil",
//			args: args{
//				content:  nil,
//				content2: nil,
//				path:     nil,
//			},
//			want:  openapi3.NewContent(),
//			want1: nil,
//		},
//		{
//			name: "non mutual contents",
//			args: args{
//				content: openapi3.Content{
//					"json": openapi3.NewMediaType().WithSchema(openapi3.NewStringSchema()),
//				},
//				content2: openapi3.Content{
//					"xml": openapi3.NewMediaType().WithSchema(openapi3.NewStringSchema()),
//				},
//				path: nil,
//			},
//			want: openapi3.Content{
//				"json": openapi3.NewMediaType().WithSchema(openapi3.NewStringSchema()),
//				"xml":  openapi3.NewMediaType().WithSchema(openapi3.NewStringSchema()),
//			},
//			want1: nil,
//		},
//		{
//			name: "mutual contents",
//			args: args{
//				content: openapi3.Content{
//					"json": openapi3.NewMediaType().WithSchema(openapi3.NewStringSchema()),
//				},
//				content2: openapi3.Content{
//					"json": openapi3.NewMediaType().WithSchema(openapi3.NewStringSchema()),
//				},
//				path: nil,
//			},
//			want: openapi3.Content{
//				"json": openapi3.NewMediaType().WithSchema(openapi3.NewStringSchema()),
//			},
//			want1: nil,
//		},
//		{
//			name: "mutual and non mutual contents",
//			args: args{
//				content: openapi3.Content{
//					"json": openapi3.NewMediaType().WithSchema(openapi3.NewStringSchema()),
//					"foo":  openapi3.NewMediaType().WithSchema(openapi3.NewInt64Schema()),
//				},
//				content2: openapi3.Content{
//					"json": openapi3.NewMediaType().WithSchema(openapi3.NewStringSchema()),
//					"xml":  openapi3.NewMediaType().WithSchema(openapi3.NewStringSchema()),
//				},
//				path: nil,
//			},
//			want: openapi3.Content{
//				"json": openapi3.NewMediaType().WithSchema(openapi3.NewStringSchema()),
//				"foo":  openapi3.NewMediaType().WithSchema(openapi3.NewInt64Schema()),
//				"xml":  openapi3.NewMediaType().WithSchema(openapi3.NewStringSchema()),
//			},
//			want1: nil,
//		},
//		{
//			name: "mutual contents solve conflicts",
//			args: args{
//				content: openapi3.Content{
//					"xml": openapi3.NewMediaType().WithSchema(openapi3.NewStringSchema()),
//				},
//				content2: openapi3.Content{
//					"xml": openapi3.NewMediaType().WithSchema(openapi3.NewInt64Schema()),
//				},
//				path: field.NewPath("start"),
//			},
//			want: openapi3.Content{
//				"xml": openapi3.NewMediaType().WithSchema(openapi3.NewStringSchema()),
//			},
//			want1: nil,
//		},
//		{
//			name: "mutual contents with conflicts",
//			args: args{
//				content: openapi3.Content{
//					"xml": openapi3.NewMediaType().WithSchema(openapi3.NewBoolSchema()),
//				},
//				content2: openapi3.Content{
//					"xml": openapi3.NewMediaType().WithSchema(openapi3.NewInt64Schema()),
//				},
//				path: field.NewPath("start"),
//			},
//			want: openapi3.Content{
//				"xml": openapi3.NewMediaType().WithSchema(openapi3.NewBoolSchema()),
//			},
//			want1: []conflict{
//				{
//					path: field.NewPath("start").Child("xml"),
//					obj1: openapi3.NewBoolSchema(),
//					obj2: openapi3.NewInt64Schema(),
//					msg:  createConflictMsg(field.NewPath("start").Child("xml"), openapi3.TypeBoolean, openapi3.TypeInteger),
//				},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, got1 := mergeContent(tt.args.content, tt.args.content2, tt.args.path)
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("mergeContent() got = %v, want %v", got, tt.want)
//			}
//			if !reflect.DeepEqual(got1, tt.want1) {
//				t.Errorf("mergeContent() got1 = %v, want %v", got1, tt.want1)
//			}
//		})
//	}
//}
