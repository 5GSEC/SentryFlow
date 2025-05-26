// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package apievent

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"

	pb "github.com/5GSEC/SentryFlow/protobuf/golang"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"

	"github.com/5GSEC/SentryFlow/sfctl/pkg/client"
)

func startEventsStreaming(ctx context.Context, config string, k8sClientset kubernetes.Interface) error {
	logger.Debug("starting port-forwarding")
	localPort, err := getFreeLocalPort()
	if err != nil {
		return err
	}
	logger.Debug("local port: ", localPort)
	if err := setupPortForwarding(ctx, config, k8sClientset, localPort); err != nil {
		return err
	}
	logger.Debug("started port-forwarding")

	return startStreaming(ctx, localPort)
}

func setupPortForwarding(ctx context.Context, config string, k8sClientset kubernetes.Interface, localPort int64) error {
	podName, err := getPodName(ctx, k8sClientset)
	if err != nil {
		return err
	}
	logger.Debug("pod name: ", podName)

	return startPortForwarding(config, podName, sentryflowNamespace, sentryflowPort, localPort)
}

func startPortForwarding(config, podName, namespace, sentryFlowPort string, localPort int64) error {
	restCfg, err := client.GetConfig(config, nil)
	if err != nil {
		return err
	}
	roundTripper, upgrader, err := spdy.RoundTripperFor(restCfg)
	if err != nil {
		return fmt.Errorf("unable to create round tripper and upgrader, error=%s", err.Error())
	}

	serverURL, err := url.Parse(restCfg.Host)
	if err != nil {
		return fmt.Errorf("failed to parse apiserver URL from kubeconfig. error=%s", err.Error())
	}
	serverURL.Path = fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", namespace, podName)

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: roundTripper}, http.MethodPost, serverURL)

	StopChan, readyChan := make(chan struct{}, 1), make(chan struct{}, 1)
	out, errOut := new(bytes.Buffer), new(bytes.Buffer)

	forwarder, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%v", localPort, sentryFlowPort)},
		StopChan, readyChan, out, errOut)
	if err != nil {
		return fmt.Errorf("unable to portforward. error=%s", err.Error())
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- forwarder.ForwardPorts()
	}()

	select {
	case err = <-errChan:
		close(errChan)
		forwarder.Close()
		return fmt.Errorf("could not create port forward %s", err)
	case <-readyChan:
		return nil
	}
}

func getFreeLocalPort() (int64, error) {
	for attempts := 0; attempts < 100; attempts++ {
		port, err := getRandomPort()
		if err != nil {
			continue
		}

		listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", strconv.FormatInt(port, 10)))
		if err != nil {
			return -1, err
		}
		if err := listener.Close(); err != nil {
			return -1, err
		}
		return port, nil
	}

	return -1, fmt.Errorf("failed to find a free local port")
}

func getRandomPort() (int64, error) {
	// Return a port number > 32767
	n, err := rand.Int(rand.Reader, big.NewInt(32900-32768))
	if err != nil {
		return -1, fmt.Errorf("failed to generate random integer for port")
	}

	portNo := n.Int64() + 32768
	return portNo, nil
}

func getPodName(ctx context.Context, k8sClientset kubernetes.Interface) (string, error) {
	logger.Debug("getting pod name")

	if k8sClientset == nil {
		logger.Warn("k8sClientset is nil")
		return "", errors.New("k8sClientset is nil")
	}
	pods, err := k8sClientset.CoreV1().Pods(sentryflowNamespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return "", err
	}
	if len(pods.Items) == 0 {
		return "", fmt.Errorf("no sentryflow pods found in %s namespace", sentryflowNamespace)
	}
	return pods.Items[0].Name, nil
}

func startStreaming(ctx context.Context, localPort int64) error {
	logger.Debug("creating SentryFlow client")

	sfClient, closer := client.NewSentryFlowClient(localPort)
	defer closer()

	hostname, err := os.Hostname()
	if err != nil {
		logger.Warnf("failed to get hostname: %v", err)
	}

	ips, err := net.LookupIP(hostname)
	if err != nil {
		logger.Warnf("failed to get current host IP address: %v", err)
	}

	clientInfo := &pb.ClientInfo{
		HostName:  hostname,
		IPAddress: ips[0].String(),
	}

	logger.Info("starting API Events streaming")
	stream, err := sfClient.GetAPIEvent(ctx, clientInfo)
	if err != nil {
		return err
	}

	logger.Info("started API Events streaming")
	for {
		event, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			if status.Code(err) != codes.Canceled {
				return err
			}
		}

		select {
		case <-ctx.Done():
			logger.Info("Shutting down API Events streaming")
			return nil
		default:
			if statusCode != "" {
				printFilteredEvents(event, prettyPrint)
				continue
			}

			if prettyPrint {
				body, err := json.MarshalIndent(event, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(body))
			} else {
				body, err := json.Marshal(event)
				if err != nil {
					return err
				}
				fmt.Println(string(body))
			}
		}
	}
}
