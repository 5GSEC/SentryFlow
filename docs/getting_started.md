# Getting Started

This guide provides a step-by-step process for deploying SentryFlow in a Kubernetes environment, aimed at enhancing API
observability. It includes detailed commands for each step along with their explanations.

> **Note**: SentryFlow is currently in the early stages of development. Please be aware that the information provided
> here may become outdated or change without notice.

## 1. Prerequisites

- A Kubernetes cluster running version 1.28 or later.
- [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) version 1.28 or later.

## 2. Deploying SentryFlow

- Add SentryFlow repo

```shell
helm repo add 5gsec https://5gsec.github.io/charts
helm repo update 5gsec
```

- Update `values.yaml` file as follows by following [this](receivers.md).

```shell
helm show values 5gsec/sentryflow > values.yaml
```

- Deploy SentryFlow

```shell
helm install --values values.yaml sentryflow 5gsec/sentryflow -n sentryflow --create-namespace 
```
