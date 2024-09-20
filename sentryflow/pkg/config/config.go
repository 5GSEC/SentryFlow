// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package config

import (
	"encoding/json"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	DefaultConfigFilePath = "config/default.yaml"
)

type Endpoint struct {
	Url  string `json:"url"`
	Port uint16 `json:"port"`
}

type Base struct {
	Name string `json:"name,omitempty"`
	// Todo: Do we really need both gRPC and http variants?
	Grpc *Endpoint `json:"grpc,omitempty"`
	Http *Endpoint `json:"http,omitempty"`
}

type ServiceMesh struct {
	Name   string `json:"name"`
	Enable bool   `json:"enable"`
}

type Receivers struct {
	ServiceMeshes []*ServiceMesh `json:"serviceMeshes,omitempty"`
	Others        []*Base        `json:"others,omitempty"`
}

type Config struct {
	Receivers *Receivers `json:"receivers"`
	Exporter  *Base      `json:"exporter"`
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

	bytes, err := json.Marshal(config)
	if err != nil {
		logger.Errorf("Failed to marshal config file: %v", err)
	}
	logger.Debugf("Config: %s", string(bytes))

	return config, nil
}
