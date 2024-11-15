// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package exporter

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/test/bufconn"
	"k8s.io/apimachinery/pkg/util/rand"

	protobuf "github.com/5GSEC/SentryFlow/protobuf/golang"
)

func Test_exporter_GetAPIEvent(t *testing.T) {
	// Use timeout to make sure this doesn't run indefinitely.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	e := getExporter()

	sfClient, closer := getSentryFlowClientAndCloser(t, e)
	defer closer()

	// Given
	stream, err := sfClient.GetAPIEvent(ctx, getClientInfo(t))
	if err != nil {
		t.Fatal(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	// Simulate some API events generation
	want := 100
	go func(numOfEvents int) {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)
		for i := 0; i < numOfEvents; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				e.apiEvents <- getDummyApiEvent(i)
				time.Sleep(10 * time.Millisecond)
			}
		}
	}(want)

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-ctx.Done():
			return
		default:
			e.putApiEventOnClientsChannel(ctx)
		}
	}()

	// When
	got := 0
	for got < want {
		_, err := stream.Recv()
		if err != nil {
			break
		}
		got++
	}

	// Then
	if got != want {
		t.Errorf("GetAPIEvent() want = %v events, got = %v", want, got)
	}

	cancel()
	wg.Wait()
}

func Test_exporter_SendAPIEvent(t *testing.T) {
	e := getExporter()

	// Given
	sfClient, closer := getSentryFlowClientAndCloser(t, e)
	defer closer()

	// When
	eventToSend := getDummyApiEvent(1)
	receivedEvent, err := sfClient.SendAPIEvent(context.Background(), eventToSend)
	if err != nil {
		t.Error(err)
	}

	// Since API events are reference types and their addresses differ, we need to
	// compare their serialized representations to ensure equality. So serialize both
	// events and then compare.
	want, _ := json.Marshal(eventToSend)
	got, _ := json.Marshal(receivedEvent)

	// Then
	if !bytes.Equal(want, got) {
		t.Errorf("SendAPIEvent() want = %v, got = %v", string(want), string(got))
	}
}

func Test_exporter_addClientToList(t *testing.T) {
	// Given
	e := exporter{
		clients: &clientList{
			Mutex:  &sync.Mutex{},
			client: make(map[string]chan *protobuf.APIEvent),
		},
	}
	uid := uuid.Must(uuid.NewRandom()).String()

	// When
	want := e.addClientToList(uid)

	// Then
	e.clients.Lock()
	got, exists := e.clients.client[uid]
	e.clients.Unlock()
	if !exists || got != want || got == nil {
		t.Errorf("addClientToList() client not added to the client list correctly")
	}
}

func Test_exporter_deleteClientFromList(t *testing.T) {
	// Given
	e := exporter{
		clients: &clientList{
			Mutex:  &sync.Mutex{},
			client: make(map[string]chan *protobuf.APIEvent),
		},
	}
	uid := uuid.Must(uuid.NewRandom()).String()

	// When
	e.deleteClientFromList(uid, e.addClientToList(uid))

	// Then
	e.clients.Lock()
	got, exists := e.clients.client[uid]
	e.clients.Unlock()
	if exists || got != nil {
		t.Errorf("deleteClientFromList() client not deleted from the client list correctly")
	}
}

func Test_exporter_add_and_delete_client_fromList_concurrently(t *testing.T) {
	// Given
	e := exporter{
		clients: &clientList{
			Mutex:  &sync.Mutex{},
			client: make(map[string]chan *protobuf.APIEvent),
		},
	}

	numOfClients := 1000
	wg := sync.WaitGroup{}
	wg.Add(numOfClients)

	// When
	for i := 0; i < numOfClients; i++ {
		go func() {
			defer wg.Done()

			uid := uuid.Must(uuid.NewRandom()).String()
			connChan := e.addClientToList(uid)

			// Simulate some work
			time.Sleep(time.Duration(rand.IntnRange(1, 100)) * time.Millisecond)

			e.deleteClientFromList(uid, connChan)
		}()
	}

	wg.Wait()

	// Then
	e.clients.Lock()
	if len(e.clients.client) != 0 {
		t.Errorf("client list is not empty after concurrent access")
	}
	e.clients.Unlock()
}

func getDummyApiEvent(ctxId int) *protobuf.APIEvent {
	return &protobuf.APIEvent{
		Metadata: &protobuf.Metadata{
			ContextId: uint32(ctxId),
			Timestamp: uint64(time.Now().Unix()),
		},
		Source: &protobuf.Workload{
			Name:      "source-workload",
			Namespace: "source-namespace",
			Ip:        "1.1.1.1",
			Port:      int32(rand.IntnRange(1025, 65536)),
		},
		Destination: &protobuf.Workload{
			Name:      "destination-workload",
			Namespace: "destination-namespace",
			Ip:        "93.184.215.14",
			Port:      int32(rand.IntnRange(80, 65536)),
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
}

func getExporter() *exporter {
	return &exporter{
		apiEvents: make(chan *protobuf.APIEvent, 100),
		logger:    zap.S(),
		clients: &clientList{
			Mutex:  &sync.Mutex{},
			client: make(map[string]chan *protobuf.APIEvent),
		},
	}
}

func getClientInfo(t *testing.T) *protobuf.ClientInfo {
	hostname, err := os.Hostname()
	if err != nil {
		t.Errorf("failed to get hostname: %v", err)
	}

	ips, err := net.LookupIP(hostname)
	if err != nil {
		t.Errorf("failed to get IP address: %v", err)
	}
	var ip string
	if len(ips) > 0 {
		ip = ips[0].String()
	}

	clientInfo := &protobuf.ClientInfo{
		HostName:  hostname,
		IPAddress: ip,
	}

	return clientInfo
}

func getSentryFlowClientAndCloser(t *testing.T, e *exporter) (protobuf.SentryFlowClient, func()) {
	listener := bufconn.Listen(101024 * 1024)
	baseServer := grpc.NewServer()
	protobuf.RegisterSentryFlowServer(baseServer, e)
	go func() {
		if err := baseServer.Serve(listener); err != nil {
			t.Errorf("failed to start exporter server: %v", err)
			return
		}
	}()

	resolver.SetDefaultScheme("passthrough")
	conn, err := grpc.NewClient("bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Errorf("failed to dial to server: %v", err)
		return nil, nil
	}

	closer := func() {
		if err := listener.Close(); err != nil {
			t.Errorf("failed to close listener: %v", err)
		}
		baseServer.Stop()
	}

	return protobuf.NewSentryFlowClient(conn), closer
}
