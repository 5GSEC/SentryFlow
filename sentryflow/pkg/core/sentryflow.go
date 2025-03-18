// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package core

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"istio.io/client-go/pkg/apis/extensions/v1alpha1"
	networkingv1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	protobuf "github.com/5GSEC/SentryFlow/protobuf/golang"
	"github.com/5GSEC/SentryFlow/sentryflow/pkg/config"
	"github.com/5GSEC/SentryFlow/sentryflow/pkg/exporter"
	"github.com/5GSEC/SentryFlow/sentryflow/pkg/k8s"
	"github.com/5GSEC/SentryFlow/sentryflow/pkg/receiver"
	"github.com/5GSEC/SentryFlow/sentryflow/pkg/util"
)

type Manager struct {
	Ctx                 context.Context
	Logger              *zap.SugaredLogger
	GrpcServer          *grpc.Server
	HttpServer          *http.Server
	K8sClient           client.Client
	Wg                  *sync.WaitGroup
	ApiEvents           chan *protobuf.APIEvent
	configChan          chan *config.Config
	receiversCtx        context.Context
	receiversCancelFunc context.CancelFunc
	receiversLock       *sync.Mutex
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

func (m *Manager) run(cfg *config.Config, kubeConfig string) {
	m.Ctx, _ = m.setupSignalHandler(make(chan os.Signal, 2))
	m.GrpcServer = grpc.NewServer()
	m.Wg = &sync.WaitGroup{}
	m.ApiEvents = make(chan *protobuf.APIEvent, 10240)

	if m.areK8sReceivers(cfg) {
		k8sClient, err := k8s.NewClient(registerAndGetScheme(), kubeConfig)
		if err != nil {
			m.Logger.Errorf("failed to create k8s client: %v", err)
			return
		}
		m.K8sClient = k8sClient
	}

	m.Wg.Add(1)
	go func() {
		defer m.Wg.Done()
		m.startHttpServer(cfg.Filters.Server.Port)
	}()

	m.receiversCtx, m.receiversCancelFunc = m.setupSignalHandler(make(chan os.Signal, 2))
	if err := receiver.Init(m.receiversCtx, m.K8sClient, cfg, m.Wg, m.receiversLock); err != nil {
		m.Logger.Errorf("failed to initialize receiver: %v", err)
		return
	}

	if err := exporter.Init(m.Ctx, m.GrpcServer, cfg, m.ApiEvents, m.Wg); err != nil {
		m.Logger.Errorf("failed to initialize exporter: %v", err)
		return
	}

	m.Wg.Add(1)
	go func() {
		defer m.Wg.Done()
		m.startGrpcServer(cfg.Exporter.Grpc.Port)
	}()

	m.Logger.Info("Started SentryFlow")

	for {
		select {
		case <-m.Ctx.Done():
			m.Logger.Info("Shutdown Signal Received. Waiting for all workers to finish.")
			m.Logger.Info("Shutting down SentryFlow")
			m.receiversCancelFunc()
			m.stopServers()
			m.Wg.Wait()
			close(m.ApiEvents)
			close(m.configChan)
			m.Logger.Info("All workers finished. Stopped SentryFlow")
			return

		case updatedConfig := <-m.configChan:
			m.receiversCancelFunc()
			if m.areK8sReceivers(updatedConfig) {
				k8sClient, err := k8s.NewClient(registerAndGetScheme(), kubeConfig)
				if err != nil {
					m.Logger.Errorf("failed to create k8s client: %v", err)
					return
				}
				m.K8sClient = k8sClient
			}
			m.receiversCtx, m.receiversCancelFunc = m.setupSignalHandler(make(chan os.Signal, 2))
			if err := receiver.Init(m.receiversCtx, m.K8sClient, updatedConfig, m.Wg, m.receiversLock); err != nil {
				m.Logger.Errorf("failed to initialize receiver: %v", err)
				return
			}
		}
	}
}

func registerAndGetScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()
	utilruntime.Must(networkingv1alpha3.AddToScheme(scheme))
	utilruntime.Must(corev1.AddToScheme(scheme))
	utilruntime.Must(appsv1.AddToScheme(scheme))
	utilruntime.Must(v1alpha1.AddToScheme(scheme))
	return scheme
}

func (m *Manager) watchConfig(configFilePath string, logger *zap.SugaredLogger) {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		cfg, err := config.New(configFilePath, logger)
		if err != nil {
			m.Logger.Errorf("failed to reload config, %v", err)
			return
		}
		m.configChan <- cfg
		m.Logger.Info("config file changed, reloading config...")
	})
}

func (m *Manager) setupSignalHandler(c chan os.Signal) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, util.LoggerContextKey{}, m.Logger)

	shutdownSignals := []os.Signal{os.Interrupt, syscall.SIGTERM}
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		cancel()
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()
	return ctx, cancel
}

func Run(configFilePath string, kubeConfig string, logger *zap.SugaredLogger) {
	mgr := &Manager{
		Logger:        logger,
		configChan:    make(chan *config.Config),
		receiversLock: &sync.Mutex{},
	}
	mgr.Logger.Info("Starting SentryFlow")

	cfg, err := config.New(configFilePath, mgr.Logger)
	if err != nil {
		mgr.Logger.Error(err)
		return
	}
	mgr.watchConfig(configFilePath, logger)

	mgr.run(cfg, kubeConfig)
}
