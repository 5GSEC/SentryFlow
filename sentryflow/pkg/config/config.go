// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package config

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	DefaultConfigFilePath             = "config/default.yaml"
	SentryFlowDefaultFilterServerPort = 8081
)

type endpoint struct {
	Url  string `json:"url"`
	Port uint16 `json:"port"`
}

type base struct {
	Name string `json:"name,omitempty"`
	// Todo: Do we really need both gRPC and http variants?
	Grpc *endpoint `json:"grpc,omitempty"`
	Http *endpoint `json:"http,omitempty"`
}

type serviceMesh struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type receivers struct {
	ServiceMeshes []*serviceMesh `json:"serviceMeshes,omitempty"`
	Others        []*base        `json:"others,omitempty"`
}

type envoyFilterConfig struct {
	Uri string `json:"uri"`
}

type server struct {
	Port uint16 `json:"port"`
}

type filters struct {
	Envoy  *envoyFilterConfig `json:"envoy,omitempty"`
	Server *server            `json:"server,omitempty"`
}

type Config struct {
	Filters   *filters   `json:"filters"`
	Receivers *receivers `json:"receivers"`
	Exporter  *base      `json:"exporter"`
}

func (c *Config) validate() error {
	if c.Filters == nil {
		return fmt.Errorf("no filter configuration provided")
	}
	if c.Filters.Envoy != nil {
		if c.Filters.Envoy.Uri == "" {
			return fmt.Errorf("no envoy filter URI provided")
		}
	}

	if c.Exporter == nil {
		return fmt.Errorf("no exporter configuration provided")
	}
	if c.Exporter.Grpc == nil {
		return fmt.Errorf("no exporter's gRPC configuration provided")
	}
	if c.Exporter.Grpc != nil && c.Exporter.Grpc.Port == 0 {
		return fmt.Errorf("no exporter's gRPC port provided")
	}
	if c.Exporter.Http != nil {
		return fmt.Errorf("http exporter is not supported")
	}

	if c.Receivers == nil {
		return fmt.Errorf("no receiver configuration provided")
	}

	for _, svcMesh := range c.Receivers.ServiceMeshes {
		if svcMesh.Name == "" {
			return fmt.Errorf("no service mesh name provided")
		}
		if svcMesh.Namespace == "" {
			return fmt.Errorf("no service mesh namespace provided")
		}
	}

	return nil
}

func New(configFilePath string, logger *zap.SugaredLogger) (*Config, error) {
	if configFilePath == "" {
		configFilePath = DefaultConfigFilePath
		logger.Warnf("Using default configfile path: %s", configFilePath)
	}

	viper.SetConfigFile(configFilePath)
	if err := viper.ReadInConfig(); err != nil {
		logger.Errorf("Failed to read config file: %v", err)
		return nil, err
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		logger.Errorf("Failed to unmarshal config file: %v", err)
		return nil, err
	}
	if config.Filters.Server == nil {
		config.Filters.Server = &server{}
	}
	if config.Filters.Server.Port == 0 {
		config.Filters.Server.Port = SentryFlowDefaultFilterServerPort
		logger.Warnf("Using default SentryFlow filter server port %d", config.Filters.Server.Port)
	}

	if err := config.validate(); err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(config)
	if err != nil {
		logger.Errorf("Failed to marshal config file: %v", err)
	}
	logger.Debugf("Config: %s", string(bytes))

	return config, nil
}
