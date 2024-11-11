// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package util

import (
	"context"
	"reflect"
	"testing"

	"go.uber.org/zap"
)

func TestLoggerFromCtx(t *testing.T) {
	logger := zap.S()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want *zap.SugaredLogger
	}{
		{
			name: "with logger in context should return logger",
			args: args{
				ctx: context.WithValue(context.Background(), LoggerContextKey{}, logger),
			},
			want: logger,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LoggerFromCtx(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoggerFromCtx() = %v, want %v", got, tt.want)
			}
		})
	}
}
