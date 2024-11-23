package util

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LevelInfo  = "info"
	LevelDebug = "debug"
)

var logger *zap.SugaredLogger

func InitLogger(debugMode bool) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	if debugMode {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	coreLogger, _ := cfg.Build()
	logger = coreLogger.Sugar()
}

func GetLogger() *zap.SugaredLogger {
	if logger == nil {
		InitLogger(false)
	}
	return logger
}
