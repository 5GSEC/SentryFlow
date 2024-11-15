# SentryFlow Receivers

SentryFlow supports following receivers:

## Kubernetes

- [Istio sidecar](https://istio.io/latest/docs/setup/) service mesh. To integrate SentryFlow with it, refer
  to [this](receivers/service-mesh/istio/istio.md).
- [Nginx Inc.](https://github.com/nginxinc/kubernetes-ingress/) ingress controller. To integrate SentryFlow with it,
  refer to [this](receivers/other/ingress-controller/nginx-inc/nginx_inc.md).

## Non-Kubernetes

- [Nginx web server](https://github.com/nginx/nginx) running on Virtual Machine or Bare-Metal. To integrate SentryFlow
  with it, refer to [this](receivers/other/web-server/nginx/nginx.md).
