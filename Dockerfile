# build env
FROM golang:1.24 AS build-env
COPY go.mod go.sum /src/
WORKDIR /src
RUN go mod download
COPY . .
ARG TARGETOS
ARG TARGETARCH
ARG release=
RUN <<EOR
  VERSION=$(git rev-parse --short HEAD)
  BUILDTIME=$(date -u '+%Y-%m-%dT%H:%M:%SZ')
  RELEASE=$release
  CGO_ENABLED=1 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -tags=with_blob_v1 -o /app/ -ldflags="-s -w -X 'github.com/ethpandaops/spamoor/utils.BuildVersion=${VERSION}' -X 'github.com/ethpandaops/spamoor/utils.BuildRelease=${RELEASE}' -X 'github.com/ethpandaops/spamoor/utils.Buildtime=${BUILDTIME}'" ./cmd/*
EOR

# final stage
FROM debian:stable-slim
WORKDIR /app
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
RUN update-ca-certificates
ENV PATH="$PATH:/app"
COPY --from=build-env /app/* /app
ENTRYPOINT ["./spamoor"]
