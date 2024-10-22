# Getting Started

This guide provides a step-by-step process for deploying SentryFlow in a Kubernetes environment, aimed at enhancing API
observability. It includes detailed commands for each step along with their explanations.

> **Note**: SentryFlow is currently in the early stages of development. Please be aware that the information provided
> here may become outdated or change without notice.

## 1. Prerequisites

- A Kubernetes cluster running version 1.28 or later.
- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) version 1.28 or later.

## 2. Deploying SentryFlow

Configure SentryFlow receiver by following [this](receivers.md). Then deploy SentryFlow by following `kubectl` command:

```shell
kubectl apply -f https://raw.githubusercontent.com/5GSEC/SentryFlow/refs/heads/main/deployments/sentryflow.yaml
```

This will create a namespace named `sentryflow` and will deploy the necessary Kubernetes resources.

Then, check if SentryFlow is up and running by:

```shell
$ kubectl -n sentryflow get pods
NAME                         READY   STATUS    RESTARTS   AGE
sentryflow-cff887bbd-rljm7   1/1     Running   0          73s
```

## 3. Deploying SentryFlow Clients

SentryFlow has now been deployed in the cluster. In addition, SentryFlow exports API access logs through `gRPC`.

For testing purposes, a client has been developed.

- `log-client`: Simply logs everything on `STDOUT` coming from SentryFlow.

It can be deployed into the cluster under namespace `sentryflow` by following the command:

```shell
kubectl apply -f https://raw.githubusercontent.com/5GSEC/SentryFlow/refs/heads/main/deployments/sentryflow-client.yaml
```

Then, check if it is up and running by:

```shell
kubectl get pods -n sentryflow
```

If you observe `log-client`, is running, the setup has been completed successfully.
