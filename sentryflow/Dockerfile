# SPDX-License-Identifier: Apache-2.0

FROM golang:1.21-alpine3.19 as builder
ARG TARGETOS
ARG TARGETARCH

RUN apk --no-cache update
RUN apk add --no-cache make protobuf
RUN go install github.com/golang/protobuf/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

WORKDIR /protobuf
COPY /protobuf .
RUN go mod tidy
RUN make build

WORKDIR /app
COPY /sentryflow .

# Build
# the GOARCH has not a default value to allow the binary be built according to the host where the command
# was called. For example, if we call make docker-build in a local env which has the Apple Silicon M1 SO
# the docker BUILDPLATFORM arg will be linux/arm64 when for Apple x86 it will be linux/amd64. Therefore,
# by leaving it empty we can ensure that the container and binary shipped on it will have the same platform.
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -ldflags="-w" -o sentryflow main.go

### Make executable image
FROM alpine:3.19 as sentryflow

RUN echo "@community http://dl-cdn.alpinelinux.org/alpine/edge/community" | tee -a /etc/apk/repositories

RUN apk --no-cache update
RUN apk add bash

COPY --from=builder /app/sentryflow .

CMD ["./sentryflow"]