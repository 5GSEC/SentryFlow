package apispec

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
)

func Test_shouldPreferType(t *testing.T) {
	type args struct {
		t1 *openapi3.Types
		t2 *openapi3.Types
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "should not prefer - bool",
			args: args{
				t1: &openapi3.Types{
					openapi3.TypeBoolean,
				},
			},
			want: false,
		},
		{
			name: "should not prefer - obj",
			args: args{
				t1: &openapi3.Types{
					openapi3.TypeObject,
				},
			},
			want: false,
		},
		{
			name: "should not prefer - array",
			args: args{
				t1: &openapi3.Types{
					openapi3.TypeArray,
				},
			},
			want: false,
		},
		{
			name: "should not prefer - number over object",
			args: args{
				t1: &openapi3.Types{
					openapi3.TypeNumber,
				},
				t2: &openapi3.Types{
					openapi3.TypeObject,
				},
			},
			want: false,
		},
		{
			name: "prefer - number over int",
			args: args{
				t1: &openapi3.Types{
					openapi3.TypeNumber,
				},
				t2: &openapi3.Types{
					openapi3.TypeInteger,
				},
			},
			want: true,
		},
		{
			name: "prefer - string over anything",
			args: args{
				t1: &openapi3.Types{
					openapi3.TypeString,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldPreferType(tt.args.t1, tt.args.t2); got != tt.want {
				t.Errorf("shouldPreferType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_conflictSolver(t *testing.T) {
	type args struct {
		t1 *openapi3.Types
		t2 *openapi3.Types
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "no conflict",
			args: args{
				t1: &openapi3.Types{openapi3.TypeNumber},
				t2: &openapi3.Types{openapi3.TypeNumber},
			},
			want: NoConflict,
		},
		{
			name: "prefer string over anything",
			args: args{
				t1: &openapi3.Types{openapi3.TypeString},
				t2: &openapi3.Types{openapi3.TypeNumber},
			},
			want: PreferType1,
		},
		{
			name: "prefer string over anything",
			args: args{
				t1: &openapi3.Types{openapi3.TypeInteger},
				t2: &openapi3.Types{openapi3.TypeString},
			},
			want: PreferType2,
		},
		{
			name: "prefer number over int",
			args: args{
				t1: &openapi3.Types{openapi3.TypeNumber},
				t2: &openapi3.Types{openapi3.TypeInteger},
			},
			want: PreferType1,
		},
		{
			name: "prefer number over int",
			args: args{
				t1: &openapi3.Types{openapi3.TypeInteger},
				t2: &openapi3.Types{openapi3.TypeNumber},
			},
			want: PreferType2,
		},
		{
			name: "conflict - bool",
			args: args{
				t1: &openapi3.Types{openapi3.TypeInteger},
				t2: &openapi3.Types{openapi3.TypeBoolean},
			},
			want: ConflictUnresolved,
		},
		{
			name: "conflict - obj",
			args: args{
				t1: &openapi3.Types{openapi3.TypeObject},
				t2: &openapi3.Types{openapi3.TypeBoolean},
			},
			want: ConflictUnresolved,
		},
		{
			name: "conflict - array",
			args: args{
				t1: &openapi3.Types{openapi3.TypeObject},
				t2: &openapi3.Types{openapi3.TypeArray},
			},
			want: ConflictUnresolved,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := conflictSolver(tt.args.t1, tt.args.t2); got != tt.want {
				t.Errorf("conflictSolver() = %v, want %v", got, tt.want)
			}
		})
	}
}
