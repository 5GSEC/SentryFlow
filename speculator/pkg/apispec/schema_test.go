package apispec

import (
	"reflect"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func Test_splitByStyle(t *testing.T) {
	type args struct {
		data  string
		style string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "empty data",
			args: args{
				data:  "",
				style: "",
			},
			want: nil,
		},
		{
			name: "unsupported serialization style",
			args: args{
				data:  "",
				style: "Unsupported",
			},
			want: nil,
		},
		{
			name: "SerializationForm",
			args: args{
				data:  "1, 2, 3",
				style: openapi3.SerializationForm,
			},
			want: []string{"1", "2", "3"},
		},
		{
			name: "SerializationSimple",
			args: args{
				data:  "1, 2, 3",
				style: openapi3.SerializationSimple,
			},
			want: []string{"1", "2", "3"},
		},
		{
			name: "SerializationSpaceDelimited",
			args: args{
				data:  "1 2  3",
				style: openapi3.SerializationSpaceDelimited,
			},
			want: []string{"1", "2", "3"},
		},
		{
			name: "SerializationPipeDelimited",
			args: args{
				data:  "1|2|3",
				style: openapi3.SerializationPipeDelimited,
			},
			want: []string{"1", "2", "3"},
		},
		{
			name: "SerializationPipeDelimited with empty space in the middle",
			args: args{
				data:  "1| |3",
				style: openapi3.SerializationPipeDelimited,
			},
			want: []string{"1", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitByStyle(tt.args.data, tt.args.style); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitByStyle() = %v, want %v", got, tt.want)
			}
		})
	}
}
