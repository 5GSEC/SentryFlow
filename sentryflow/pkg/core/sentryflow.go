// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package core

import (
	"context"
	"net/http"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"istio.io/client-go/pkg/apis/extensions/v1alpha1"
	networkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/5GSEC/SentryFlow/pkg/config"
	"github.com/5GSEC/SentryFlow/pkg/exporter"
	"github.com/5GSEC/SentryFlow/pkg/k8s"
	"github.com/5GSEC/SentryFlow/pkg/receiver"
	"github.com/5GSEC/SentryFlow/pkg/util"
	protobuf "github.com/5GSEC/SentryFlow/protobuf/golang"
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

func (m *Manager) areK8sReceivers(cfg *config.Config) bool {
	if len(cfg.Receivers.ServiceMeshes) > 0 {
		return true
	}

	for _, other := range cfg.Receivers.Others {
		if other.Name == util.NginxIncorporationIngressController {
			return true
		}
	}

	return false
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
		mgr.Logger.Error(err)
		return
	}

	if mgr.areK8sReceivers(cfg) {
		k8sClient, err := k8s.NewClient(registerAndGetScheme(), kubeConfig)
		if err != nil {
			mgr.Logger.Errorf("failed to create k8s client: %v", err)
			return
		}
		mgr.K8sClient = k8sClient
	}

	mgr.Wg.Add(1)
	go func() {
		defer mgr.Wg.Done()
		mgr.startHttpServer(cfg.Filters.Server.Port)
	}()

	if err := receiver.Init(mgr.Ctx, mgr.K8sClient, cfg, mgr.ApiEvents, mgr.GrpcServer, mgr.Wg); err != nil {
		mgr.Logger.Errorf("failed to initialize receiver: %v", err)
		return
	}

	if err := exporter.Init(mgr.Ctx, mgr.GrpcServer, cfg, mgr.ApiEvents, mgr.Wg); err != nil {
		mgr.Logger.Errorf("failed to initialize exporter: %v", err)
		return
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
	utilruntime.Must(networkingv1alpha3.AddToScheme(scheme))
	utilruntime.Must(corev1.AddToScheme(scheme))
	utilruntime.Must(appsv1.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))
	return scheme
}
