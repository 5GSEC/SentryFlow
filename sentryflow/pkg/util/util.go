package util

import (
	"context"

	"go.uber.org/zap"
)

type LoggerContextKey struct{}

const (
	ServiceMeshIstioSidecar = "istio-sidecar"
	ServiceMeshIstioAmbient = "istio-ambient"
	ServiceMeshKong         = "kong"
	ServiceMeshConsul       = "consul"
	ServiceMeshLinkerd      = "linkerd"
	OpenTelemetry           = "otel"
)

func LoggerFromCtx(ctx context.Context) *zap.SugaredLogger {
	logger, _ := ctx.Value(LoggerContextKey{}).(*zap.SugaredLogger)
	return logger
}
