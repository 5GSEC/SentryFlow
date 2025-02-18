// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package config

import (
	"path/filepath"
	"reflect"
	"testing"

	"go.uber.org/zap"

	"github.com/5GSEC/SentryFlow/pkg/util"
)

func TestConfig_validate(t *testing.T) {
	type fields struct {
		Filters   *Filters
		Receivers *Receivers
		Exporter  *ExporterConfig
	}
	tests := []struct {
		name               string
		fields             fields
		wantErr            bool
		expectedErrMessage string
	}{
		{
			name: "with nil filter config should return error",
			fields: fields{
				Filters: nil,
				Receivers: &Receivers{
					ServiceMeshes: []*NameAndNamespace{
						{
							Name:      "istio-sidecar",
							Namespace: "istio-system",
						},
					},
				},
				Exporter: &ExporterConfig{
					Grpc: &Server{
						Port: 11111,
					},
				},
			},
			wantErr:            true,
			expectedErrMessage: "no filter configuration provided",
		},
		{
			name: "with empty envoy URI in filter should return error",
			fields: fields{
				Filters: &Filters{
					Envoy: &EnvoyFilterConfig{
						Uri: "",
					},
				},
				Receivers: &Receivers{
					ServiceMeshes: []*NameAndNamespace{
						{
							Name:      "istio-sidecar",
							Namespace: "istio-system",
						},
					},
				},
				Exporter: &ExporterConfig{
					Grpc: &Server{
						Port: 11111,
					},
				},
			},
			wantErr:            true,
			expectedErrMessage: "no envoy filter URI provided",
		},
		{
			name: "with nil exporter config should return error",
			fields: fields{
				Filters: &Filters{
					Envoy: &EnvoyFilterConfig{
						Uri: "5gsec/http-filter:v0.1",
					},
					Server: &Server{
						Port: SentryFlowDefaultFilterServerPort,
					},
				},
				Receivers: &Receivers{
					ServiceMeshes: []*NameAndNamespace{
						{
							Name:      "istio-sidecar",
							Namespace: "istio-system",
						},
					},
				},
				Exporter: nil,
			},
			wantErr:            true,
			expectedErrMessage: "no exporter configuration provided",
		},
		{
			name: "with nil exporter gRPC config should return error",
			fields: fields{
				Filters: &Filters{
					Envoy: &EnvoyFilterConfig{
						Uri: "5gsec/http-filter:v0.1",
					},
					Server: &Server{
						Port: SentryFlowDefaultFilterServerPort,
					},
				},
				Receivers: &Receivers{
					ServiceMeshes: []*NameAndNamespace{
						{
							Name:      "istio-sidecar",
							Namespace: "istio-system",
						},
					},
				},
				Exporter: &ExporterConfig{
					Grpc: nil,
				},
			},
			wantErr:            true,
			expectedErrMessage: "no exporter's gRPC configuration provided",
		},
		{
			name: "without exporter's gRPC port config should return error",
			fields: fields{
				Filters: &Filters{
					Envoy: &EnvoyFilterConfig{
						Uri: "5gsec/http-filter:v0.1",
					},
					Server: &Server{
						Port: SentryFlowDefaultFilterServerPort,
					},
				},
				Receivers: &Receivers{
					ServiceMeshes: []*NameAndNamespace{
						{
							Name:      "istio-sidecar",
							Namespace: "istio-system",
						},
					},
				},
				Exporter: &ExporterConfig{
					Grpc: &Server{},
				},
			},
			wantErr:            true,
			expectedErrMessage: "no exporter's gRPC port provided",
		},
		{
			name: "with nil receiver config should return error",
			fields: fields{
				Filters: &Filters{
					Envoy: &EnvoyFilterConfig{
						Uri: "5gsec/http-filter:v0.1",
					},
					Server: &Server{
						Port: SentryFlowDefaultFilterServerPort,
					},
				},
				Receivers: nil,
				Exporter: &ExporterConfig{
					Grpc: &Server{
						Port: 11111,
					},
				},
			},
			wantErr:            true,
			expectedErrMessage: "no receiver configuration provided",
		},
		{
			name: "with istio-sidecar svcmesh and without envoy URI should return error",
			fields: fields{
				Filters: &Filters{
					Envoy: nil,
					Server: &Server{
						Port: SentryFlowDefaultFilterServerPort,
					},
				},
				Receivers: &Receivers{
					ServiceMeshes: []*NameAndNamespace{
						{
							Name:      util.ServiceMeshIstioSidecar,
							Namespace: "istio-system",
						},
					},
				},
				Exporter: &ExporterConfig{
					Grpc: &Server{
						Port: 11111,
					},
				},
			},
			wantErr:            true,
			expectedErrMessage: "no envoy filter configuration provided for istio sidecar servicemesh",
		},
		{
			name: "with empty service mesh name receiver should return error",
			fields: fields{
				Filters: &Filters{
					Envoy: &EnvoyFilterConfig{
						Uri: "5gsec/http-filter:v0.1",
					},
					Server: &Server{
						Port: SentryFlowDefaultFilterServerPort,
					},
				},
				Receivers: &Receivers{
					ServiceMeshes: []*NameAndNamespace{
						{
							Namespace: "istio-system",
						},
					},
				},
				Exporter: &ExporterConfig{
					Grpc: &Server{
						Port: 11111,
					},
				},
			},
			wantErr:            true,
			expectedErrMessage: "no service mesh name provided",
		},
		{
			name: "with empty service mesh namespace receiver should return error",
			fields: fields{
				Filters: &Filters{
					Envoy: &EnvoyFilterConfig{
						Uri: "5gsec/http-filter:v0.1",
					},
					Server: &Server{
						Port: SentryFlowDefaultFilterServerPort,
					},
				},
				Receivers: &Receivers{
					ServiceMeshes: []*NameAndNamespace{
						{
							Name: "istio-sidecar",
						},
					},
				},
				Exporter: &ExporterConfig{
					Grpc: &Server{
						Port: 11111,
					},
				},
			},
			wantErr:            true,
			expectedErrMessage: "no service mesh namespace provided",
		},
		{
			name: "with valid config should not return error",
			fields: fields{
				Filters: &Filters{
					Envoy: &EnvoyFilterConfig{
						Uri: "5gsec/http-filter:v0.1",
					},
					Server: &Server{
						Port: SentryFlowDefaultFilterServerPort,
					},
				},
				Receivers: &Receivers{
					ServiceMeshes: []*NameAndNamespace{
						{
							Name:      "istio-sidecar",
							Namespace: "istio-system",
						},
					},
				},
				Exporter: &ExporterConfig{
					Grpc: &Server{
						Port: 11111,
					},
				},
			},
			wantErr:            false,
			expectedErrMessage: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				Filters:   tt.fields.Filters,
				Receivers: tt.fields.Receivers,
				Exporter:  tt.fields.Exporter,
			}

			err := c.validate()
			if tt.wantErr && err == nil {
				t.Errorf("validate() expected error but got nil")
			} else if !tt.wantErr && err != nil {
				t.Errorf("validate() expected no error but got error = %v", err)
			} else if tt.wantErr && err != nil && tt.expectedErrMessage != err.Error() {
				t.Errorf("validate() expected error message to be %v but got %v", tt.expectedErrMessage, err.Error())
			}
		})
	}
}

func TestNew(t *testing.T) {
	logger := zap.S()

	type args struct {
		configFilePath string
		logger         *zap.SugaredLogger
	}
	tests := []struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}{
		{
			name: "with valid configFilePath should return config",
			args: args{
				configFilePath: filepath.Join(".", "test-configs", "default-config.yaml"),
				logger:         logger,
			},
			want: &Config{
				Filters: &Filters{
					Envoy: &EnvoyFilterConfig{
						Uri: "anuragrajawat/httpfilter:v0.1",
					},
					Server: &Server{
						Port: 8081,
					},
				},
				Receivers: &Receivers{
					ServiceMeshes: []*NameAndNamespace{
						{
							Name:      "istio-sidecar",
							Namespace: "istio-system",
						},
					},
				},
				Exporter: &ExporterConfig{
					Grpc: &Server{
						Port: 8080,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "with invalid configFilePath should return error",
			args: args{
				configFilePath: filepath.Join(".", "path-doesnt-exist", "invalid-config.yaml"),
				logger:         logger,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "with nil filter server configFilePath should return config with default filter server",
			args: args{
				configFilePath: filepath.Join(".", "test-configs", "without-filter-server.yaml"),
				logger:         logger,
			},
			want: &Config{
				Filters: &Filters{
					Envoy: &EnvoyFilterConfig{
						Uri: "anuragrajawat/httpfilter:v0.1",
					},
					Server: &Server{
						Port: 9999,
					},
				},
				Receivers: &Receivers{
					ServiceMeshes: []*NameAndNamespace{
						{
							Name:      "istio-sidecar",
							Namespace: "istio-system",
						},
					},
				},
				Exporter: &ExporterConfig{
					Grpc: &Server{
						Port: 8080,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "without filter server port config should return config with default port",
			args: args{
				configFilePath: filepath.Join(".", "test-configs", "without-filter-server.yaml"),
				logger:         logger,
			},
			want: &Config{
				Filters: &Filters{
					Envoy: &EnvoyFilterConfig{
						Uri: "anuragrajawat/httpfilter:v0.1",
					},
					Server: &Server{
						Port: 9999,
					},
				},
				Receivers: &Receivers{
					ServiceMeshes: []*NameAndNamespace{
						{
							Name:      "istio-sidecar",
							Namespace: "istio-system",
						},
					},
				},
				Exporter: &ExporterConfig{
					Grpc: &Server{
						Port: 8080,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "with invalid config should return error",
			args: args{
				configFilePath: filepath.Join(".", "test-configs", "invalid-config.yaml"),
				logger:         logger,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.configFilePath, tt.args.logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}
