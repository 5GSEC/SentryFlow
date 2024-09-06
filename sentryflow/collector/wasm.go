// SPDX-License-Identifier: Apache-2.0

package collector

import (
	"encoding/json"
	"github.com/5gsec/SentryFlow/processor"
	"github.com/5gsec/SentryFlow/protobuf"
	"io/ioutil"
	"log"
	"net/http"
)

// Handler for the HTTP endpoint to receive api events from WASM filter
func DataHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the request is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	// Parse the JSON data into the TelemetryData struct
	var apiLog *protobuf.APILogV2
	err = json.Unmarshal(body, &apiLog)
	if err != nil {
		log.Print("failed to parse json")
		log.Print(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Check protocol version
	if r.ProtoMajor == 2 {
		apiLog.Protocol = "HTTP/2.0"
	} else if r.ProtoMajor == 1 && r.ProtoMinor == 1 {
		apiLog.Protocol = "HTTP/1.1"
	} else if r.ProtoMajor == 1 && r.ProtoMinor == 0 {
		apiLog.Protocol = "HTTP/1.0"
	} else {
		apiLog.Protocol = "Unknown"
	}
	processor.InsertAPILog(apiLog)

	// Log the received telemetry data
	log.Printf("Received data: %+v\n", apiLog)
}
