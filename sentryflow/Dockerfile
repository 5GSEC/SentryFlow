FROM golang:1.23 AS builder

ARG TARGETOS
ARG TARGETARCH

RUN mkdir -p /protobuf/golang

COPY protobuf/golang /protobuf/golang

WORKDIR /sentryflow

COPY sentryflow/go.mod .
COPY sentryflow/go.sum .
RUN go mod download

COPY sentryflow/cmd cmd/
COPY sentryflow/pkg pkg/
COPY sentryflow/main.go main.go
COPY sentryflow/Makefile Makefile

RUN GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} make build

FROM gcr.io/distroless/static-debian12

COPY --from=builder /sentryflow/bin/sentryflow /

ENTRYPOINT ["/sentryflow"]
