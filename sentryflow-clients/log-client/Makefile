# SPDX-License-Identifier: Apache-2.0

IMAGE_NAME = 5gsec/sentryflow-log-client
TAG = v0.1

# CONTAINER_TOOL defines the container tool to be used for building images.
# Be aware that the target commands are only tested with Docker which is
# scaffolded by default. However, you might want to replace it to use other
# tools. (i.e. podman)
CONTAINER_TOOL ?= docker

.PHONY: image
image:
	$(CONTAINER_TOOL) build -t $(IMAGE_NAME):$(TAG) -t $(IMAGE_NAME):latest -f ./Dockerfile ../../

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
	- $(CONTAINER_TOOL) buildx build --push --platform=$(PLATFORMS) --build-arg VERSION=$(TAG) --tag $(IMAGE_NAME):$(TAG) -f Dockerfile.cross ../../ || { $(CONTAINER_TOOL) buildx rm builder; rm Dockerfile.cross; exit 1; }
	- $(CONTAINER_TOOL) buildx rm builder
	rm Dockerfile.cross
