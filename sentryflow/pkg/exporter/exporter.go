// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package exporter

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"

	"github.com/5GSEC/SentryFlow/pkg/config"
	"github.com/5GSEC/SentryFlow/pkg/util"
	protobuf "github.com/5GSEC/SentryFlow/protobuf/golang"
)

// clientList represents a list of gRPC clients and their associated channels for
// sending API events. It uses a mutex to synchronize access to the client map.
type clientList struct {
	*sync.Mutex
	client map[string]chan *protobuf.APIEvent
}

type exporter struct {
	protobuf.UnimplementedSentryFlowServer
	apiEvents chan *protobuf.APIEvent
	logger    *zap.SugaredLogger
	clients   *clientList
}

// GetAPIEvent streams generated API events to connected clients. Each client is
// assigned a unique identifier (UID) and a dedicated channel to receive events.
// This ensures that all connected clients receive the same API events in
// real-time.
func (e *exporter) GetAPIEvent(clientInfo *protobuf.ClientInfo, stream grpc.ServerStreamingServer[protobuf.APIEvent]) error {
	uid := uuid.Must(uuid.NewRandom()).String()

	connChan := e.addClientToList(uid)
	defer e.deleteClientFromList(uid, connChan)

	e.logger.Infof("Client: %s %s (%s) connected", uid, clientInfo.HostName, clientInfo.IPAddress)

	for {
		select {
		case <-stream.Context().Done():
			e.logger.Infof("Client: %s %s (%s) disconnected", uid, clientInfo.HostName, clientInfo.IPAddress)
			return stream.Context().Err()
		case apiEvent, ok := <-connChan:
			if !ok {
				e.logger.Warn("Channel closed")
				return nil
			}
			if status, ok := grpcstatus.FromError(stream.Send(apiEvent)); !ok {
				if status.Code() == codes.Canceled {
					e.logger.Infof("Client: %s %s (%s) cancelled the operation", uid, clientInfo.HostName, clientInfo.IPAddress)
					return nil
				}
				e.logger.Errorf("Failed to send APIEvent: %v", status.Err())
				return status.Err()
			}
		}
	}
}

func (e *exporter) addClientToList(uid string) chan *protobuf.APIEvent {
	e.clients.Lock()
	connChan := make(chan *protobuf.APIEvent)
	e.clients.client[uid] = connChan
	e.clients.Unlock()
	return connChan
}

func (e *exporter) deleteClientFromList(uid string, connChan chan *protobuf.APIEvent) {
	e.clients.Lock()
	close(connChan)
	delete(e.clients.client, uid)
	e.clients.Unlock()
}

// SendAPIEvent ingests an API event received from the source and publishes it to
// the `apiEvents` channel for subscribed clients to consume.
func (e *exporter) SendAPIEvent(ctx context.Context, apiEvent *protobuf.APIEvent) (*protobuf.APIEvent, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		e.apiEvents <- apiEvent
		e.logger.Info("Received APIEvent")
		return apiEvent, nil
	}
}

// putApiEventOnClientsChannel continuously listens to the `apiEvents` channel
// and forwards incoming API events to all connected clients. If the context is
// canceled, the function returns.
func (e *exporter) putApiEventOnClientsChannel(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case apiEvent, ok := <-e.apiEvents:
			if !ok {
				e.logger.Warn("Channel closed")
				continue
			}
			eventToSend := apiEvent
			e.clients.Lock()
			for _, clientChan := range e.clients.client {
				select {
				case clientChan <- eventToSend:
				default:
					e.logger.Warn("Event dropped")
				}
			}
			e.clients.Unlock()
		}
	}
}

// Init initializes and registers the gRPC-based exporter with the provided
// server. This allows clients to connect and consume the generated API events
// streamed through the server.
func Init(ctx context.Context, server *grpc.Server, cfg *config.Config, events chan *protobuf.APIEvent, wg *sync.WaitGroup) error {
	logger := util.LoggerFromCtx(ctx).Named("exporter")
	logger.Info("Starting exporter")

	e := &exporter{
		apiEvents: events,
		logger:    logger,
		clients: &clientList{
			Mutex:  &sync.Mutex{},
			client: make(map[string]chan *protobuf.APIEvent),
		},
	}

	protobuf.RegisterSentryFlowServer(server, e)

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		e.putApiEventOnClientsChannel(ctx)
	}(ctx)

	return nil
}
