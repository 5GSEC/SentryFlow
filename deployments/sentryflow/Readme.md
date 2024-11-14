# Install SentryFlow

Install SentryFlow using the official 5GSEC Helm charts.

```shell
helm repo add 5gsec https://5gsec.github.io/charts
helm repo update 5gsec
helm upgrade --install sentryflow 5gsec/sentryflow-n sentryflow --create-namespace
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
