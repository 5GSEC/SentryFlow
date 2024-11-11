use log::error;
use proxy_wasm::traits::{Context, HttpContext, RootContext};
use proxy_wasm::types::{Action, ContextType, LogLevel};
use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::time::{Duration, UNIX_EPOCH};

#[derive(Default)]
struct Plugin {
    _context_id: u32,
    config: PluginConfig,
    api_event: APIEvent,
}

#[derive(Deserialize, Clone, Default)]
struct PluginConfig {
    upstream_name: String,
    api_path: String,
    authority: String,
}

#[derive(Serialize, Default, Clone)]
struct APIEvent {
    metadata: Metadata,
    request: Reqquest,
    response: Ressponse,
    source: Workload,
    destination: Workload,
    protocol: String,
}

#[derive(Serialize, Default, Clone)]
struct Metadata {
    context_id: u32,
    timestamp: u64,
    istio_version: String,
    mesh_id: String,
    node_name: String,
}

#[derive(Serialize, Default, Clone)]
struct Workload {
    name: String,
    namespace: String,
    ip: String,
    port: u16,
}

#[derive(Serialize, Clone, Default)]
struct Reqquest {
    headers: HashMap<String, String>,
    body: String,
}

#[derive(Serialize, Clone, Default, Debug)]
struct Ressponse {
    headers: HashMap<String, String>,
    body: String,
}

const MAX_BODY_SIZE: usize = 1_000_000; // 1 MB

fn _start() {
    proxy_wasm::main! {{
        proxy_wasm::set_log_level(LogLevel::Warn);
        proxy_wasm::set_root_context(|_| -> Box<dyn RootContext> {Box::new(Plugin::default())});
    }}
}

impl Context for Plugin {
    fn on_done(&mut self) -> bool {
        dispatch_http_call_to_upstream(self);
        true
    }
}

impl RootContext for Plugin {
    fn on_configure(&mut self, _plugin_configuration_size: usize) -> bool {
        if let Some(config_bytes) = self.get_plugin_configuration() {
            if let Ok(config) = serde_json::from_slice::<PluginConfig>(&config_bytes) {
                self.config = config;
            } else {
                error!("Failed to parse plugin config");
            }
        } else {
            error!("No plugin config found");
        }
        true
    }

    fn create_http_context(&self, _context_id: u32) -> Option<Box<dyn HttpContext>> {
        Some(Box::new(Plugin {
            _context_id,
            config: self.config.clone(),
            api_event: Default::default(),
        }))
    }

    fn get_type(&self) -> Option<ContextType> {
        Some(ContextType::HttpContext)
    }
}

impl HttpContext for Plugin {
    fn on_http_request_headers(&mut self, _num_headers: usize, _end_of_stream: bool) -> Action {
        let (src_ip, src_port) = get_url_and_port(
            String::from_utf8(
                self.get_property(vec!["source", "address"])
                    .unwrap_or_default(),
            )
            .unwrap_or_default(),
        );

        let req_headers = self.get_http_request_headers();
        let mut headers: HashMap<String, String> = HashMap::with_capacity(req_headers.len());
        for header in req_headers {
            // Don't include Envoy's pseudo headers
            // https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_conn_man/headers#id13
            if !header.0.starts_with("x-envoy") {
                headers.insert(header.0, header.1);
            }
        }

        self.api_event.metadata.timestamp = self
            .get_current_time()
            .duration_since(UNIX_EPOCH)
            .unwrap_or_default()
            .as_secs();
        self.api_event.metadata.context_id = self._context_id;
        self.api_event.request.headers = headers;

        let protocol = String::from_utf8(
            self.get_property(vec!["request", "protocol"])
                .unwrap_or_default(),
        )
        .unwrap_or_default();
        self.api_event.protocol = protocol;

        self.api_event.source.ip = src_ip;
        self.api_event.source.port = src_port;
        self.api_event.source.name = String::from_utf8(
            self.get_property(vec!["node", "metadata", "NAME"])
                .unwrap_or_default(),
        )
        .unwrap_or_default();
        self.api_event.source.namespace = String::from_utf8(
            self.get_property(vec!["node", "metadata", "NAMESPACE"])
                .unwrap_or_default(),
        )
        .unwrap_or_default();

        Action::Continue
    }

    fn on_http_request_body(&mut self, _body_size: usize, _end_of_stream: bool) -> Action {
        let body = String::from_utf8(
            self.get_http_request_body(0, _body_size)
                .unwrap_or_default(),
        )
        .unwrap_or_default();

        if !body.is_empty() && body.len() <= MAX_BODY_SIZE {
            self.api_event.request.body = body;
        }
        Action::Continue
    }

    fn on_http_response_headers(&mut self, _num_headers: usize, _end_of_stream: bool) -> Action {
        let (dest_ip, dest_port) = get_url_and_port(
            String::from_utf8(
                self.get_property(vec!["destination", "address"])
                    .unwrap_or_default(),
            )
            .unwrap_or_default(),
        );

        let res_headers = self.get_http_response_headers();
        let mut headers: HashMap<String, String> = HashMap::with_capacity(res_headers.len());
        for res_header in res_headers {
            // Don't include Envoy's pseudo headers
            // https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_conn_man/headers#id13
            if !res_header.0.starts_with("x-envoy") {
                headers.insert(res_header.0, res_header.1);
            }
        }

        self.api_event.response.headers = headers;
        self.api_event.destination.ip = dest_ip;
        self.api_event.destination.port = dest_port;
        find_and_update_dest_namespace(self);

        Action::Continue
    }

    fn on_http_response_body(&mut self, _body_size: usize, _end_of_stream: bool) -> Action {
        let body = String::from_utf8(
            self.get_http_response_body(0, _body_size)
                .unwrap_or_default(),
        )
        .unwrap_or_default();
        if !body.is_empty() && body.len() <= MAX_BODY_SIZE {
            self.api_event.response.body = body;
        }
        Action::Continue
    }
}

fn find_and_update_dest_namespace(obj: &mut Plugin) {
    let dest_ns = String::from_utf8(
        obj.get_property(vec![
            "upstream_host_metadata",
            "filter_metadata",
            "istio",
            "workload",
        ])
        .unwrap_or_default(),
    )
    .unwrap_or_default();

    // e.g., filterserver;sentryflow;filterserver;;Kubernetes
    if !dest_ns.is_empty() {
        let parts: Vec<&str> = dest_ns.split(";").collect();
        if parts.len() == 5 || parts.len() == 4 {
            obj.api_event.destination.namespace = parts[1].to_string();
        }
    }
}

fn dispatch_http_call_to_upstream(obj: &mut Plugin) {
    update_metadata(obj);
    let telemetry_json = serde_json::to_string(&obj.api_event).unwrap_or_default();

    let headers = vec![
        (":method", "POST"),
        (":authority", &obj.config.authority),
        (":path", &obj.config.api_path),
        ("accept", "*/*"),
        ("Content-Type", "application/json"),
    ];

    let http_call_res = obj.dispatch_http_call(
        &obj.config.upstream_name,
        headers,
        Some(telemetry_json.as_bytes()),
        vec![],
        Duration::from_secs(1),
    );

    if http_call_res.is_err() {
        error!(
            "Failed to dispatch HTTP call, to '{}' status: {http_call_res:#?}",
            &obj.config.upstream_name,
        );
    }
}

fn update_metadata(obj: &mut Plugin) {
    obj.api_event.metadata.node_name = String::from_utf8(
        obj.get_property(vec!["node", "metadata", "NODE_NAME"])
            .unwrap_or_default(),
    )
    .unwrap_or_default();
    obj.api_event.metadata.mesh_id = String::from_utf8(
        obj.get_property(vec!["node", "metadata", "MESH_ID"])
            .unwrap_or_default(),
    )
    .unwrap_or_default();
    obj.api_event.metadata.istio_version = String::from_utf8(
        obj.get_property(vec!["node", "metadata", "ISTIO_VERSION"])
            .unwrap_or_default(),
    )
    .unwrap_or_default();
}

fn get_url_and_port(address: String) -> (String, u16) {
    let parts: Vec<&str> = address.split(':').collect();

    let mut url = "".to_string();
    let mut port = 0;

    if parts.len() == 2 {
        url = parts[0].parse().unwrap();
        port = parts[1].parse().unwrap();
    } else {
        error!("Invalid address");
    }

    (url, port)
}
