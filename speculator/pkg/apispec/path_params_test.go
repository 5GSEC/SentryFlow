package apispec

import (
	"reflect"
	"sort"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func Test_createParameterizedPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no suspect params",
			args: args{
				path: "/api/user/hello",
			},
			want: "/api/user/hello",
		},
		{
			name: "1 suspect param",
			args: args{
				path: "/api/123/hello",
			},
			want: "/api/{param1}/hello",
		},
		{
			name: "2 suspect param",
			args: args{
				path: "/api/123/hello/234",
			},
			want: "/api/{param1}/hello/{param2}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createParameterizedPath(tt.args.path); got != tt.want {
				t.Errorf("createParameterizedPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isSuspectPathParam(t *testing.T) {
	type args struct {
		pathPart string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "number",
			args: args{
				pathPart: "1234",
			},
			want: true,
		},
		{
			name: "big number",
			args: args{
				pathPart: "123456789001234567890023456789",
			},
			want: true,
		},
		{
			name: "uuid",
			args: args{
				pathPart: "3d9f2779-264f-4930-9196-e60c8a3610d2",
			},
			want: true,
		},
		{
			name: "mixed type - numbers are more than 20%",
			args: args{
				pathPart: "abcdefghij123",
			},
			want: true,
		},
		{
			name: "mixed type - numbers are less than 20%",
			args: args{
				pathPart: "abcdefghijk12",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSuspectPathParam(tt.args.pathPart); got != tt.want {
				t.Errorf("isSuspectPathParam() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_countDigitsInString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "4",
			args: args{
				s: "abcdefg1234hijk",
			},
			want: 4,
		},
		{
			name: "0",
			args: args{
				s: "abcdefghijk",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countDigitsInString(tt.args.s); got != tt.want {
				t.Errorf("countDigitsInString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getOnlyIndexedPartFromPaths(t *testing.T) {
	type args struct {
		paths map[string]bool
		i     int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "2 numbers",
			args: args{
				paths: map[string]bool{
					"/api/1/foo": true,
					"/api/2/foo": true,
				},
				i: 1,
			},
			want: []string{"1", "2"},
		},
		{
			name: "number and string",
			args: args{
				paths: map[string]bool{
					"/api/1/foo": true,
					"/api/foo/2": true,
				},
				i: 1,
			},
			want: []string{"1", "foo"},
		},
		{
			name: "get first part",
			args: args{
				paths: map[string]bool{
					"/api/1/foo": true,
					"/api/2/foo": true,
				},
				i: 0,
			},
			want: []string{"api", "api"},
		},
		{
			name: "get last part",
			args: args{
				paths: map[string]bool{
					"/api/1/foo": true,
					"/api/2/foo": true,
				},
				i: 2,
			},
			want: []string{"foo", "foo"},
		},
		{
			name: "index is bigger than paths len",
			args: args{
				paths: map[string]bool{
					"/api/1/foo": true,
					"/api/2/foo": true,
				},
				i: 3,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getOnlyIndexedPartFromPaths(tt.args.paths, tt.args.i)
			sort.Slice(got, func(i, j int) bool {
				return got[i] < got[j]
			})
			sort.Slice(tt.want, func(i, j int) bool {
				return tt.want[i] < tt.want[j]
			})
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getOnlyIndexedPartFromPaths() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getParamTypeAndFormat(t *testing.T) {
	type args struct {
		paramsList []string
	}
	tests := []struct {
		name string
		args args
		want *openapi3.Schema
	}{
		{
			name: "mixed",
			args: args{
				paramsList: []string{"str", "1234", "77e1c83b-7bb0-437b-bc50-a7a58e5660ac"},
			},
			want: openapi3.NewStringSchema(),
		},
		{
			name: "uuid",
			args: args{
				paramsList: []string{"77e1c83b-7bb0-437b-bc50-a7a58e5660a3", "77e1c83b-7bb0-437b-bc50-a7a58e5660a8", "77e1c83b-7bb0-437b-bc50-a7a58e5660ac"},
			},
			want: openapi3.NewUUIDSchema(),
		},
		{
			name: "number",
			args: args{
				paramsList: []string{"7776", "78", "123"},
			},
			want: openapi3.NewInt64Schema(),
		},
		{
			name: "string",
			args: args{
				paramsList: []string{"strone", "strtwo", "strthree"},
			},
			want: openapi3.NewStringSchema(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schema := getParamSchema(tt.args.paramsList)
			assertEqual(t, schema, tt.want)
		})
	}
}

func Test_createPathParam(t *testing.T) {
	type args struct {
		name   string
		schema *openapi3.Schema
	}
	tests := []struct {
		name string
		args args
		want *PathParam
	}{
		{
			name: "create",
			args: args{
				name:   "param1",
				schema: openapi3.NewUUIDSchema(),
			},
			want: &PathParam{
				Parameter: openapi3.NewPathParameter("param1").WithSchema(openapi3.NewUUIDSchema()),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createPathParam(tt.args.name, tt.args.schema); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createPathParam() = %v, want %v", got, tt.want)
			}
		})
	}
}
