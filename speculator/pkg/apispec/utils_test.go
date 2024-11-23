package apispec

import (
	"testing"
)

func TestGetPathAndQuery(t *testing.T) {
	type args struct {
		fullPath string
	}
	tests := []struct {
		name      string
		args      args
		wantPath  string
		wantQuery string
	}{
		{
			name: "no query params",
			args: args{
				fullPath: "/path",
			},
			wantPath:  "/path",
			wantQuery: "",
		},
		{
			name: "path with ? in last index",
			args: args{
				fullPath: "/path?",
			},
			wantPath:  "/path?",
			wantQuery: "",
		},
		{
			name: "path with query",
			args: args{
				fullPath: "/path?query=param",
			},
			wantPath:  "/path",
			wantQuery: "query=param",
		},
		{
			name: "path with query and several ?",
			args: args{
				fullPath: "/path?query=param?stam=foo",
			},
			wantPath:  "/path",
			wantQuery: "query=param?stam=foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPath, gotQuery := GetPathAndQuery(tt.args.fullPath)
			if gotPath != tt.wantPath {
				t.Errorf("GetPathAndQuery() gotPath = %v, want %v", gotPath, tt.wantPath)
			}
			if gotQuery != tt.wantQuery {
				t.Errorf("GetPathAndQuery() gotQuery = %v, want %v", gotQuery, tt.wantQuery)
			}
		})
	}
}
