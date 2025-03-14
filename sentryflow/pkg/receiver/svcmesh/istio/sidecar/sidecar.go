// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package sidecar

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"text/template"

	_struct "github.com/golang/protobuf/ptypes/struct"
	"go.uber.org/zap"
	extensionsv1alpha1 "istio.io/api/extensions/v1alpha1"
	"istio.io/api/type/v1beta1"
	"istio.io/client-go/pkg/apis/extensions/v1alpha1"
	networkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/5GSEC/SentryFlow/sentryflow/pkg/config"
	"github.com/5GSEC/SentryFlow/sentryflow/pkg/util"
)

const (
	FilterName             = "http-filter"
	UpstreamAndClusterName = "sentryflow"
	ApiPath                = "/api/v1/events"
)

type envoyFilterData struct {
	FilterName                 string
	IstioRootNs                string
	UpstreamAndClusterName     string
	SentryFlowFilterServerPort uint16
}

// StartMonitoring begins monitoring API calls within the Istio (sidecar based)
// service mesh deployed in a Kubernetes cluster. It achieves this by creating a
// custom EnvoyFilter resource in Kubernetes.
func StartMonitoring(ctx context.Context, cfg *config.Config, k8sClient client.Client, lock *sync.Mutex) {
	logger := util.LoggerFromCtx(ctx).Named("istio-sidecar")
	logger.Info("Starting istio sidecar mesh monitoring")

	lock.Lock()
	if err := createResources(ctx, cfg, k8sClient); err != nil {
		logger.Error(err)
		lock.Unlock()
		return
	}
	logger.Info("Started istio sidecar mesh monitoring")
	lock.Unlock()

	<-ctx.Done()
	logger.Info("Shutting down istio sidecar mesh monitoring")

	lock.Lock()
	doCleanup(logger, k8sClient, getIstioRootNamespaceFromConfig(cfg))
	lock.Unlock()

	logger.Info("Stopped istio sidecar mesh monitoring")
}

func createResources(ctx context.Context, cfg *config.Config, k8sClient client.Client) error {
	if err := createEnvoyFilter(ctx, cfg, k8sClient); err != nil {
		return fmt.Errorf("failed to create EnvoyFilter. Stopping istio sidecar mesh monitoring, error: %v", err)
	}

	if err := createWasmPlugin(ctx, cfg, k8sClient); err != nil {
		return fmt.Errorf("failed to create WasmPlugin. Stopping istio sidecar monitoring, error: %v", err)
	}

	return nil
}

func doCleanup(logger *zap.SugaredLogger, k8sClient client.Client, istioRootNs string) {
	if err := deleteEnvoyFilter(logger, k8sClient, istioRootNs); err != nil {
		logger.Errorf("failed to delete EnvoyFilter, error: %v", err)
	}
	if err := deleteWasmPlugin(logger, k8sClient, istioRootNs); err != nil {
		logger.Errorf("failed to delete WasmPlugin, error: %v", err)
	}
}

func createWasmPlugin(ctx context.Context, cfg *config.Config, k8sClient client.Client) error {
	logger := util.LoggerFromCtx(ctx)

	wasmPlugin := &v1alpha1.WasmPlugin{
		TypeMeta: metav1.TypeMeta{
			Kind:       "WasmPlugin",
			APIVersion: "extensions.istio.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      FilterName,
			Namespace: getIstioRootNamespaceFromConfig(cfg),
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "sentryflow",
			},
		},
		Spec: extensionsv1alpha1.WasmPlugin{
			Url: cfg.Filters.Envoy.Uri,
			PluginConfig: &_struct.Struct{
				Fields: map[string]*_struct.Value{
					"upstream_name": {
						Kind: &_struct.Value_StringValue{
							StringValue: UpstreamAndClusterName,
						},
					},
					"authority": {
						Kind: &_struct.Value_StringValue{
							StringValue: UpstreamAndClusterName,
						},
					},
					"api_path": {
						Kind: &_struct.Value_StringValue{
							StringValue: ApiPath,
						},
					},
				},
			},
			PluginName:   FilterName,
			FailStrategy: extensionsv1alpha1.FailStrategy_FAIL_OPEN,
			Match: []*extensionsv1alpha1.WasmPlugin_TrafficSelector{
				{
					Mode: v1beta1.WorkloadMode_CLIENT,
				},
			},
			Type: extensionsv1alpha1.PluginType_HTTP,
		},
	}

	existingWasmPlugin := &v1alpha1.WasmPlugin{}
	if err := k8sClient.Get(ctx, client.ObjectKeyFromObject(wasmPlugin), existingWasmPlugin); err != nil {
		if errors.IsNotFound(err) {
			if err := k8sClient.Create(ctx, wasmPlugin); err != nil {
				return err
			}
			logger.Infow("Created WasmPlugin", "name", wasmPlugin.Name, "namespace", wasmPlugin.Namespace)
			return nil
		}
		return err
	}
	logger.Infow("Found existing WasmPlugin", "name", wasmPlugin.Name, "namespace", wasmPlugin.Namespace)
	return nil
}

func deleteWasmPlugin(logger *zap.SugaredLogger, k8sClient client.Client, istioRootNs string) error {
	existingWasmPlugin := &v1alpha1.WasmPlugin{}
	if err := k8sClient.Get(context.Background(), types.NamespacedName{Name: FilterName, Namespace: istioRootNs}, existingWasmPlugin); err != nil {
		return err
	}

	if err := k8sClient.Delete(context.Background(), existingWasmPlugin); err != nil {
		return err
	}
	logger.Infow("Deleted WasmPlugin", "name", FilterName, "namespace", istioRootNs)
	return nil
}

func createEnvoyFilter(ctx context.Context, cfg *config.Config, k8sClient client.Client) error {
	logger := util.LoggerFromCtx(ctx)

	// Istio feature stages for reference to keep trace when they plan to move their
	// alpha APIs to beta then stable.
	// https://istio.io/latest/docs/releases/feature-stages/#extensibility
	// https://istio.io/latest/docs/releases/feature-stages/#traffic-management [Enabling custom filters in Envoy]
	const httpFilter = `
apiVersion: networking.istio.io/v1alpha3
kind: EnvoyFilter
metadata:
  name: {{ .FilterName }}
  # Deploy the filter to whatever istio considers its "root" namespace so that we
  # don't have to create the ConfigMap(s) containing the WASM filter binary,
  # and the associated annotations/configuration for the Istio sidecar(s).
  # https://istio.io/latest/docs/reference/config/istio.mesh.v1alpha1/#MeshConfig:~:text=No-,rootNamespace,-string
  namespace: {{ .IstioRootNs }}
  labels:
    app.kubernetes.io/managed-by: sentryflow
spec:
  configPatches:
    - applyTo: CLUSTER
      match:
        context: ANY
      patch:
        operation: ADD
        value:
          name: {{ .UpstreamAndClusterName }}
          type: LOGICAL_DNS
          connect_timeout: 1s
          lb_policy: ROUND_ROBIN
          load_assignment:
            cluster_name: {{ .UpstreamAndClusterName }}
            endpoints:
              - lb_endpoints:
                  - endpoint:
                      address:
                        socket_address:
                          protocol: TCP
                          address: {{ .UpstreamAndClusterName }}.{{ .UpstreamAndClusterName }}
                          port_value: {{ .SentryFlowFilterServerPort }}
`

	data := envoyFilterData{
		FilterName:                 FilterName,
		IstioRootNs:                getIstioRootNamespaceFromConfig(cfg),
		UpstreamAndClusterName:     UpstreamAndClusterName,
		SentryFlowFilterServerPort: cfg.Filters.Server.Port,
	}

	tmpl, err := template.New("envoyHttpFilter").Parse(httpFilter)
	if err != nil {
		logger.Errorf("Failed to parse EnvoyFilter template: %v", err)
		return err
	}

	envoyFilter := &bytes.Buffer{}
	if err := tmpl.Execute(envoyFilter, data); err != nil {
		logger.Errorf("Failed to execute EnvoyFilter template: %v", err)
		return err
	}

	filterToCreate := &networkingv1alpha3.EnvoyFilter{
		TypeMeta: metav1.TypeMeta{
			Kind:       "EnvoyFilter",
			APIVersion: "networking.istio.io/v1alpha3",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      FilterName,
			Namespace: getIstioRootNamespaceFromConfig(cfg),
		},
	}
	if err := yaml.UnmarshalStrict(envoyFilter.Bytes(), filterToCreate); err != nil {
		logger.Errorf("Failed to unmarshal EnvoyFilter: %v", err)
		return err
	}

	existingFilter := &networkingv1alpha3.EnvoyFilter{}
	if err := k8sClient.Get(ctx, client.ObjectKeyFromObject(filterToCreate), existingFilter); err != nil {
		if errors.IsNotFound(err) {
			if err := k8sClient.Create(ctx, filterToCreate); err != nil {
				return err
			}
			logger.Infow("Created EnvoyFilter", "name", filterToCreate.Name, "namespace", filterToCreate.Namespace)
			return nil
		}
		return err
	}
	logger.Infow("Found existing EnvoyFilter", "name", filterToCreate.Name, "namespace", filterToCreate.Namespace)

	return nil
}

func deleteEnvoyFilter(logger *zap.SugaredLogger, k8sClient client.Client, istioRootNs string) error {
	existingFilter := &networkingv1alpha3.EnvoyFilter{}
	if err := k8sClient.Get(context.Background(), types.NamespacedName{Name: FilterName, Namespace: istioRootNs}, existingFilter); err != nil {
		return err
	}

	if err := k8sClient.Delete(context.Background(), existingFilter); err != nil {
		return err
	}
	logger.Infow("Deleted EnvoyFilter", "name", FilterName, "namespace", "istio-system")

	return nil
}

func getIstioRootNamespaceFromConfig(cfg *config.Config) string {
	for _, svcMesh := range cfg.Receivers.ServiceMeshes {
		switch svcMesh.Name {
		case util.ServiceMeshIstioSidecar:
			return svcMesh.Namespace
		}
	}
	// The `namespace` field is always non-empty due to validation during config
	// initialization.
	return ""
}
