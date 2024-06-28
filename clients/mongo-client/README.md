# Mongo Client
Mongo client collects AccessLogs and Metrics from SentryFlow and stores them to database.

## Mongo Client Deployment
Mongo client can be deployed using kubectl command. The deployment can be accomplished with the following
commands:
```bash
$ cd SentryFlow/deployments
$ kubectl apply -f mongo-client.yaml
```

## Mongo client options
These are the default env value.
```bash
env:
- LOG_CFG: "stdout"
- METRIC_CFG: "stdout"
- METRIC_FILTER: "api"
```

If you want to change the default env value, you can refer to the following options.
```bash
env:
- name: LOG_CFG
  value: {"mongodb"|"none"}
- name: METRIC_CFG
  value: {"mongodb"|"none"}
- name: METRIC_FILTER
  value: {"all"|"api"|"envoy"}
```
