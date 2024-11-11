// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package core

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.uber.org/zap"

	protobuf "github.com/5GSEC/SentryFlow/protobuf/golang"
)

func Test_healthzHandler(t *testing.T) {
	tests := []struct {
		name   string
		method string
		want   int
	}{
		{
			name:   "with GET method should return StatusOK",
			method: http.MethodGet,
			want:   http.StatusOK,
		},
		{
			name:   "with POST method should return StatusMethodNotAllowed",
			method: http.MethodPost,
			want:   http.StatusMethodNotAllowed,
		},
		{
			name:   "with PUT method should return StatusMethodNotAllowed",
			method: http.MethodPut,
			want:   http.StatusMethodNotAllowed,
		},
		{
			name:   "with DELETE method should return StatusMethodNotAllowed",
			method: http.MethodDelete,
			want:   http.StatusMethodNotAllowed,
		},
		{
			name:   "with PATCH method should return StatusMethodNotAllowed",
			method: http.MethodPatch,
			want:   http.StatusMethodNotAllowed,
		},
		{
			name:   "with HEAD method should return StatusMethodNotAllowed",
			method: http.MethodHead,
			want:   http.StatusMethodNotAllowed,
		},
		{
			name:   "with TRACE method should return StatusMethodNotAllowed",
			method: http.MethodTrace,
			want:   http.StatusMethodNotAllowed,
		},
		{
			name:   "with OPTIONS method should return StatusMethodNotAllowed",
			method: http.MethodOptions,
			want:   http.StatusMethodNotAllowed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := &Manager{}
			request := httptest.NewRequest(tt.method, "/healthz", nil)
			response := httptest.NewRecorder()
			manager.healthzHandler(response, request)
			if got := response.Code; got != tt.want {
				t.Errorf("healthzHandler() gotStatusCode = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_eventsHandler(t *testing.T) {
	logger := zap.S()

	validApiEvent := getDummyValidApiEvent()
	invalidApiEvent := getDummyInvalidApiEvent()

	type fields struct {
		Logger    *zap.SugaredLogger
		ApiEvents chan *protobuf.APIEvent
	}
	tests := []struct {
		name           string
		fields         fields
		body           []byte
		method         string
		wantApiEvent   []byte
		wantStatusCode int
	}{
		{
			name: "with valid apiEvent and GET method should return StatusMethodNotAllowed",
			fields: fields{
				Logger:    logger,
				ApiEvents: make(chan *protobuf.APIEvent, 1),
			},
			body:           validApiEvent,
			method:         http.MethodGet,
			wantApiEvent:   nil,
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name: "with valid apiEvent and POST method should return StatusAccepted",
			fields: fields{
				Logger:    logger,
				ApiEvents: make(chan *protobuf.APIEvent, 1),
			},
			body:           validApiEvent,
			method:         http.MethodPost,
			wantApiEvent:   validApiEvent,
			wantStatusCode: http.StatusAccepted,
		},
		{
			name: "with valid apiEvent and PUT method should return StatusMethodNotAllowed",
			fields: fields{
				Logger:    logger,
				ApiEvents: make(chan *protobuf.APIEvent, 1),
			},
			body:           validApiEvent,
			method:         http.MethodPut,
			wantApiEvent:   nil,
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name: "with valid apiEvent and DELETE method should return StatusMethodNotAllowed",
			fields: fields{
				Logger:    logger,
				ApiEvents: make(chan *protobuf.APIEvent, 1),
			},
			body:           validApiEvent,
			method:         http.MethodDelete,
			wantApiEvent:   nil,
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name: "with valid apiEvent and PATCH method should return StatusMethodNotAllowed",
			fields: fields{
				Logger:    logger,
				ApiEvents: make(chan *protobuf.APIEvent, 1),
			},
			body:           validApiEvent,
			method:         http.MethodPatch,
			wantApiEvent:   nil,
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name: "with valid apiEvent and OPTIONS method should return StatusMethodNotAllowed",
			fields: fields{
				Logger:    logger,
				ApiEvents: make(chan *protobuf.APIEvent, 1),
			},
			body:           validApiEvent,
			method:         http.MethodOptions,
			wantApiEvent:   nil,
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name: "with valid apiEvent and TRACE method should return StatusMethodNotAllowed",
			fields: fields{
				Logger:    logger,
				ApiEvents: make(chan *protobuf.APIEvent, 1),
			},
			body:           validApiEvent,
			method:         http.MethodTrace,
			wantApiEvent:   nil,
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name: "with valid apiEvent and TRACE method should return StatusMethodNotAllowed",
			fields: fields{
				Logger:    logger,
				ApiEvents: make(chan *protobuf.APIEvent, 1),
			},
			body:           validApiEvent,
			method:         http.MethodHead,
			wantApiEvent:   nil,
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name: "with empty apiEvent body should return StatusBadRequest",
			fields: fields{
				Logger:    logger,
				ApiEvents: make(chan *protobuf.APIEvent, 1),
			},
			body:           nil,
			method:         http.MethodPost,
			wantApiEvent:   nil,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "with invalid apiEvent body should return StatusBadRequest",
			fields: fields{
				Logger:    logger,
				ApiEvents: make(chan *protobuf.APIEvent, 1),
			},
			body:           invalidApiEvent,
			method:         http.MethodPost,
			wantApiEvent:   nil,
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Avoid using `t.Cleanup()` as it executes after all subtests complete. To
			// prevent processing stale API events, close the `apiEvent` channel proactively.
			defer close(tt.fields.ApiEvents)

			m := &Manager{
				Logger:    tt.fields.Logger,
				ApiEvents: tt.fields.ApiEvents,
			}

			request := httptest.NewRequest(tt.method, "/api/v1/events", bytes.NewReader(tt.body))
			response := httptest.NewRecorder()
			m.eventsHandler(response, request)

			if gotStatusCode := response.Code; gotStatusCode != tt.wantStatusCode {
				t.Errorf("eventsHandler() gotStatusCode = %v, want %v", gotStatusCode, tt.wantStatusCode)
			}
			if len(tt.fields.ApiEvents) > 0 {
				gotApiEvent, _ := json.Marshal(<-m.ApiEvents)
				if !bytes.Equal(gotApiEvent, tt.wantApiEvent) {
					t.Errorf("eventsHandler() gotApiEvent = %v, want %v", gotApiEvent, tt.wantApiEvent)
				}
			}
		})
	}
}

func getDummyInvalidApiEvent() []byte {
	apiEvent := `
{
  "metadata": {
	"timestamp": 1729179252
  },
  "http": {
	"request": {
	  "method": "GET",
	  "path": "/_-flags-1x1-hr.svg"
	},
	"response": {
	  "status_code": 200
	}
  },
}
`
	body, _ := json.Marshal(apiEvent)
	return body
}

func getDummyValidApiEvent() []byte {
	apiEvent := &protobuf.APIEvent{
		Metadata: &protobuf.Metadata{
			ContextId: 1,
			Timestamp: uint64(time.Now().Unix()),
		},
		Source: &protobuf.Workload{
			Name:      "source-workload",
			Namespace: "source-namespace",
			Ip:        "1.1.1.1",
			Port:      11111,
		},
		Destination: &protobuf.Workload{
			Name:      "destination-workload",
			Namespace: "destination-namespace",
			Ip:        "93.184.215.14",
			Port:      80,
		},
		Request: &protobuf.Request{
			Headers: map[string]string{
				":authority": "example.com",
				":method":    "GET",
				":path":      "/",
				":scheme":    "http",
			},
			Body: "request body",
		},
		Response: &protobuf.Response{
			Headers: map[string]string{
				":status": "200",
			},
			Body: "response body",
		},
		Protocol: "HTTP/1.1",
	}
	body, _ := json.Marshal(apiEvent)
	return body
}
