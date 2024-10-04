# SPDX-License-Identifier: Apache-2.0
# Copyright 2024 Authors of SentryFlow

PROTOS:="sentryflow.proto sentryflow_metrics.proto"

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

.PHONY: all
all: golang python ## Generate code for Go and Python

.PHONY: golang
golang: go-plugins ## Generate code for Go
	@protoc --go_out=golang --go_opt=paths=source_relative --go-grpc_out=golang --go-grpc_opt=paths=source_relative "${PROTOS}"
	@cd golang; go mod tidy

.PHONY: python
python: ## Generate code for Python
	@protoc -I=. --python_out=pyi_out:python "${PROTOS}"

.PHONY: clean
clean: ## Remove the generated code
	@cd golang; rm -f *.pb.go; go mod tidy
	@cd python; rm -f *.py*

.PHONY:
go-plugins: ## Download Go plugins locally if necessary.
	@test -s protoc-gen-go || go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@test -s protoc-gen-go-grpc || go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
