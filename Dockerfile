# syntax=docker/dockerfile:1.7
#
# Multi-stage build for evm-oracle-demo-indexer-service.
#
# Stage 1 — build:
#   * pinned codegen toolchain (buf, protoc-gen-go, protoc-gen-go-grpc)
#     per architecture rule 9 — never `@latest`.
#   * `make build` runs `proto-gen` (writes into `internal/genproto/`)
#     and then compiles the static binary.
#
# Stage 2 — runtime: distroless static-debian12:nonroot, no shell, no
# package manager, no root user.

FROM golang:1.25 AS builder

ARG BUF_VERSION=v1.55.0
ARG PROTOC_GEN_GO_VERSION=v1.36.0
ARG PROTOC_GEN_GO_GRPC_VERSION=v1.5.1

RUN apt-get update \
    && apt-get install -y --no-install-recommends make ca-certificates git \
    && rm -rf /var/lib/apt/lists/*

# Pinned protobuf toolchain. Never `@latest` (architecture rule 9
# "How to apply" — pin every codegen plugin to a known version).
RUN go install github.com/bufbuild/buf/cmd/buf@${BUF_VERSION} \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@${PROTOC_GEN_GO_VERSION} \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@${PROTOC_GEN_GO_GRPC_VERSION}

ENV PATH=/go/bin:${PATH}

WORKDIR /src

# Pull in module metadata first so docker caches the dep layer across
# source-only changes.
COPY go.mod go.sum ./
RUN go mod download

COPY .git ./.git
COPY . .

# `make build` depends on `proto-gen`, so `internal/genproto/` is
# regenerated from `./protocols/` inside the image (those .pb.go
# files are gitignored — they never leave this layer).
RUN make build

# ------------------------------------------------------------------

FROM gcr.io/distroless/static-debian12:nonroot

LABEL org.opencontainers.image.title="evm-oracle-demo-indexer-service" \
      org.opencontainers.image.description="Single-chain event indexer for the EVM oracle demo. Watches PriceAggregator + OracleRegistry, persists past 5 confirmations, serves gRPC ListEvents/GetRequest/StreamEvents." \
      org.opencontainers.image.source="https://github.com/asolovov/evm-oracle-demo-indexer-service" \
      org.opencontainers.image.licenses="MIT"

COPY --from=builder /src/evm-oracle-demo-indexer-service /usr/local/bin/indexer-service
COPY --from=builder /src/migrations /migrations

EXPOSE 9090 8080
USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/indexer-service"]
CMD ["serve"]
