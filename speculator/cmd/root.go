// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package cmd

import (
	"github.com/spf13/cobra"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/5gsec/sentryflow/speculator/pkg/core"
	"github.com/5gsec/sentryflow/speculator/pkg/util"
)

var (
	configFilePath string
	//kubeConfig     string
	debugMode bool
)

func init() {
	RootCmd.PersistentFlags().StringVar(&configFilePath, "config", "", "config file path")
	//RootCmd.PersistentFlags().StringVar(&kubeConfig, "kubeconfig", "", "kubeconfig file path")
	RootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "run in debug mode")
}

var RootCmd = &cobra.Command{
	Use: "speculator",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func run() {
	util.InitLogger(debugMode)
	logBuildInfo(util.GetLogger())
	ctx := ctrl.SetupSignalHandler()
	core.Run(ctx, configFilePath)
}
