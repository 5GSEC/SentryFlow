# Install SentryFlow

Install SentryFlow using the official 5GSEC Helm charts.

### Add SentryFlow repo

```shell
helm repo add 5gsec https://5gsec.github.io/charts
helm repo update 5gsec
```

### Update `values.yaml` file according to your requirements.

```shell
helm show values 5gsec/sentryflow > values.yaml
```

Configure SentryFlow receiver by following [this](../../docs/receivers.md).

### Deploy SentryFlow

```shell
helm install --values values.yaml sentryflow 5gsec/sentryflow-n sentryflow --create-namespace 
```

Install SentryFlow using Helm charts locally (for testing)

```bash
cd deployments/sentryflow/
helm upgrade --install sentryflow . -n sentryflow --create-namespace
```

## Uninstall

To uninstall, just run:

```bash
helm uninstall sentryflow -n sentryflow
```
