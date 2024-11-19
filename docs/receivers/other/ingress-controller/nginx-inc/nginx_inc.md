# Nginx Incorporation Ingress Controller

## Description

This guide provides a step-by-step process to integrate SentryFlow
with [Nginx Inc.](https://docs.nginx.com/nginx-ingress-controller/) Ingress Controller, aimed at enhancing API
observability. It includes detailed commands for each step along with their explanations.

SentryFlow make use of following to provide visibility into API calls:

- [Nginx njs](https://nginx.org/en/docs/njs/) module.
- [Njs filter](../../../../../filter/nginx).

## Prerequisites

- Nginx Inc. Ingress Controller.
  Follow [this](https://docs.nginx.com/nginx-ingress-controller/installation/installing-nic/) to deploy it.

## How to

To Observe API calls of your workloads served by Nginx inc. ingress controller in Kubernetes environment, follow
the below
steps:

1. Create the following configmap in the same namespace as ingress controller.

```shell
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ConfigMap
metadata:
  name: sentryflow-nginx-inc
  namespace: <ingress-controller-namespace>
data:
  sentryflow.js: |
    const DEFAULT_KEY = "sentryFlow";
    const ResStatusKey = ":status"
    const MAX_BODY_SIZE = 1_000_000; // 1 MB
    
    function requestHandler(r, data, flags) {
        r.sendBuffer(data, flags);
        r.done();
        
        let responseBody = ""
        try {
            responseBody = new TextDecoder("utf-8")
                .decode(new Uint8Array(data));
        } catch (error) {
            r.error(`failed to decode data, error: ${error}`)
        }
        
        if (responseBody.length > MAX_BODY_SIZE) {
            responseBody = ""
        }
        
        let apiEvent = {
            "metadata": {
                "timestamp": Date.parse(r.variables.time_iso8601.split("+")[0]) / 1000,
                "receiver_name": "nginx",
                "receiver_version": ngx.version,
            },
            "source": {
                "ip": r.remoteAddress,
                "port": r.variables.remote_port,
            },
            "destination": {
                "ip": r.variables.server_addr,
                "port": r.variables.server_port,
            },
            "request": {
                "headers": {},
                "body": r.requestText || "",
            },
            "response": {
                "headers": {},
                "body": responseBody,
            },
            "protocol": r.variables.server_protocol,
        };
        
        for (const header in r.headersIn) {
            apiEvent.request.headers[header] = r.headersIn[header];
        }
        
        apiEvent.request.headers[":scheme"] = r.variables.scheme
        apiEvent.request.headers[":path"] = r.uri
        apiEvent.request.headers[":method"] = r.variables.request_method
        
        apiEvent.request.headers["body_bytes_sent"] = r.variables.body_bytes_sent
        
        apiEvent.request.headers["request_length"] = r.variables.request_length
        
        apiEvent.request.headers["request_time"] = r.variables.request_time
        
        apiEvent.request.headers["query"] = r.variables.query_string
        
        for (const header in r.headersOut) {
            apiEvent.response.headers[header] = r.headersOut[header];
        }
        apiEvent.response.headers[ResStatusKey] = r.variables.status
        
        ngx.shared.apievents.set(DEFAULT_KEY, JSON.stringify(apiEvent));
    }
    
    async function dispatchHttpCall(r) {
        try {
            let apiEvent = ngx.shared.apievents.get(DEFAULT_KEY);
            await r.subrequest("/sentryflow", {
                method: "POST", body: apiEvent, detached: true
            })
        } catch (error) {
            r.error(`failed to dispatch HTTP call to SentryFlow, error: ${error}`)
            return;
        } finally {
            ngx.shared.apievents.clear();
        }
        
        r.return(200, "OK");
    }
    
    export default {requestHandler, dispatchHttpCall};
EOF
```

2. Add the following volume and volume-mount in ingress controller deployment:

```yaml
...
volumes:
  - name: sentryflow-nginx-inc
    configMap:
      name: sentryflow-nginx-inc
...
...
volumeMounts:
  - mountPath: /etc/nginx/njs/sentryflow.js
    name: sentryflow-nginx-inc
    subPath: sentryflow.js
```

3. Update ingress controller configmap as follows:

```yaml
...
data:
  http-snippets: |
    js_path "/etc/nginx/njs/";
    subrequest_output_buffer_size 8k;
    js_shared_dict_zone zone=apievents:1M timeout=300s evict;
    js_import main from sentryflow.js;
  location-snippets: |
    js_body_filter main.requestHandler buffer_type=buffer;
    mirror      /mirror_request;
    mirror_request_body on;
  server-snippets: |
    location /mirror_request {
      internal;
      js_content main.dispatchHttpCall;
    }
    location /sentryflow {
      internal;
      # Update SentryFlow URL with path to ingest access logs if required.
      proxy_pass http://sentryflow.sentryflow:8081/api/v1/events;
      proxy_method      POST;
      proxy_set_header accept "application/json";
      proxy_set_header Content-Type "application/json";
    }
```

4. Download SentryFlow manifest file

  ```shell
  curl -sO https://raw.githubusercontent.com/5GSEC/SentryFlow/refs/heads/main/deployments/sentryflow.yaml
  ```

5. Update the `.receivers` configuration in `sentryflow` [configmap](../../../../../deployments/sentryflow.yaml) as
   follows:

  ```yaml
  filters:
    server:
      port: 8081
    # Following is required for `nginx-inc-ingress-controller` receiver.  
    nginxIngress:
      deploymentName: <nginx-ingress-controller-deploy-name>
      configMapName: <nginx-ingress-configmap-name>
      sentryFlowNjsConfigMapName: <sentryflow-nginx-inc-configmap-name>

  receivers:
    others:
      - name: nginx-inc-ingress-controller # SentryFlow makes use of `name` to configure receivers. DON'T CHANGE IT.
        namespace: <ingress-controller-namespace> # Kubernetes namespace in which you've deployed the ingress controller.
    ...
  ```

6. Deploy SentryFlow

  ```shell
  kubectl apply -f sentryflow.yaml
  ```

7. Trigger API calls to generate traffic.

8. Use SentryFlow [log client](../../../../client) to see the API Events.
