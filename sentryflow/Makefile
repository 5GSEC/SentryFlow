# SPDX-License-Identifier: Apache-2.0

IMAGE_NAME = 5gsec/sentryflow
TAG = v0.1

# CONTAINER_TOOL defines the container tool to be used for building images.
# Be aware that the target commands are only tested with Docker which is
# scaffolded by default. However, you might want to replace it to use other
# tools. (i.e. podman)
CONTAINER_TOOL ?= docker

.PHONY: image
image:
	$(CONTAINER_TOOL) build -t $(IMAGE_NAME):$(TAG) -t $(IMAGE_NAME):latest -f ./Dockerfile ../

.PHONY: push-image
push-image:
	$(CONTAINER_TOOL) push $(IMAGE_NAME):$(TAG)
	$(CONTAINER_TOOL) push $(IMAGE_NAME):latest

.PHONY: rm-image
rm-image:
	docker rmi $(IMAGE_NAME):$(TAG)

.PHONY: run
run:
	docker run -it --rm $(IMAGE_NAME):$(TAG)

PLATFORMS ?= linux/arm64,linux/amd64
.PHONY: docker-buildx
docker-buildx:
	# copy existing Dockerfile and insert --platform=${BUILDPLATFORM} into Dockerfile.cross, and preserve the original Dockerfile
	sed -e '1 s/\(^FROM\)/FROM --platform=\$$\{BUILDPLATFORM\}/; t' -e ' 1,// s//FROM --platform=\$$\{BUILDPLATFORM\}/' Dockerfile > Dockerfile.cross
	- $(CONTAINER_TOOL) buildx create --name builder
	$(CONTAINER_TOOL) buildx use builder
	- $(CONTAINER_TOOL) buildx build --push --platform=$(PLATFORMS) --build-arg VERSION=$(TAG) --tag $(IMAGE_NAME):$(TAG) -f Dockerfile.cross ../ || { $(CONTAINER_TOOL) buildx rm builder; rm Dockerfile.cross; exit 1; }
	- $(CONTAINER_TOOL) buildx rm builder
	rm Dockerfile.cross

.PHONY: gofmt
gofmt:
	cd $(CURDIR); gofmt -w -s -d $(shell find . -type f -name '*.go' -print)

.PHONY: golint
golint:
ifeq (, $(shell which golint))
	@{ \
	set -e ;\
	GOLINT_TEMP_DIR=$$(mktemp -d) ;\
	cd $$GOLINT_TEMP_DIR ;\
	go mod init tmp ;\
	go get golang.org/x/lint/golint ;\
	go install golang.org/x/lint/golint ;\
	rm -rf $$GOLINT_TEMP_DIR ;\
	}
endif
	cd $(CURDIR); golint ./...

.PHONY: gosec
gosec:
ifeq (, $(shell which gosec))
	@{ \
	set -e ;\
	GOSEC_TEMP_DIR=$$(mktemp -d) ;\
	cd $$GOSEC_TEMP_DIR ;\
	go mod init tmp ;\
	go get github.com/securego/gosec/v2/cmd/gosec ;\
	go install github.com/securego/gosec/v2/cmd/gosec ;\
	rm -rf $$GOSEC_TEMP_DIR ;\
	}
endif
	cd $(CURDIR); gosec -exclude=G402 ./...
