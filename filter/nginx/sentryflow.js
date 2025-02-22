const DEFAULT_KEY = "sentryFlow";
const ResStatusKey = ":status"
const MAX_BODY_SIZE = 1_000_000; // 1 MB

function requestHandler(r, data, flags) {
    // https://nginx.org/en/docs/njs/reference.html#r_sendbuffer
    r.sendBuffer(data, flags);

    // https://nginx.org/en/docs/njs/reference.html#r_done
    r.done();

    let responseBody = ""
    try {
        responseBody = new TextDecoder("utf-8")
            .decode(new Uint8Array(data));
    } catch (error) {
        r.error(`failed to decode data, error: ${error}`)
        // Do not return, process other info even without body.
    }

    if (responseBody.length > MAX_BODY_SIZE) {
        responseBody = ""
    }

    let apiEvent = {
        "metadata": {
            // Divide by 1000 converts the timestamp from milliseconds to seconds.
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

    // https://nginx.org/en/docs/http/ngx_http_core_module.html#variables
    apiEvent.request.headers[":scheme"] = r.variables.scheme
    apiEvent.request.headers[":path"] = r.uri
    apiEvent.request.headers[":method"] = r.variables.request_method

    // Number of bytes sent to a client, not counting the response header; this
    // variable is compatible with the “%B” parameter of the mod_log_config Apache module.
    apiEvent.request.headers["body_bytes_sent"] = r.variables.body_bytes_sent

    // Request length including request line, header, and request body.
    apiEvent.request.headers["request_length"] = r.variables.request_length

    // Request processing time in seconds with a milliseconds resolution;
    // Time elapsed since the first bytes were read from the client.
    apiEvent.request.headers["request_time"] = r.variables.request_time

    // Query (args) in the request line.
    apiEvent.request.headers["query"] = r.variables.query_string

    for (const header in r.headersOut) {
        apiEvent.response.headers[header] = r.headersOut[header];
    }
    apiEvent.response.headers[ResStatusKey] = r.variables.status

    // https://nginx.org/en/docs/njs/reference.html#ngx_shared
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
