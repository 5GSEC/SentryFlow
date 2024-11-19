# Nginx

This Nginx JavaScript script actively monitors API calls made to an application running on a virtual machine (VM) behind
a Nginx web server.

## Sample API Event:

```json
{
  "metadata": {
    "timestamp": 1728722194,
    "receiver_name": "nginx",
    "receiver_version": "1.26.2"
  },
  "source": {
    "ip": "192.168.64.1",
    "port": "58242"
  },
  "destination": {
    "ip": "192.168.64.19",
    "port": "80"
  },
  "request": {
    "headers": {
      "Host": "192.168.64.19",
      "Connection": "keep-alive",
      "Upgrade-Insecure-Requests": "1",
      "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36",
      "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8",
      "Sec-GPC": "1",
      "Accept-Language": "en-GB,en",
      "Accept-Encoding": "gzip, deflate",
      ":scheme": "http",
      ":path": "/api/v1/healthz",
      ":method": "GET",
      "body_bytes_sent": "0",
      "request_length": "440",
      "request_time": "0.000",
      "query": "just=for&testing=purpose"
    },
    "body": ""
  },
  "response": {
    "headers": {
      "Content-Type": "text/html",
      "Content-Length": "555",
      "status": "404"
    },
    "body": "<html>\r\n<head><title>404 Not Found</title></head>\r\n<body>\r\n<center><h1>404 Not Found</h1></center>\r\n"
  },
  "protocol": "HTTP/1.1"
}
```

# Getting Started

- Install [nginx-njs-module](https://github.com/nginx/njs?tab=readme-ov-file#downloading-and-installing)
- Copy [sentryflow.js](sentryflow.js) file to `/etc/nginx/njs/` directory as `sentryflow.js`.
- Edit `nginx.conf` file located in `/etc/nginx/` directory as follows:

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

Here is the sample [nginx.conf](nginx.conf) file for reference.

- Reload `nginx`:

```shell
$ sudo nginx -s reload
```

- Trigger API calls to generate traffic.
- Verify that the recorded API events in SentryFlow are similar to [sample event](#sample-api-event).
