// SPDX-License-Identifier: Apache-2.0

package collector

import (
	"io"
	"log"

	"github.com/5GSEC/SentryFlow/core"
	envoyAls "github.com/envoyproxy/go-control-plane/envoy/service/accesslog/v3"
	envoyMetrics "github.com/envoyproxy/go-control-plane/envoy/service/metrics/v3"
	"google.golang.org/grpc"
)

// EnvoyMetricsServer Structure
type EnvoyMetricsServer struct {
	envoyMetrics.UnimplementedMetricsServiceServer
	collectorInterface
}

// newEnvoyMetricsServer Function
func newEnvoyMetricsServer() *EnvoyMetricsServer {
	ret := &EnvoyMetricsServer{}
	return ret
}

// registerService Function
func (ems *EnvoyMetricsServer) registerService(server *grpc.Server) {
	envoyMetrics.RegisterMetricsServiceServer(server, ems)
}

// StreamMetrics Function
func (ems *EnvoyMetricsServer) StreamMetrics(stream envoyMetrics.MetricsService_StreamMetricsServer) error {
	event, err := stream.Recv()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		log.Printf("[Envoy] Something went on wrong when receiving event: %v", err)
		return err
	}

	err = event.ValidateAll()
	if err != nil {
		log.Printf("[Envoy] Failed to validate stream: %v", err)
	}

	// @todo parse this event entry into our format
	identifier := event.GetIdentifier()
	identifier.GetNode().GetMetadata()

	if identifier != nil {
		log.Printf("[Envoy] Received EnvoyMetric - ID: %s, %s", identifier.GetNode().GetId(), identifier.GetNode().GetCluster())
		metaData := identifier.GetNode().GetMetadata().AsMap()

		envoyMetric := core.GenerateMetricFromEnvoy(event, metaData)

		core.Lh.InsertLog(envoyMetric)
	}

	return nil
}

// EnvoyAccessLogsServer Structure
type EnvoyAccessLogsServer struct {
	envoyAls.UnimplementedAccessLogServiceServer
	collectorInterface
}

// newEnvoyAccessLogsServer Function
func newEnvoyAccessLogsServer() *EnvoyAccessLogsServer {
	ret := &EnvoyAccessLogsServer{}
	return ret
}

// registerService Function
func (eas *EnvoyAccessLogsServer) registerService(server *grpc.Server) {
	envoyAls.RegisterAccessLogServiceServer(server, eas)
}

// StreamAccessLogs Function
func (eas *EnvoyAccessLogsServer) StreamAccessLogs(stream envoyAls.AccessLogService_StreamAccessLogsServer) error {
	for {
		event, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Printf("[Envoy] Something went on wrong when receiving event: %v", err)
			return err
		}

		// Check HTTP logs
		if event.GetHttpLogs() != nil {
			for _, entry := range event.GetHttpLogs().LogEntry {
				envoyAccessLog := core.GenerateAccessLogsFromEnvoy(entry)
				core.Lh.InsertLog(envoyAccessLog)
			}
		}
	}
}
