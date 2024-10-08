# SPDX-License-Identifier: Apache-2.0
# Copyright 2024 Authors of SentryFlow

FROM rust:1.81.0 AS builder

WORKDIR /envoy-plugin

COPY . .

RUN make toolchain build

FROM scratch

COPY --from=builder /envoy-plugin/target/wasm32-unknown-unknown/release/httpfilters.wasm ./plugin.wasm
