// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package nginxinc

import (
	"context"
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/5GSEC/SentryFlow/sentryflow/pkg/config"
	"github.com/5GSEC/SentryFlow/sentryflow/pkg/util"
)

func Start(ctx context.Context, cfg *config.Config, k8sClient client.Client) {
	logger := util.LoggerFromCtx(ctx)
	defer doCleanup(logger, cfg, k8sClient)

	logger.Info("Starting nginx-incorporation ingress controller receiver")

	if autoConfigure(cfg) {
		logger.Warn("nginx inc ingress deployment will restart")
		if err := deployResources(ctx, k8sClient, cfg); err != nil {
			logger.Errorf("failed to configure resources for nginx inc ingress receiver: %v", err)
			return
		}
		logger.Info("successfully configured nginx-incorporation ingress controller")
	}

	if err := validateResources(ctx, cfg, k8sClient); err != nil {
		logger.Errorf("%v. Stopped nginx-incorporation ingress controller receiver", err)
		return
	}
	logger.Info("Started nginx-incorporation ingress controller receiver")

	<-ctx.Done()
}

func validateResources(ctx context.Context, cfg *config.Config, k8sClient client.Client) error {
	ingressDeployNamespace := getIngressControllerDeploymentNamespace(cfg)
	sentryFlowNjsCm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cfg.Filters.NginxIngress.SentryFlowNjsConfigMapName,
			Namespace: ingressDeployNamespace,
		},
	}

	if err := k8sClient.Get(ctx, client.ObjectKeyFromObject(sentryFlowNjsCm), sentryFlowNjsCm); err != nil {
		return fmt.Errorf("failed to get sentryflow configmap: %w", err)
	}

	if err := validateIngressDeployAndConfigMap(ctx, cfg, k8sClient, ingressDeployNamespace); err != nil {
		return err
	}

	return nil
}

func validateIngressDeployAndConfigMap(ctx context.Context, cfg *config.Config, k8sClient client.Client, ingressNamespace string) error {
	ingressDeploy := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cfg.Filters.NginxIngress.DeploymentName,
			Namespace: ingressNamespace,
		},
	}
	if err := k8sClient.Get(ctx, client.ObjectKeyFromObject(ingressDeploy), ingressDeploy); err != nil {
		return fmt.Errorf("failed to get nginx-incorporation ingress controller deployment: %w", err)
	}

	// Just check ingress controller deployment volume mount because if volume
	// itself doesn't exist, the container will not start.
	volumeMountFound := false
	for _, container := range ingressDeploy.Spec.Template.Spec.Containers {
		for _, volumeMount := range container.VolumeMounts {
			// Volume-mount name could be different so only check mount-path.
			if volumeMount.MountPath == "/etc/nginx/njs/sentryflow.js" {
				volumeMountFound = true
			}
		}
	}
	if !volumeMountFound {
		return fmt.Errorf("sentryflow-njs volume-mount not found")
	}

	ingressCm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cfg.Filters.NginxIngress.ConfigMapName,
			Namespace: ingressNamespace,
		},
	}
	if err := k8sClient.Get(ctx, client.ObjectKeyFromObject(ingressCm), ingressCm); err != nil {
		return fmt.Errorf("failed to get nginx-incorporation ingress controller configmap: %w", err)
	}

	httpSnippets, exists := ingressCm.Data["http-snippets"]
	if !exists {
		return fmt.Errorf("sentryflow http-snippets not found in nginx-incorporation ingress configmap")
	}
	expectedHttpSnippets := strings.TrimSpace(`js_path "/etc/nginx/njs/";
subrequest_output_buffer_size 8k;
js_shared_dict_zone zone=apievents:1M timeout=300s evict;
js_import main from sentryflow.js;
`)
	if !strings.Contains(httpSnippets, expectedHttpSnippets) {
		return fmt.Errorf("sentryflow http-snippets were not properly configured in nginx-incorporation ingress configmap.\nGOT\n%v, \nEXPECTED\n%v", httpSnippets, expectedHttpSnippets)
	}

	locationSnippets, exists := ingressCm.Data["location-snippets"]
	if !exists {
		return fmt.Errorf("sentryflow location-snippets not found in nginx-incorporation ingress configmap")
	}
	expectedLocationSnippets := strings.TrimSpace(`js_body_filter main.requestHandler buffer_type=buffer;
mirror      /mirror_request;
mirror_request_body on;
`)
	if !strings.Contains(locationSnippets, expectedLocationSnippets) {
		return fmt.Errorf("sentryflow location-snippets were not properly configured in nginx-incorporation ingress configmap.\nGOT\n%v, \nEXPECTED\n%v", locationSnippets, expectedLocationSnippets)
	}

	serverSnippets, exists := ingressCm.Data["server-snippets"]
	if !exists {
		return fmt.Errorf("sentryflow server-snippets not found in nginx-incorporation ingress configmap")
	}
	expectedServerSnippets := strings.TrimSpace(`location /mirror_request {
  internal;
  js_content main.dispatchHttpCall;
}
location /sentryflow {
  internal;
  proxy_method      POST;
  proxy_set_header accept "application/json";
  proxy_set_header Content-Type "application/json";
}
`)
	// The server snippet might have different SentryFlow URL in `proxy_pass`
	// directive. To avoid potential conflicts, check without that directive.
	if !strings.ContainsAny(serverSnippets, expectedServerSnippets) {
		return fmt.Errorf("sentryflow server-snippets were not properly configured in nginx-incorporation ingress configmap.\nGOT\n%v, \nEXPECTED\n%v", serverSnippets, expectedServerSnippets)
	}

	return nil
}

func getIngressControllerDeploymentNamespace(cfg *config.Config) string {
	for _, other := range cfg.Receivers.Others {
		switch other.Name {
		case util.NginxIncorporationIngressController:
			return other.Namespace
		}
	}
	return ""
}
