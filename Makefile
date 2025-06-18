# spamoor
VERSION := $(shell git rev-parse --short HEAD)
BUILDTIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GOLDFLAGS += -X 'github.com/ethpandaops/spamoor/utils.BuildVersion="$(VERSION)"'
GOLDFLAGS += -X 'github.com/ethpandaops/spamoor/utils.BuildTime="$(BUILDTIME)"'
GOLDFLAGS += -X 'github.com/ethpandaops/spamoor/utils.BuildRelease="$(RELEASE)"'

.PHONY: all docs build test clean

all: docs build

test:
	go test ./...

build:
	@echo version: $(VERSION)
	env CGO_ENABLED=1 CGO_CFLAGS="-O2 -D__BLST_PORTABLE__" go build -v -tags=with_blob_v1,ckzg -o bin/ -ldflags="-s -w $(GOLDFLAGS)" ./cmd/*

build-lib:
	@echo version: $(VERSION)
	cat go.mod | sed -E "s/^replace/\/\/replace/" > go.lib.mod
	go mod tidy -modfile=go.lib.mod
	env CGO_ENABLED=1 go build -modfile=go.lib.mod -v -o bin/ -ldflags="-s -w $(GOLDFLAGS)" ./cmd/*

docs:
	go install github.com/swaggo/swag/cmd/swag@v1.16.3 && swag init -g handler.go -d webui/handlers/api --parseDependency -o webui/handlers/docs

clean:
	rm -f bin/*

devnet:
	.hack/devnet/run.sh

devnet-run: devnet build
	bin/spamoor-daemon --rpchost-file .hack/devnet/generated-hosts.txt --privkey 3fd98b5187bf6526734efaa644ffbb4e3670d66f5d0268ce0323ec09124bff61 --port 8080 --db .hack/devnet/custom-spamoor.db

devnet-clean:
	.hack/devnet/cleanup.sh
