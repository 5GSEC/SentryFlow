package util

import (
	"testing"
)

func TestIsApplicationJsonMediaType(t *testing.T) {
	type args struct {
		mediaType string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "application/json",
			args: args{
				mediaType: "application/json",
			},
			want: true,
		},
		{
			name: "application/hal+json",
			args: args{
				mediaType: "application/hal+json",
			},
			want: true,
		},
		{
			name: "not application json mime",
			args: args{
				mediaType: "test/html",
			},
			want: false,
		},
		{
			name: "empty mediaType",
			args: args{
				mediaType: "",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsApplicationJSONMediaType(tt.args.mediaType); got != tt.want {
				t.Errorf("IsApplicationJSONMediaType() = %v, want %v", got, tt.want)
			}
		})
	}
}
