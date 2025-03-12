// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package client

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// NewClientset returns a new Kubernetes Clientset.
func NewClientset(kubeConfig string, kubeCtx *string) (kubernetes.Interface, error) {
	config, err := GetConfig(kubeConfig, kubeCtx)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %v", err)
	}
	return kubernetes.NewForConfig(config)
}

func GetConfig(kubeConfig string, kubeCtx *string) (*rest.Config, error) {
	if kubeConfig == "" {
		kubeConfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}

	restClienGetter := genericclioptions.ConfigFlags{
		KubeConfig: &kubeConfig,
		Context:    kubeCtx,
	}
	rawKubeConfigLoader := restClienGetter.ToRawKubeConfigLoader()

	config, err := rawKubeConfigLoader.ClientConfig()
	if err != nil {
		return nil, err
	}

	return config, nil
}
