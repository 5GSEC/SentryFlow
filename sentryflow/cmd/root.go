// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/5GSEC/SentryFlow/sentryflow/pkg/core"
)

var (
	configFilePath string
	kubeConfig     string
	development    bool
	logger         *zap.SugaredLogger
)

func init() {
	RootCmd.PersistentFlags().StringVar(&configFilePath, "config", "", "config file path")
	RootCmd.PersistentFlags().StringVar(&kubeConfig, "kubeconfig", "", "kubeconfig file path")
	RootCmd.PersistentFlags().BoolVar(&development, "development", false, "run in development mode")
}

var RootCmd = &cobra.Command{
	Use:   "sentryflow",
	Short: "API observability",
	Long: `
SentryFlow provides real-time monitoring of API calls made to and from your system. 
`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func run() {
	initLogger(development)
	logBuildInfo()
	core.Run(configFilePath, kubeConfig, logger)
}

func initLogger(development bool) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	if development {
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	coreLogger, _ := cfg.Build()
	logger = coreLogger.Sugar()
}
