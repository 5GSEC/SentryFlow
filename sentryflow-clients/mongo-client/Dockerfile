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
COPY /sentryflow-clients/mongo-client .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} go build -ldflags="-w" -o mongo-client main.go

### Make executable image
FROM alpine:3.19 as client

RUN echo "@community http://dl-cdn.alpinelinux.org/alpine/edge/community" | tee -a /etc/apk/repositories

RUN apk --no-cache update
RUN apk add bash

COPY --from=builder /app/mongo-client /

CMD ["/mongo-client"]
