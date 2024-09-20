// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package receiver

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/grpc"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/5GSEC/SentryFlow/pkg/config"
	istiosidecar "github.com/5GSEC/SentryFlow/pkg/receiver/svcmesh/istio/sidecar"
	"github.com/5GSEC/SentryFlow/pkg/util"
	"github.com/5GSEC/SentryFlow/protobuf"
)

// Init initializes the API event sources based on the provided configuration. It
// starts monitoring from configured sources and supports adding other sources in
// the future.
func Init(ctx context.Context, k8sClient client.Client, cfg *config.Receivers,
	apiEvents chan *protobuf.APIEvent, server *grpc.Server, wg *sync.WaitGroup) error {
	//logger := util.LoggerFromCtx(ctx).Named("receiver")

	for _, serviceMesh := range cfg.ServiceMeshes {
		if serviceMesh.Name != "" && serviceMesh.Enable {
			switch serviceMesh.Name {
			case util.ServiceMeshIstioSidecar:
				wg.Add(1)
				go func() {
					defer wg.Done()
					istiosidecar.StartMonitoring(ctx, k8sClient)
				}()
			default:
				return fmt.Errorf("unsupported Service Mesh, %v", serviceMesh.Name)
			}
		}
	}

	// Placeholder for other sources (To be implemented based on requirements)
	// TODO: Implement initialization for other telemetry sources based on the
	// `cfg.Others` configuration.
	//	This would involve handling gRPC or HTTP configs
	// for each supported source type and potentially adding new subdirectories in
	// `pkg/receiver/other` for each source.

	return nil
}
