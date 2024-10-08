# SPDX-License-Identifier: Apache-2.0
# Copyright 2024 Authors of SentryFlow

REGISTRY ?= docker.io/5gsec
DOCKER_IMAGE ?= $(REGISTRY)/sentryflow-httpfilter
DOCKER_TAG ?= latest
CONTAINER_TOOL ?= docker

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

.PHONY: toolchain
toolchain: ## Install Rust WASM toolchain
	@test rustup || curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
	@rustup target add wasm32-unknown-unknown

.PHONY: build
build: ## Build plugin.
	@cargo build --target wasm32-unknown-unknown --release

.PHONY: clean
clean: ## Remove generated stuff.
	@cargo clean

.PHONY: image
image: ## Build Plugin's container image
	$(CONTAINER_TOOL) build -t ${DOCKER_IMAGE}:${DOCKER_TAG} .

.PHONY: push
push: ## Push Plugin's container image
	$(CONTAINER_TOOL) push ${DOCKER_IMAGE}:${DOCKER_TAG}

.PHONY: imagex
imagex: ## Build and push Plugin's multi-platform container image.
	$(CONTAINER_TOOL) buildx build --push --platform=linux/arm64,linux/amd64 -t ${DOCKER_IMAGE}:${DOCKER_TAG} .
