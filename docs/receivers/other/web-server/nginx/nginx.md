# Nginx Web Server

## Description

This guide provides a step-by-step process to integrate SentryFlow
with [Nginx webserver](https://nginx.org/), aimed at enhancing API
observability. It includes detailed commands for each step along with their explanations.

SentryFlow make use of following to provide visibility into API calls:

- [Nginx njs](https://nginx.org/en/docs/njs/) module.
- [Njs filter](../../../../../filter/nginx).

## Prerequisites

- Nginx web server.
- [Nginx-njs-module](https://github.com/nginx/njs?tab=readme-ov-file#downloading-and-installing).

## How to

To Observe API calls of your application running on a virtual machine (VM) behind a Nginx web server, follow the below
steps:

1. Copy [sentryflow.js](../../../../../filter/nginx/sentryflow.js) file to `/etc/nginx/njs/` directory as
   `sentryflow.js`.
2. Edit `nginx.conf` file located in `/etc/nginx/` directory as follows:

```nginx configuration
load_module /etc/nginx/modules/ngx_http_js_module.so;
...
http {
    ...
    subrequest_output_buffer_size 8k;
    js_path "/etc/nginx/njs/";
    js_shared_dict_zone zone=apievents:1M timeout=60s evict;
    js_import main from sentryflow.js;
    ...
    server {
        location / {
            js_body_filter main.requestHandler buffer_type=buffer;
            mirror      /mirror_request;
            mirror_request_body on;
        }
        
        location /mirror_request {
            internal;
            js_content main.dispatchHttpCall;
        }
        
        location /sentryflow {
            internal;
            
            # SentryFlow URL with path to ingest access logs.
            proxy_pass http://<sentryflow_url>/api/v1/events;
            
            proxy_method      POST;
            proxy_set_header accept "application/json";
            proxy_set_header Content-Type "application/json";
        }
        ...
    }
} 
```

Here is the sample [nginx.conf](../../../../../filter/nginx/nginx.conf) file for reference.

3. Reload `nginx`:

```shell
$ sudo nginx -s reload
```

4. Deploy SentryFlow

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
receivers:
  others:
    - name: nginx-webserver # SentryFlow makes use of `name` to configure receivers. DON'T CHANGE IT.
    # Existing snippets
```

- Deploy SentryFlow

```shell
helm install --values values.yaml sentryflow 5gsec/sentryflow -n sentryflow --create-namespace 
```

5. Trigger API calls to generate traffic.
