package core

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/5GSEC/SentryFlow/protobuf"
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

func (m *Manager) startHttpServer() {
	m.Logger.Info("Starting HTTP server")
	const port = 8081
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

		if r.ProtoMajor == 2 {
			if strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
				apiEvent.Protocol = "grpc"
			} else {
				apiEvent.Protocol = "HTTP/2.0"
			}
		} else if r.ProtoMajor == 1 && r.ProtoMinor == 1 {
			apiEvent.Protocol = "HTTP/1.1"
		} else if r.ProtoMajor == 1 && r.ProtoMinor == 0 {
			apiEvent.Protocol = "HTTP/1.0"
		}
		m.ApiEvents <- apiEvent
	})
}
