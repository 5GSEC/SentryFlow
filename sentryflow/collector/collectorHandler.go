// SPDX-License-Identifier: Apache-2.0

package collector

import (
	"fmt"
	"github.com/5gsec/SentryFlow/config"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

// == //

// ColH global reference for Collector Handler
var ColH *ColHandler

// init Function
func init() {
	ColH = NewCollectorHandler()
}

// ColHandler Structure
type ColHandler struct {
	colService net.Listener
	grpcServer *grpc.Server
	collectors []collectorInterface
}

// NewCollectorHandler Function
func NewCollectorHandler() *ColHandler {
	ch := &ColHandler{
		collectors: make([]collectorInterface, 0),
	}
	return ch
}

// == //

// StartCollector Function
func StartCollector() bool {
	// Make a string with the given collector address and port
	collectorService := fmt.Sprintf("%s:%s", config.GlobalConfig.CollectorAddr, config.GlobalConfig.CollectorPort)

	// Start listening gRPC port
	colService, err := net.Listen("tcp", collectorService)
	if err != nil {
		log.Printf("[Collector] Failed to listen at %s: %v", collectorService, err)
		return false
	}
	ColH.colService = colService

	log.Printf("[Collector] Listening Collector gRPC services (%s)", collectorService)

	// Create gRPC Service
	gRPCServer := grpc.NewServer()
	ColH.grpcServer = gRPCServer

	// initialize OpenTelemetry collector
	ColH.collectors = append(ColH.collectors, newOpenTelemetryLogsServer())

	// initialize Envoy collectors for AccessLogs and Metrics
	ColH.collectors = append(ColH.collectors, newEnvoyAccessLogsServer())
	ColH.collectors = append(ColH.collectors, newEnvoyMetricsServer())

	// register services
	for _, col := range ColH.collectors {
		col.registerService(ColH.grpcServer)
	}

	log.Print("[Collector] Initialized Collector gRPC services")

	// Serve gRPC Service
	go ColH.grpcServer.Serve(ColH.colService)

	log.Print("[Collector] Serving Collector gRPC services")

	// Start the http server
	address := fmt.Sprintf("%s:%s", config.GlobalConfig.ApiLogCollectorAddr, config.GlobalConfig.ApiLogCollectorPort)
	log.Print("[Collector] Serving Collector http service on ", address)
	go func() {
		// Create a new HTTP server
		http.HandleFunc("/api/v1/events", DataHandler)
		err = http.ListenAndServe(address, nil)
		if err != nil {
			log.Println("[Collector] Error serving Collector http service on ", err.Error())
			panic(err)
		}
	}()

	return true
}

// StopCollector Function
func StopCollector() bool {
	ColH.grpcServer.GracefulStop()

	log.Print("[Collector] Gracefully stopped Collector gRPC services")

	return true
}

// == //
