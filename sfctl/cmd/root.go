// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/5GSEC/SentryFlow/sfctl/pkg/apievent"
	"github.com/5GSEC/SentryFlow/sfctl/pkg/version"
)

var (
	kubeConfig string
	context    string
)

// Todo: Improve the description

func init() {
	rootCmd.PersistentFlags().StringVar(&kubeConfig, "kubeconfig", "", "path to the kubeconfig file")
	rootCmd.PersistentFlags().StringVar(&context, "context", "", "name of the kubeconfig context to use")
	rootCmd.AddCommand(version.VersionCmd, apievent.EventCmd)
}

var rootCmd = &cobra.Command{
	Use:   "sfctl",
	Short: "Manage SentryFlow API events.",
	Long:  `SentryFlow command-line utility for managing captured API events.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
