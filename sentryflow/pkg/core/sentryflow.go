// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package core

import (
	"context"
	"net/http"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	istionet "istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	"github.com/5GSEC/SentryFlow/pkg/config"
	"github.com/5GSEC/SentryFlow/pkg/exporter"
	"github.com/5GSEC/SentryFlow/pkg/k8s"
	"github.com/5GSEC/SentryFlow/pkg/receiver"
	"github.com/5GSEC/SentryFlow/pkg/util"
	"github.com/5GSEC/SentryFlow/protobuf"
)

type Manager struct {
	Ctx        context.Context
	Logger     *zap.SugaredLogger
	GrpcServer *grpc.Server
	HttpServer *http.Server
	K8sClient  client.Client
	Wg         *sync.WaitGroup
	ApiEvents  chan *protobuf.APIEvent
}

func Run(ctx context.Context, configFilePath string, kubeConfig string) {
	mgr := &Manager{
		Ctx:        ctx,
		Logger:     util.LoggerFromCtx(ctx),
		GrpcServer: grpc.NewServer(),
		Wg:         &sync.WaitGroup{},
		ApiEvents:  make(chan *protobuf.APIEvent, 10240),
	}
	mgr.Logger.Info("Starting SentryFlow")

	cfg, err := config.New(configFilePath, mgr.Logger)
	if err != nil {
		mgr.Logger.Fatal(err)
	}

	k8sClient, err := k8s.NewClient(registerAndGetScheme(), kubeConfig)
	if err != nil {
		mgr.Logger.Fatalf("Failed to create k8s client: %v", err)
	}
	mgr.K8sClient = k8sClient

	mgr.Wg.Add(1)
	go func() {
		defer mgr.Wg.Done()
		mgr.startHttpServer()
	}()

	if err := receiver.Init(mgr.Ctx, mgr.K8sClient, cfg.Receivers, mgr.ApiEvents, mgr.GrpcServer, mgr.Wg); err != nil {
		mgr.Logger.Fatalf("failed to initialize receiver: %v", err)
	}

	if err := exporter.Init(mgr.Ctx, mgr.GrpcServer, cfg.Exporter, mgr.ApiEvents, mgr.Wg); err != nil {
		mgr.Logger.Fatalf("Failed to initialize exporter: %v", err)
	}

	mgr.Wg.Add(1)
	go func() {
		defer mgr.Wg.Done()
		mgr.startGrpcServer(cfg.Exporter.Grpc.Port)
	}()

	mgr.Logger.Info("Started SentryFlow")

	<-ctx.Done()
	mgr.Logger.Info("Shutdown Signal Received. Waiting for all workers to finish.")
	mgr.Logger.Info("Shutting down SentryFlow")

	mgr.stopServers()
	mgr.Wg.Wait()
	close(mgr.ApiEvents)

	mgr.Logger.Info("All workers finished. Stopped SentryFlow")
}

func registerAndGetScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	utilruntime.Must(istionet.AddToScheme(scheme))
	return scheme
}