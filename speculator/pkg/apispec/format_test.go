package apispec

import (
	"testing"
)

// format taken from time/format.go.
func Test_isDateFormat(t *testing.T) {
	type args struct {
		input interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "RFC3339 should not match",
			args: args{
				input: "2021-08-23T06:52:48Z03:00",
			},
			want: false,
		},
		{
			name: "StampNano",
			args: args{
				input: "Aug 23 06:52:48.000000000",
			},
			want: true,
		},
		{
			name: "StampMicro",
			args: args{
				input: "Aug 23 06:52:48.000000",
			},
			want: true,
		},
		{
			name: "StampMilli",
			args: args{
				input: "Aug 23 06:52:48.000",
			},
			want: true,
		},
		{
			name: "Stamp",
			args: args{
				input: "Aug 23 06:52:48",
			},
			want: true,
		},
		{
			name: "RFC1123Z",
			args: args{
				input: "Mon, 23 Aug 2021 06:52:48 -0300",
			},
			want: true,
		},
		{
			name: "RFC1123",
			args: args{
				input: "Mon, 23 Aug 2021 06:52:48 GMT",
			},
			want: true,
		},
		{
			name: "RFC850",
			args: args{
				input: "Monday, 23-Aug-21 06:52:48 GMT",
			},
			want: true,
		},
		{
			name: "RFC822Z",
			args: args{
				input: "23 Aug 21 06:52 -0300",
			},
			want: true,
		},
		{
			name: "RFC822",
			args: args{
				input: "23 Aug 21 06:52 GMT",
			},
			want: true,
		},
		{
			name: "RubyDate",
			args: args{
				input: "Mon Aug 23 06:52:48 -0300 2021",
			},
			want: true,
		},
		{
			name: "UnixDate",
			args: args{
				input: "Mon Aug 23 06:52:48 GMT 2021",
			},
			want: true,
		},
		{
			name: "ANSIC",
			args: args{
				input: "Mon Aug 23 06:52:48 2021",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isDateFormat(tt.args.input); got != tt.want {
				t.Errorf("isDateFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}
