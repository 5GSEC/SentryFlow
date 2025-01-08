// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package apievent

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/5GSEC/SentryFlow/sfctl/pkg/client"
	"github.com/5GSEC/SentryFlow/sfctl/pkg/util"
)

var (
	prettyPrint         bool
	sentryflowPort      string
	sentryflowNamespace string
	k8sClientset        kubernetes.Interface
	kubeConfigFilePath  string
	logger              = util.GetLogger()
)

func init() {
	EventCmd.PersistentFlags().BoolVar(&prettyPrint, "pretty", false, "pretty print API Events in JSON format")
	EventCmd.PersistentFlags().StringVar(&sentryflowPort, "port", "8888", "port to connect to SentryFlow")
	EventCmd.Flags().StringVar(&sentryflowNamespace, "namespace", "sentryflow", "namespace to connect to SentryFlow")
	EventCmd.AddCommand(filterCmd)
}

var EventCmd = &cobra.Command{
	Use:   "event",
	Short: "Print the captured API Events",
	Long: `Print captured API events from SentryFlow.

This prints API events captured by SentryFlow. By default, events are printed in standard JSON format.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return printEvents(cmd.Flags())
	},
	Example: `# Print API events in standard JSON format
sfctl event

# Pretty print API events in JSON format
sfctl event --pretty`,
	SilenceUsage: true,
}

func printEvents(flags *pflag.FlagSet) error {
	kubeCtx, err := flags.GetString("context")
	if err != nil {
		return err
	}

	kubeConfigFilePath, err = flags.GetString("kubeconfig")
	if err != nil {
		logger.Errorf("failed to get kubeconfig file path: %v", err)
		return err
	}

	k8sClientset, err = client.NewClientset(kubeConfigFilePath, &kubeCtx)
	if err != nil {
		logger.Errorf("failed to create Kubernetes clientset: %v", err)
		return err
	}

	logger.Debug("starting events stream")
	return startEventsStreaming(ctrl.SetupSignalHandler(), kubeConfigFilePath, k8sClientset)
}
