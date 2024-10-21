// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package util

import (
	"context"

	"go.uber.org/zap"
)

type LoggerContextKey struct{}

const (
	ServiceMeshIstioSidecar             = "istio-sidecar"
	ServiceMeshIstioAmbient             = "istio-ambient"
	ServiceMeshKong                     = "kong"
	ServiceMeshConsul                   = "consul"
	ServiceMeshLinkerd                  = "linkerd"
	OpenTelemetry                       = "otel"
	NginxWebServer                      = "nginx-webserver"
	NginxIncorporationIngressController = "nginx-inc-ingress-controller" // https://github.com/nginxinc/kubernetes-ingress/
)

func LoggerFromCtx(ctx context.Context) *zap.SugaredLogger {
	logger, _ := ctx.Value(LoggerContextKey{}).(*zap.SugaredLogger)
	return logger
}
