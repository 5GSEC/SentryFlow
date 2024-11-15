// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package config

import (
	"path/filepath"
	"reflect"
	"testing"

	"go.uber.org/zap"
)

func TestConfig_validate(t *testing.T) {
	type fields struct {
		Filters   *filters
		Receivers *receivers
		Exporter  *base
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
				Receivers: &receivers{
					ServiceMeshes: []*serviceMesh{
						{
							Name:      "istio-sidecar",
							Namespace: "istio-system",
						},
					},
				},
				Exporter: &base{
					Name: "default gRPC exporter",
					Grpc: &endpoint{
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
				Filters: &filters{
					Envoy: &envoyFilterConfig{
						Uri: "",
					},
				},
				Receivers: &receivers{
					ServiceMeshes: []*serviceMesh{
						{
							Name:      "istio-sidecar",
							Namespace: "istio-system",
						},
					},
				},
				Exporter: &base{
					Name: "default gRPC exporter",
					Grpc: &endpoint{
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
				Filters: &filters{
					Envoy: &envoyFilterConfig{
						Uri: "5gsec/http-filter:v0.1",
					},
					Server: &server{
						Port: SentryFlowDefaultFilterServerPort,
					},
				},
				Receivers: &receivers{
					ServiceMeshes: []*serviceMesh{
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
				Filters: &filters{
					Envoy: &envoyFilterConfig{
						Uri: "5gsec/http-filter:v0.1",
					},
					Server: &server{
						Port: SentryFlowDefaultFilterServerPort,
					},
				},
				Receivers: &receivers{
					ServiceMeshes: []*serviceMesh{
						{
							Name:      "istio-sidecar",
							Namespace: "istio-system",
						},
					},
				},
				Exporter: &base{
					Name: "default gRPC exporter",
					Grpc: nil,
				},
			},
			wantErr:            true,
			expectedErrMessage: "no exporter's gRPC configuration provided",
		},
		{
			name: "without exporter's gRPC port config should return error",
			fields: fields{
				Filters: &filters{
					Envoy: &envoyFilterConfig{
						Uri: "5gsec/http-filter:v0.1",
					},
					Server: &server{
						Port: SentryFlowDefaultFilterServerPort,
					},
				},
				Receivers: &receivers{
					ServiceMeshes: []*serviceMesh{
						{
							Name:      "istio-sidecar",
							Namespace: "istio-system",
						},
					},
				},
				Exporter: &base{
					Name: "gRPC exporter without port",
					Grpc: &endpoint{},
				},
			},
			wantErr:            true,
			expectedErrMessage: "no exporter's gRPC port provided",
		},
		{
			name: "with HTTP exporter should return error",
			fields: fields{
				Filters: &filters{
					Envoy: &envoyFilterConfig{
						Uri: "5gsec/http-filter:v0.1",
					},
					Server: &server{
						Port: SentryFlowDefaultFilterServerPort,
					},
				},
				Receivers: &receivers{
					ServiceMeshes: []*serviceMesh{
						{
							Name:      "istio-sidecar",
							Namespace: "istio-system",
						},
					},
				},
				Exporter: &base{
					Name: "unsupported HTTP exporter",
					Grpc: &endpoint{
						Port: 11111,
					},
					Http: &endpoint{
						Port: 65432,
					},
				},
			},
			wantErr:            true,
			expectedErrMessage: "http exporter is not supported",
		},
		{
			name: "with nil receiver config should return error",
			fields: fields{
				Filters: &filters{
					Envoy: &envoyFilterConfig{
						Uri: "5gsec/http-filter:v0.1",
					},
					Server: &server{
						Port: SentryFlowDefaultFilterServerPort,
					},
				},
				Receivers: nil,
				Exporter: &base{
					Name: "default gRPC exporter",
					Grpc: &endpoint{
						Port: 11111,
					},
				},
			},
			wantErr:            true,
			expectedErrMessage: "no receiver configuration provided",
		},
		{
			name: "with empty service mesh name receiver should return error",
			fields: fields{
				Filters: &filters{
					Envoy: &envoyFilterConfig{
						Uri: "5gsec/http-filter:v0.1",
					},
					Server: &server{
						Port: SentryFlowDefaultFilterServerPort,
					},
				},
				Receivers: &receivers{
					ServiceMeshes: []*serviceMesh{
						{
							Namespace: "istio-system",
						},
					},
				},
				Exporter: &base{
					Name: "default gRPC exporter",
					Grpc: &endpoint{
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
				Filters: &filters{
					Envoy: &envoyFilterConfig{
						Uri: "5gsec/http-filter:v0.1",
					},
					Server: &server{
						Port: SentryFlowDefaultFilterServerPort,
					},
				},
				Receivers: &receivers{
					ServiceMeshes: []*serviceMesh{
						{
							Name: "istio-sidecar",
						},
					},
				},
				Exporter: &base{
					Name: "default gRPC exporter",
					Grpc: &endpoint{
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
				Filters: &filters{
					Envoy: &envoyFilterConfig{
						Uri: "5gsec/http-filter:v0.1",
					},
					Server: &server{
						Port: SentryFlowDefaultFilterServerPort,
					},
				},
				Receivers: &receivers{
					ServiceMeshes: []*serviceMesh{
						{
							Name:      "istio-sidecar",
							Namespace: "istio-system",
						},
					},
				},
				Exporter: &base{
					Name: "default gRPC exporter",
					Grpc: &endpoint{
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
				Filters: &filters{
					Envoy: &envoyFilterConfig{
						Uri: "anuragrajawat/httpfilter:v0.1",
					},
					Server: &server{
						Port: 8081,
					},
				},
				Receivers: &receivers{
					ServiceMeshes: []*serviceMesh{
						{
							Name:      "istio-sidecar",
							Namespace: "istio-system",
						},
					},
				},
				Exporter: &base{
					Grpc: &endpoint{
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
				Filters: &filters{
					Envoy: &envoyFilterConfig{
						Uri: "anuragrajawat/httpfilter:v0.1",
					},
					Server: &server{
						Port: 8081,
					},
				},
				Receivers: &receivers{
					ServiceMeshes: []*serviceMesh{
						{
							Name:      "istio-sidecar",
							Namespace: "istio-system",
						},
					},
				},
				Exporter: &base{
					Grpc: &endpoint{
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
				Filters: &filters{
					Envoy: &envoyFilterConfig{
						Uri: "anuragrajawat/httpfilter:v0.1",
					},
					Server: &server{
						Port: 8081,
					},
				},
				Receivers: &receivers{
					ServiceMeshes: []*serviceMesh{
						{
							Name:      "istio-sidecar",
							Namespace: "istio-system",
						},
					},
				},
				Exporter: &base{
					Grpc: &endpoint{
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
