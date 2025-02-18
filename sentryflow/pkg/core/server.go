// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package core

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"google.golang.org/protobuf/encoding/protojson"

	protobuf "github.com/5GSEC/SentryFlow/protobuf/golang"
)

func (m *Manager) startGrpcServer(port uint16) {
	m.Logger.Info("Starting gRPC server")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		m.Logger.Fatalf("Failed to listen on %v port, error: %v", port, err)
	}

	m.Logger.Infof("gRPC server listening on port %d", port)
	if err := m.GrpcServer.Serve(listener); err != nil {
		m.Logger.Fatalf("Failed to serve gRPC server on port %d, error: %v", port, err)
	}
}

func (m *Manager) startHttpServer(port uint16) {
	m.Logger.Info("Starting HTTP server")

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", m.healthzHandler)
	mux.HandleFunc("/api/v1/events", m.eventsHandler)

	m.HttpServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		ReadTimeout:       3 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	m.Logger.Infof("HTTP server listening on port %d", port)
	if err := m.HttpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		m.Logger.Fatalf("Failed to serve http server, error: %v", err)
	}
}

func (m *Manager) eventsHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if request.Body == nil {
		m.Logger.Info("Body is nil")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(request.Body)
	if err != nil {
		m.Logger.Errorf("failed to read request body, error: %v", err)
		http.Error(writer, "failed to read request body", http.StatusInternalServerError)
		return
	}

	apiEvent := &protobuf.APIEvent{}
	if err := protojson.Unmarshal(body, apiEvent); err != nil {
		m.Logger.Info("failed to unmarshal api event, error:", err)
		http.Error(writer, "failed to unmarshal request body", http.StatusBadRequest)
		return
	}

	m.ApiEvents <- apiEvent
	writer.WriteHeader(http.StatusAccepted)
	m.Logger.Debug(apiEvent)
}

func (m *Manager) healthzHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	writer.WriteHeader(http.StatusOK)
}

func (m *Manager) stopServers() {
	m.Logger.Info("Stopping servers")
	if err := m.HttpServer.Shutdown(context.Background()); err != nil {
		m.Logger.Errorf("Failed to shutdown http server, error: %v", err)
	}
	m.GrpcServer.GracefulStop()
	m.Logger.Info("Stopped servers")
}
