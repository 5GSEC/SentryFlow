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

- Add SentryFlow repo

```shell
helm repo add 5gsec https://5gsec.github.io/charts
helm repo update 5gsec
```

- Update `values.yaml` file as follows.

```shell
helm show values 5gsec/sentryflow > values.yaml
```

```yaml
filters:
  server:
  # Existing snippets

  # Envoy filter is required for `istio-sidecar` service-mesh receiver.
  # Leave it as it is unless you want to use your filter.
  envoy:
    uri: 5gsec/sentryflow-httpfilter:latest

receivers:
  serviceMeshes:
    - name: istio-sidecar # SentryFlow makes use of `name` to configure receivers. DON'T CHANGE IT.
      namespace: istio-system # Kubernetes namespace in which you've deployed Istio.
  # Existing snippets
```

- Deploy SentryFlow

```shell
helm install --values values.yaml sentryflow 5gsec/sentryflow -n sentryflow --create-namespace 
```

- Trigger API calls to generate traffic.
