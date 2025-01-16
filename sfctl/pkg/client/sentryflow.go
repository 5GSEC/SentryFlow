// SPDX-License-Identifier: Apache-2.0
// Copyright 2024 Authors of SentryFlow

package client

import (
	"fmt"

	pb "github.com/5GSEC/SentryFlow/protobuf/golang"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/5GSEC/SentryFlow/sfctl/pkg/util"
)

func NewSentryFlowClient(port int64) (pb.SentryFlowClient, func()) {
	logger := util.GetLogger()

	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", port),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)

	if err != nil {
		logger.Errorf("failed to connect to SentryFlow: %v", err)
		return nil, nil
	}

	return pb.NewSentryFlowClient(conn), func() {
		if err := conn.Close(); err != nil {
			logger.Warnf("failed to close connection to SentryFlow: %v", err)
		}
	}
}
