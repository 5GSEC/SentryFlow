// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package sidecar

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"
	"text/template"

	_struct "github.com/golang/protobuf/ptypes/struct"
	"go.uber.org/zap"
	extensionsv1alpha1 "istio.io/api/extensions/v1alpha1"
	"istio.io/api/type/v1beta1"
	"istio.io/client-go/pkg/apis/extensions/v1alpha1"
	networkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/5GSEC/SentryFlow/sentryflow/pkg/config"
	"github.com/5GSEC/SentryFlow/sentryflow/pkg/util"
)

var istioRootNs = getIstioRootNamespaceFromConfig(getConfig())

func Test_createEnvoyFilter(t *testing.T) {
	cfg := getConfig()
	ctx := context.WithValue(context.Background(), util.LoggerContextKey{}, zap.S())
	fakeClient := getFakeClient()

	t.Run("when filter doesn't exist should create it", func(t *testing.T) {
		// Given
		envoyFilter := getEnvoyFilter()
		want, _ := json.Marshal(envoyFilter)
		defer func() {
			if err := fakeClient.Delete(ctx, envoyFilter); err != nil {
				t.Errorf("createEnvoyFilter() failed to delete filter = %v", err)
			}
		}()

		// When
		if err := createEnvoyFilter(ctx, cfg, fakeClient); err != nil {
			t.Errorf("createEnvoyFilter() error = %v, wantErr = nil", err)
		}

		// Then
		latestFilter := &networkingv1alpha3.EnvoyFilter{}
		_ = fakeClient.Get(ctx, client.ObjectKeyFromObject(envoyFilter), latestFilter)
		got, _ := json.Marshal(latestFilter)
		if !bytes.Equal(got, want) {
			t.Errorf("createEnvoyFilter() got = %v, want = %v", string(got), string(want))
		}
	})

	t.Run("when filter already exists should not create new one", func(t *testing.T) {
		// Given
		envoyFilter := getEnvoyFilter()

		want, _ := json.Marshal(envoyFilter)
		envoyFilter.ResourceVersion = ""

		if err := fakeClient.Create(ctx, envoyFilter); err != nil {
			t.Errorf("createEnvoyFilter() failed to create filter = %v", err)
		}
		defer func() {
			if err := fakeClient.Delete(ctx, getEnvoyFilter()); err != nil {
				t.Errorf("createEnvoyFilter() failed to delete filter = %v", err)
			}
		}()

		// When
		if err := createEnvoyFilter(ctx, cfg, fakeClient); err != nil {
			t.Errorf("createEnvoyFilter() error = %v, wantErr = nil", err)
		}

		// Then
		filter := &networkingv1alpha3.EnvoyFilter{}
		_ = fakeClient.Get(ctx, client.ObjectKeyFromObject(envoyFilter), filter)
		got, _ := json.Marshal(filter)
		if !bytes.Equal(got, want) {
			t.Errorf("createEnvoyFilter() got = %v, want = %v", string(got), string(want))
		}
	})
}

func Test_createWasmPlugin(t *testing.T) {
	cfg := getConfig()
	ctx := context.WithValue(context.Background(), util.LoggerContextKey{}, zap.S())
	fakeClient := getFakeClient()

	t.Run("when wasm plugin doesn't exist should create it", func(t *testing.T) {
		// Given
		wasmPlugin := getWasmPlugin()
		want, _ := json.Marshal(wasmPlugin)
		defer func() {
			if err := fakeClient.Delete(ctx, wasmPlugin); err != nil {
				t.Errorf("createWasmPlugin() failed to delete plugin = %v", err)
			}
		}()

		// When
		if err := createWasmPlugin(ctx, cfg, fakeClient); err != nil {
			t.Errorf("createWasmPlugin() error = %v, wantErr = nil", err)
		}

		// Then
		latestWasmPlugin := &v1alpha1.WasmPlugin{}
		_ = fakeClient.Get(ctx, client.ObjectKeyFromObject(wasmPlugin), latestWasmPlugin)
		got, _ := json.Marshal(latestWasmPlugin)
		if !bytes.Equal(got, want) {
			t.Errorf("createWasmPlugin() got = %v, want = %v", string(got), string(want))
		}
	})

	t.Run("when wasm plugin already exist should not create new one", func(t *testing.T) {
		// Given
		wasmPlugin := getWasmPlugin()

		want, _ := json.Marshal(wasmPlugin)
		wasmPlugin.ResourceVersion = ""

		if err := fakeClient.Create(ctx, wasmPlugin); err != nil {
			t.Errorf("createWasmPlugin() failed to create error = %v, wantErr = nil", err)
		}
		defer func() {
			if err := fakeClient.Delete(ctx, wasmPlugin); err != nil {
				t.Errorf("createWasmPlugin() failed to delete plugin = %v", err)
			}
		}()

		// When
		if err := createWasmPlugin(ctx, cfg, fakeClient); err != nil {
			t.Errorf("createWasmPlugin() error = %v, wantErr = nil", err)
		}

		// Then
		latestWasmPlugin := &v1alpha1.WasmPlugin{}
		_ = fakeClient.Get(ctx, client.ObjectKeyFromObject(wasmPlugin), latestWasmPlugin)
		got, _ := json.Marshal(latestWasmPlugin)
		if !bytes.Equal(got, want) {
			t.Errorf("createWasmPlugin() got = %v, want = %v", string(got), string(want))
		}
	})
}

func Test_deleteEnvoyFilter(t *testing.T) {
	ctx := context.WithValue(context.Background(), util.LoggerContextKey{}, zap.S())
	fakeClient := getFakeClient()

	t.Run("when filter exists should delete it and return no error", func(t *testing.T) {
		// Given
		envoyFilter := getEnvoyFilter()
		envoyFilter.ResourceVersion = ""
		if err := fakeClient.Create(ctx, envoyFilter); err != nil {
			t.Errorf("deleteEnvoyFilter() failed to create filter error = %v, wantErr = nil", err)
		}

		// When & Then
		if err := deleteEnvoyFilter(zap.S(), fakeClient, istioRootNs); err != nil {
			t.Errorf("deleteEnvoyFilter() error = %v, wantErr = nil", err)
		}
	})

	t.Run("when filter doesn't exist should return error", func(t *testing.T) {
		// Given
		errMessage := `envoyfilters.networking.istio.io "http-filter" not found`

		// When
		err := deleteEnvoyFilter(zap.S(), fakeClient, istioRootNs)

		// Then
		if err == nil {
			t.Errorf("deleteEnvoyFilter() error = nil, wantErr = %v", errMessage)
		}
		if err.Error() != errMessage {
			t.Errorf("deleteEnvoyFilter() errorMessage = %v, wantErrMessage = %v", err, errMessage)
		}

	})
}

func Test_deleteWasmPlugin(t *testing.T) {
	ctx := context.WithValue(context.Background(), util.LoggerContextKey{}, zap.S())
	fakeClient := getFakeClient()

	t.Run("when wasm plugin exists should delete it and return no error", func(t *testing.T) {
		// Given
		wasmPlugin := getWasmPlugin()
		wasmPlugin.ResourceVersion = ""
		if err := fakeClient.Create(ctx, wasmPlugin); err != nil {
			t.Errorf("deleteWasmPlugin() failed to create wasm plugin error = %v, wantErr = nil", err)
		}

		// When & Then
		if err := deleteWasmPlugin(zap.S(), fakeClient, istioRootNs); err != nil {
			t.Errorf("deleteWasmPlugin() error = %v, wantErr = nil", err)
		}
	})

	t.Run("when wasm plugin doesn't exist should return error", func(t *testing.T) {
		// Given
		errMessage := `wasmplugins.extensions.istio.io "http-filter" not found`

		// When
		err := deleteWasmPlugin(zap.S(), fakeClient, istioRootNs)

		// Then
		if err == nil {
			t.Errorf("deleteWasmPlugin() error = nil, wantErr = %v", errMessage)
		}
		if err.Error() != errMessage {
			t.Errorf("deleteWasmPlugin() errorMessage = %v, wantErrMessage = %v", err, errMessage)
		}

	})
}

func Test_getIstioRootNamespaceFromConfig(t *testing.T) {
	t.Run("with valid istio-sidecar receiver config should return its namespace", func(t *testing.T) {
		if got := getIstioRootNamespaceFromConfig(getConfig()); got != "istio-system" {
			t.Errorf("getIstioRootNamespaceFromConfig() got = %v, want %v", got, "istio-system")
		}
	})
}

func getConfig() *config.Config {
	configFilePath, err := filepath.Abs(filepath.Join("..", "..", "..", "..", "config", "test-configs", "default-config.yaml"))
	if err != nil {
		panic(fmt.Errorf("failed to get absolute path of config file: %v", err))
	}

	cfg, err := config.New(configFilePath, zap.S())
	if err != nil {
		panic(fmt.Errorf("failed to create config: %v", err))
	}

	return cfg
}

func getFakeClient() client.WithWatch {
	scheme := runtime.NewScheme()
	utilruntime.Must(networkingv1alpha3.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))

	return fake.
		NewClientBuilder().
		WithScheme(scheme).
		Build()
}

func getEnvoyFilter() *networkingv1alpha3.EnvoyFilter {
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
		IstioRootNs:                "istio-system",
		UpstreamAndClusterName:     UpstreamAndClusterName,
		SentryFlowFilterServerPort: 8081,
	}

	tmpl, err := template.New("envoyHttpFilter").Parse(httpFilter)
	if err != nil {
		return nil
	}

	envoyFilter := &bytes.Buffer{}
	if err := tmpl.Execute(envoyFilter, data); err != nil {
		return nil
	}

	filter := &networkingv1alpha3.EnvoyFilter{
		TypeMeta: metav1.TypeMeta{
			Kind:       "EnvoyFilter",
			APIVersion: "networking.istio.io/v1alpha3",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:            FilterName,
			Namespace:       istioRootNs,
			ResourceVersion: "1",
		},
	}
	if err := yaml.UnmarshalStrict(envoyFilter.Bytes(), filter); err != nil {
		return nil
	}
	return filter
}

func getWasmPlugin() *v1alpha1.WasmPlugin {
	cfg := getConfig()

	return &v1alpha1.WasmPlugin{
		TypeMeta: metav1.TypeMeta{
			Kind:       "WasmPlugin",
			APIVersion: "extensions.istio.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      FilterName,
			Namespace: istioRootNs,
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "sentryflow",
			},
			ResourceVersion: "1",
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
}
