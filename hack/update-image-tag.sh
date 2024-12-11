#!/usr/bin/env bash
# SPDX-License-Identifier: Apache-2.0
# Copyright 2024 Authors of SentryFlow

if ! command -v yq > /dev/null; then
  echo "Installing yq..."
  go install github.com/mikefarah/yq/v4@latest
fi

TAG=$1
DEPLOYMENT_ROOT_DIR="deployments/sentryflow"

echo "Updating tag to ${TAG}"
yq -i ".image.tag = \"$TAG\"" "${DEPLOYMENT_ROOT_DIR}/values.yaml"
