# SPDX-License-Identifier: Apache-2.0
# Copyright 2024 Authors of SentryFlow

BINARY_NAME ?= sfctl
VERSION ?= $(shell git rev-parse HEAD)
BUILD_TS ?= $(shell date)

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

##@ Development
.PHONY: fmt
fmt: ## Run go fmt against code
	@go fmt ./...

.PHONY: vet
vet: ## Run go vet against code
	@go vet ./...

GOLANGCI_LINT = $(shell pwd)/bin/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.63.0
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
	@../hack/add-license-header

##@ Build

.PHONY: build
build: fmt vet ## Build sfctl binary
	@go mod tidy
	@CGO_ENABLED=0 go build -ldflags="-s" -o bin/"${BINARY_NAME}" .
