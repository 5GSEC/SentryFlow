package apispec

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func Test_extractQueryParams(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    url.Values
		wantErr bool
	}{
		{
			name: "no query params",
			args: args{
				path: "path",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "no query params with ?",
			args: args{
				path: "path?",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "with query params",
			args: args{
				path: "path?foo=bar&foo=bar2",
			},
			want:    map[string][]string{"foo": {"bar", "bar2"}},
			wantErr: false,
		},
		{
			name: "invalid query params",
			args: args{
				path: "path?foo%2",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractQueryParams(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractQueryParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractQueryParams() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addQueryParam(t *testing.T) {
	type args struct {
		operation *openapi3.Operation
		key       string
		values    []string
	}
	tests := []struct {
		name string
		args args
		want *openapi3.Operation
	}{
		{
			name: "sanity",
			args: args{
				operation: openapi3.NewOperation(),
				key:       "key",
				values:    []string{"val1"},
			},
			want: createTestOperation().WithParameter(openapi3.NewQueryParameter("key").WithSchema(openapi3.NewStringSchema())).Op,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addQueryParam(tt.args.operation, tt.args.key, tt.args.values); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addQueryParam() = %v, want %v", got, tt.want)
			}
		})
	}
}
