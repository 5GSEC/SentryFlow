# SPDX-License-Identifier: Apache-2.0
# Copyright 2024 Authors of SentryFlow

BINARY_NAME ?= sentryflow
REGISTRY ?= docker.io/5gsec
VERSION ?= $(shell git rev-parse HEAD)
BUILD_TS ?= $(shell date)
DOCKER_IMAGE ?= $(REGISTRY)/$(BINARY_NAME)
DOCKER_TAG ?= latest
CONTAINER_TOOL ?= docker

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

##@ Development
.PHONY: run
run: fmt vet ## Run SentryFlow on your host
	@go mod tidy;go run main.go

.PHONY: fmt
fmt: ## Run go fmt against code
	@go fmt ./...

.PHONY: vet
vet: ## Run go vet against code
	@go vet ./...

GOLANGCI_LINT = $(shell pwd)/bin/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.60.3
golangci-lint:
	@[ -f $(GOLANGCI_LINT) ] || { \
	set -e ;\
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell dirname $(GOLANGCI_LINT)) $(GOLANGCI_LINT_VERSION) ;\
	}

.PHONY: lint
lint: golangci-lint ## Run golangci-lint linter
	@$(GOLANGCI_LINT) run

.PHONY: license
license: ## Check and fix license header on all go files
	@../scripts/add-license-header

##@ Build

.PHONY: build
build: fmt vet ## Build SentryFlow binary
	@CGO_ENABLED=0 go build -ldflags="-s \
	-X 'main.Version=${VERSION}' \
	-X 'main.CommitHash=${VERSION}' \
	-X 'main.BuildTimestamp=${BUILD_TS}'" \
 	-o bin/"${BINARY_NAME}" main.go

.PHONY: image
image: ## Build SentryFlow's container image
	$(CONTAINER_TOOL) build -t ${DOCKER_IMAGE}:${DOCKER_TAG} -f Dockerfile ../

.PHONY: push
push: ## Push SentryFlow's container image
	$(CONTAINER_TOOL) push ${DOCKER_IMAGE}:${DOCKER_TAG}

PLATFORMS ?= linux/arm64,linux/amd64
.PHONY: imagex
imagex: ## Build and push SentryFlow's container image for cross-platform support
	# copy existing Dockerfile and insert --platform=${BUILDPLATFORM} into Dockerfile.cross, and preserve the original Dockerfile
	sed -e '1 s/\(^FROM\)/FROM --platform=\$$\{BUILDPLATFORM\}/; t' -e ' 1,// s//FROM --platform=\$$\{BUILDPLATFORM\}/' Dockerfile > Dockerfile.cross
	- $(CONTAINER_TOOL) buildx create --name project-v3-builder
	$(CONTAINER_TOOL) buildx use project-v3-builder
	- $(CONTAINER_TOOL) buildx build --push --platform=$(PLATFORMS) --tag ${DOCKER_IMAGE}:${DOCKER_TAG} -f Dockerfile.cross ../ || { $(CONTAINER_TOOL) buildx rm project-v3-builder; rm Dockerfile.cross; exit 1; }
	- $(CONTAINER_TOOL) buildx rm project-v3-builder
	rm Dockerfile.cross
