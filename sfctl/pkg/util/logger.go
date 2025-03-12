// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package util

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func InitLogger() {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	coreLogger, _ := cfg.Build()
	logger = coreLogger.Sugar()
}

func GetLogger() *zap.SugaredLogger {
	if logger == nil {
		InitLogger()
	}
	return logger
}
