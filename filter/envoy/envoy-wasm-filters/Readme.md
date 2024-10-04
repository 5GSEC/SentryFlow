# Envoy Wasm HTTP Filter

HTTP filter to observe RESTful and gRPC API calls made to/from a k8s workload.

## Sample API Event:

```json
{
  "metadata": {
    "context_id": 3,
    "timestamp": 1726211548,
    "istio_version": "1.23.0",
    "mesh_id": "cluster.local",
    "node_name": "kind-control-plane"
  },
  "request": {
    "headers": {
      ":scheme": "http",
      ":method": "GET",
      "x-envoy-decorator-operation": "filterserver.sentryflow.svc.cluster.local:80/*",
      ":authority": "filterserver.sentryflow",
      "user-agent": "Wget",
      "x-forwarded-proto": "http",
      "x-request-id": "6b2e87df-257c-931e-a996-5517b8155b4a"
    },
    "body": ""
  },
  "response": {
    "headers": {
      "date": "Fri, 13 Sep 2024 07:12:28 GMT",
      "content-type": "application/json; charset=utf-8",
      "content-length": "63",
      ":status": "404"
    },
    "body": "{\"message\":\"The specified route / not found\",\"status\":\"failed\"}"
  },
  "source": {
    "name": "httpd-c6d6cb94b-v259g",
    "namespace": "default",
    "ip": "10.244.0.27",
    "port": 54636
  },
  "destination": {
    "name": "",
    "namespace": "sentryflow",
    "ip": "10.96.158.214",
    "port": 80
  },
  "protocol": "HTTP/1.1"
}
```

# Getting Started

## Install development tools

You'll need these tools for a smooth development experience:

- [Make](https://www.gnu.org/software/make/#download),
- [Rust](https://www.rust-lang.org/tools/install) toolchain,
- An IDE ([RustRover](https://www.jetbrains.com/rust/) / [VS Code](https://code.visualstudio.com/download)),
- Container tools ([Docker](https://www.docker.com/) / [Podman](https://podman.io/)),
- [Kubernetes cluster](https://kubernetes.io/docs/setup/) running version 1.26 or later,
- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) version 1.26 or later.

## In Envoy alone

This example can be run with docker compose and has a matching [envoy configuration](envoy.yaml) file.

- Build the plugin
  ```shell
  make
  ```

- Start the envoy container
  ```shell
  docker compose up
  ```

- See the Raw API Events in `server` cluster configured in [envoy configuration](envoy.yaml).

# In Kubernetes

- [Install Istio](https://istio.io/latest/docs/setup/install/)
- Build the plugin
  ```shell
  make
  ```

- Build and push plugin's container image
  ```shell
  make imagex push
  ```

- Update the value of `filters.envoy.uri` with the latest image in
  SentryFlow's [configMap](https://github.com/5GSEC/SentryFlow/blob/main/deployments/sentryflow.yaml#L68)

- Deploy SentryFlow
  ```shell
  kubectl apply -f https://raw.githubusercontent.com/5GSEC/SentryFlow/refs/heads/main/deployments/sentryflow.yaml
  ```

- Enable the envoy proxy injection by labeling the namespace in which you'll deploy workload:
  ```shell
  kubectl label ns <namespace_name> istio-injection=enabled
  ```
- Deploy some workload and generate traffic by calling some APIs. For e.g., you can use
  Google's [microservices-demo](https://github.com/GoogleCloudPlatform/microservices-demo).

- Use SentryFlow's client to see the API Events. 