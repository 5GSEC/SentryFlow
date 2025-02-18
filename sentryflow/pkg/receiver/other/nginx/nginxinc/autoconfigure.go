// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package nginxinc

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/5GSEC/SentryFlow/pkg/config"
	"github.com/5GSEC/SentryFlow/pkg/util"
)

const volumeName = "sentryflow-nginx-inc"

func deployResources(ctx context.Context, k8sClient client.Client, cfg *config.Config) error {
	ingressDeployNs := getIngressControllerDeploymentNamespace(cfg)
	if err := patchIngressDeploy(ctx, k8sClient, cfg.Filters.NginxIngress.DeploymentName, ingressDeployNs); err != nil {
		return err
	}
	if err := patchIngressConfigMap(ctx, k8sClient, cfg, ingressDeployNs); err != nil {
		return err
	}
	return nil
}

func patchIngressConfigMap(ctx context.Context, k8sClient client.Client, cfg *config.Config, namespace string) error {
	nginxConfigMap := &corev1.ConfigMap{}
	if err := k8sClient.Get(ctx, types.NamespacedName{Name: cfg.Filters.NginxIngress.ConfigMapName, Namespace: namespace}, nginxConfigMap); err != nil {
		return err
	}

	const data = `
http-snippets:
js_path "/etc/nginx/njs/";
subrequest_output_buffer_size 8k;
js_shared_dict_zone zone=apievents:1M timeout=300s evict;
js_import main from sentryflow.js;

location-snippets:
js_body_filter main.requestHandler buffer_type=buffer;
mirror      /mirror_request;
mirror_request_body on;

server-snippets:

location /mirror_request {
  internal;
  js_content main.dispatchHttpCall;
}

location /sentryflow {
  internal;
  # Update SentryFlow URL with path to ingest access logs if required.
  proxy_pass http://{{ .SentryFlowSvcName }}.{{ .SentryFlowSvcNamespace }}:{{ .SentryFlowFilterServerPort }}/api/v1/events;
  proxy_method      POST;
  proxy_set_header accept "application/json";
  proxy_set_header Content-Type "application/json";
}
`
	sentryflowSvcName, sentryflowSvcNamespace, err := sentryFlowSvcNameAndNs(ctx, k8sClient)
	if err != nil {
		return err
	}

	values := struct {
		SentryFlowSvcName          string
		SentryFlowSvcNamespace     string
		SentryFlowFilterServerPort uint16
	}{
		SentryFlowSvcName:          sentryflowSvcName,
		SentryFlowSvcNamespace:     sentryflowSvcNamespace,
		SentryFlowFilterServerPort: cfg.Filters.Server.Port,
	}

	tmpl, err := template.New("nginx.tmpl").Parse(data)
	if err != nil {
		return err
	}

	cmData := &bytes.Buffer{}
	if err := tmpl.Execute(cmData, values); err != nil {
		return err
	}
	cmDataStr := cmData.String()

	parts := strings.SplitN(cmDataStr, "http-snippets:", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid template format: missing http-snippets section")
	}
	httpSnippetsPart := parts[1]

	parts = strings.SplitN(httpSnippetsPart, "location-snippets:", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid template format: missing location-snippets section")
	}
	httpSnippets := strings.TrimSpace(parts[0])
	locationSnippetsPart := parts[1]

	parts = strings.SplitN(locationSnippetsPart, "server-snippets:", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid template format: missing server-snippets section")
	}

	locationSnippets := strings.TrimSpace(parts[0])
	serverSnippets := strings.TrimSpace(parts[1])

	if nginxConfigMap.Data == nil {
		nginxConfigMap.Data = map[string]string{}
	}
	nginxConfigMap.Data["http-snippets"] = httpSnippets
	nginxConfigMap.Data["location-snippets"] = locationSnippets
	nginxConfigMap.Data["server-snippets"] = serverSnippets

	return k8sClient.Update(ctx, nginxConfigMap)
}

func sentryFlowSvcNameAndNs(ctx context.Context, k8sClient client.Client) (string, string, error) {
	svcs := &corev1.ServiceList{}
	if err := k8sClient.List(ctx, svcs, &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(map[string]string{
			"app.kubernetes.io/name": "sentryflow",
		}),
	}); err != nil {
		return "", "", err
	}
	if len(svcs.Items) == 0 {
		return "", "", fmt.Errorf("sentryFlow svc was not found")
	}
	return svcs.Items[0].Name, svcs.Items[0].Namespace, nil
}

func patchIngressDeploy(ctx context.Context, k8sClient client.Client, deploymentName, namespace string) error {
	nginxDeploy := &appsv1.Deployment{}
	if err := k8sClient.Get(ctx, types.NamespacedName{Name: deploymentName, Namespace: namespace}, nginxDeploy); err != nil {
		return err
	}

	if !containsVolumeAndVolumeMount(nginxDeploy.Spec.Template.Spec) {
		nginxDeploy.Spec.Template.Spec.Volumes = append(nginxDeploy.Spec.Template.Spec.Volumes, corev1.Volume{
			Name: volumeName,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: volumeName,
					},
				},
			},
		})

		nginxDeploy.Spec.Template.Spec.Containers[0].VolumeMounts = append(nginxDeploy.Spec.Template.Spec.Containers[0].VolumeMounts,
			corev1.VolumeMount{
				Name:      volumeName,
				MountPath: "/etc/nginx/njs/sentryflow.js",
				SubPath:   "sentryflow.js",
			},
		)

		if err := k8sClient.Update(ctx, nginxDeploy); err != nil {
			return err
		}
	}

	return nil
}

func containsVolumeAndVolumeMount(spec corev1.PodSpec) bool {
	volumeFound, volumeMountFound := false, false

	for _, volume := range spec.Volumes {
		if volume.Name == volumeName {
			volumeFound = true
			break
		}
	}

	for _, container := range spec.Containers {
		for _, volumeMount := range container.VolumeMounts {
			if volumeMount.MountPath == "/etc/nginx/njs/sentryflow.js" && volumeMount.SubPath == "sentryflow.js" {
				volumeMountFound = true
				break
			}
		}
	}

	return volumeFound && volumeMountFound
}

func autoConfigure(cfg *config.Config) bool {
	for _, other := range cfg.Receivers.Others {
		if other.Name == util.NginxIncorporationIngressController && other.AutoConfigure {
			return true
		}
	}
	return false
}

func doCleanup(logger *zap.SugaredLogger, cfg *config.Config, k8sClient client.Client) {
	logger.Info("shutting down nginx-incorporation ingress controller receiver")
	ctx := context.Background()

	if autoConfigure(cfg) {
		logger.Warn("nginx inc ingress deployment will restart")
		ingressDeployNs := getIngressControllerDeploymentNamespace(cfg)

		var err error
		if err = removePatchFromIngressConfigMap(ctx, k8sClient, cfg, ingressDeployNs); err != nil {
			logger.Error("failed to remove patch from ingress-controller deployment", err)
			// Do not return, always remove patches, even if errors occur.
		}
		if err = removePatchFromIngressDeploy(ctx, k8sClient, cfg.Filters.NginxIngress.DeploymentName, ingressDeployNs); err != nil {
			logger.Error("failed to remove patch from ingress-controller deployment", err)
			// Do not return, always remove patches, even if errors occur.
		}

		if err == nil {
			logger.Info("successfully removed patches from ingress-controller deployment and configmap")
		}
	}

	logger.Info("stopped nginx-incorporation ingress controller receiver")
}

func removePatchFromIngressDeploy(ctx context.Context, k8sClient client.Client, deploymentName, namespace string) error {
	nginxDeploy := &appsv1.Deployment{}
	if err := k8sClient.Get(ctx, types.NamespacedName{Name: deploymentName, Namespace: namespace}, nginxDeploy); err != nil {
		return err
	}
	if !containsVolumeAndVolumeMount(nginxDeploy.Spec.Template.Spec) {
		return nil
	}

	var newVolumes []corev1.Volume
	for _, volume := range nginxDeploy.Spec.Template.Spec.Volumes {
		if volume.Name != volumeName {
			newVolumes = append(newVolumes, volume)
		}
	}
	nginxDeploy.Spec.Template.Spec.Volumes = newVolumes

	var newVolumeMounts []corev1.VolumeMount
	for _, volumeMount := range nginxDeploy.Spec.Template.Spec.Containers[0].VolumeMounts {
		if volumeMount.Name != volumeName {
			newVolumeMounts = append(newVolumeMounts, volumeMount)
		}
	}
	nginxDeploy.Spec.Template.Spec.Containers[0].VolumeMounts = newVolumeMounts

	return k8sClient.Update(ctx, nginxDeploy)
}

func removePatchFromIngressConfigMap(ctx context.Context, k8sClient client.Client, cfg *config.Config, namespace string) error {
	nginxConfigMap := &corev1.ConfigMap{}
	if err := k8sClient.Get(ctx, types.NamespacedName{Name: cfg.Filters.NginxIngress.ConfigMapName, Namespace: namespace}, nginxConfigMap); err != nil {
		return err
	}

	var needToUpdate bool
	if _, exists := nginxConfigMap.Data["http-snippets"]; exists {
		delete(nginxConfigMap.Data, "http-snippets")
		needToUpdate = true
	}
	if _, exists := nginxConfigMap.Data["location-snippets"]; exists {
		delete(nginxConfigMap.Data, "location-snippets")
		needToUpdate = true
	}
	if _, exists := nginxConfigMap.Data["server-snippets"]; exists {
		delete(nginxConfigMap.Data, "server-snippets")
		needToUpdate = true
	}

	if needToUpdate {
		return k8sClient.Update(ctx, nginxConfigMap)
	}

	return nil
}
