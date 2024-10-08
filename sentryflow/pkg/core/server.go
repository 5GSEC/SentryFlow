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
	m.HttpServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: 3 * time.Second,
	}
	m.registerRoutes()

	m.Logger.Infof("HTTP server listening on port %d", port)
	if err := m.HttpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		m.Logger.Fatalf("Failed to serve http server, error: %v", err)
	}
}

func (m *Manager) stopServers() {
	m.Logger.Info("Stopping servers")
	if err := m.HttpServer.Shutdown(context.Background()); err != nil {
		m.Logger.Errorf("Failed to shutdown http server, error: %v", err)
	}
	m.GrpcServer.GracefulStop()
	m.Logger.Info("Stopped servers")
}

func (m *Manager) registerRoutes() {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	// Register an endpoint to receive API events from EnvoyFilter
	http.HandleFunc("/api/v1/events", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if r.Body == nil {
			m.Logger.Info("Body is nil")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			m.Logger.Errorf("failed to read request body, error: %v", err)
			http.Error(w, "failed to read request body", http.StatusInternalServerError)
			return
		}

		apiEvent := &protobuf.APIEvent{}
		if err := protojson.Unmarshal(body, apiEvent); err != nil {
			m.Logger.Info("failed to unmarshal api event, error:", err)
			http.Error(w, "failed to parse request body", http.StatusBadRequest)
			return
		}

		m.ApiEvents <- apiEvent
	})
}
