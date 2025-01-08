# Getting Started

This guide provides a step-by-step process for deploying SentryFlow in a Kubernetes environment, aimed at enhancing API
observability. It includes detailed commands for each step along with their explanations.

> **Note**: SentryFlow is currently in the early stages of development. Please be aware that the information provided
> here may become outdated or change without notice.

## 1. Prerequisites

- A Kubernetes cluster running version 1.28 or later.
- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) version 1.28 or later.

## 2. Deploying SentryFlow

Configure SentryFlow receiver by following [this](receivers.md). Then deploy updated SentryFlow manifest by following
`kubectl` command:

```shell
kubectl apply -f sentryflow.yaml
```

This will create a namespace named `sentryflow` and will deploy the necessary Kubernetes resources.

Then, check if SentryFlow is up and running by:

```shell
$ kubectl -n sentryflow get pods
NAME                         READY   STATUS    RESTARTS   AGE
sentryflow-cff887bbd-rljm7   1/1     Running   0          73s
```

## 3. Viewing Captured API Access Events Clients

SentryFlow has now been deployed in the cluster. In addition, SentryFlow exports API access events through `gRPC`.

You can use `sfctl` the SentryFlow client to view or filter captured API access events

```shell
$ sfctl event
{"level":"INFO","timestamp":"2025-01-08T18:15:31.720+0530","caller":"apievent/common.go:165","msg":"starting API Events streaming"}
{"level":"INFO","timestamp":"2025-01-08T18:15:31.771+0530","caller":"apievent/common.go:171","msg":"started API Events streaming"}
# API Access Events
{"metadata":{"context_id":9,"timestamp":1736340391,"istio_version":"1.24.1","mesh_id":"cluster.local","node_name":"kind-control-plane"},"source":{"name":"server-c7669846-w5v8m","namespace":"default","ip":"10.244.0.8","port":57754},"destination":{"namespace":"sentryflow","ip":"10.96.79.211","port":9999},"request":{"headers":{":authority":"sentryflow.sentryflow:9999",":method":"HEAD",":path":"/",":scheme":"http","accept":"*/*","user-agent":"curl/7.88.1","x-forwarded-proto":"http","x-request-id":"9ff1f0fb-adca-4cbb-bfbb-7927d5aa02ae"}},"response":{"headers":{":status":"404","content-length":"19","content-type":"text/plain; charset=utf-8","date":"Wed, 08 Jan 2025 12:46:31 GMT","x-content-type-options":"nosniff"}},"protocol":"HTTP/1.1"}
...
```

### Filter API Events based on some Response Status Code

```shell
$ sfctl event filter --status "4xx"
{"level":"INFO","timestamp":"2025-01-08T18:20:37.096+0530","caller":"apievent/common.go:165","msg":"starting API Events streaming"}
{"level":"INFO","timestamp":"2025-01-08T18:20:37.151+0530","caller":"apievent/common.go:171","msg":"started API Events streaming"}
# API Access Events
{"metadata":{"context_id":10,"timestamp":1736340639,"istio_version":"1.24.1","mesh_id":"cluster.local","node_name":"kind-control-plane"},"source":{"name":"server-c7669846-w5v8m","namespace":"default","ip":"10.244.0.8","port":37154},"destination":{"namespace":"sentryflow","ip":"10.96.79.211","port":9999},"request":{"headers":{":authority":"sentryflow.sentryflow:9999",":method":"HEAD",":path":"/",":scheme":"http","accept":"*/*","user-agent":"curl/7.88.1","x-forwarded-proto":"http","x-request-id":"e20a1002-09d1-4f3f-936e-ce688652ea4d"}},"response":{"headers":{":status":"404","content-length":"19","content-type":"text/plain; charset=utf-8","date":"Wed, 08 Jan 2025 12:50:39 GMT","x-content-type-options":"nosniff"}},"protocol":"HTTP/1.1"}
```

For more info check [this](../sfctl/README.md).
