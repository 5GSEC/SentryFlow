# Istio Sidecar Service Mesh

## Description

This guide provides a step-by-step process to integrate SentryFlow with Istio, aimed at enhancing API observability. It
includes detailed commands for each step along with their explanations.

SentryFlow makes use of following to provide visibility into API calls:

- [Envoy Wasm Filter](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/wasm_filter)
- [Istio Wasm Plugin](https://istio.io/latest/docs/reference/config/proxy_extensions/wasm-plugin/)
- [Istio EnvoyFilter](https://istio.io/latest/docs/reference/config/networking/envoy-filter/)

## Prerequisites

- Deploy Istio service mesh. Follow [this](https://istio.io/latest/docs/setup/install/) to deploy it if you've not
  deployed.
- Enable the envoy proxy injection by labeling the namespace in which you'll deploy your workloads:
  ```shell
  kubectl label ns <namespace_name> istio-injection=enabled
  ```

## How to

To Observe API calls of your workloads running on top of Istio Service Mesh in Kubernetes environment, follow the below
steps:

1. Download SentryFlow manifest file

  ```shell
  curl -sO https://raw.githubusercontent.com/5GSEC/SentryFlow/refs/heads/main/deployments/sentryflow.yaml
  ```

2. Update the `.receivers` configuration in `sentryflow` [configmap](../../../../deployments/sentryflow.yaml) as
   follows:

  ```yaml
  filters:
    server:
      port: 8081

    # Envoy filter is required for `istio-sidecar` service-mesh receiver.
    # Leave it as it is unless you want to use your filter.
    envoy:
      uri: 5gsec/sentryflow-httpfilter:v0.1

  receivers:
    serviceMeshes:
      - name: istio-sidecar # SentryFlow makes use of `name` to configure receivers. DON'T CHANGE IT.
        namespace: istio-system # Kubernetes namespace in which you've deployed Istio.
    ...
  ```

3. Apply the updated manifest file:

```shell
kubectl apply -f sentryflow.yaml
```

3. Trigger API calls to generate traffic.

4. Use SentryFlow [log client](../../../../client) to see the API Events.
