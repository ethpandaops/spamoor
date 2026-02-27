# spamoor
VERSION := $(shell git rev-parse --short HEAD)
BUILDTIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
GOLDFLAGS += -X 'github.com/ethpandaops/spamoor/utils.BuildVersion="$(VERSION)"'
GOLDFLAGS += -X 'github.com/ethpandaops/spamoor/utils.BuildTime="$(BUILDTIME)"'
GOLDFLAGS += -X 'github.com/ethpandaops/spamoor/utils.BuildRelease="$(RELEASE)"'

.PHONY: all docs build test clean generate-spammer-index generate-symbols plugins

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

generate-spammer-index:
	@echo "Generating spammer configuration index..."
	scripts/generate-spammer-index.sh spammer-configs

generate-symbols:
	@echo "Regenerating Yaegi symbol files for dynamic scenario loading..."
	go generate ./plugin/symbols/...

clean:
	rm -f bin/*

devnet:
	.hack/devnet/run.sh

devnet-run: devnet docs build
	bin/spamoor-daemon --rpchost-file .hack/devnet/generated-hosts.txt --privkey 3fd98b5187bf6526734efaa644ffbb4e3670d66f5d0268ce0323ec09124bff61 --port 8080 --db .hack/devnet/custom-spamoor.db --startup-delay 10

devnet-clean:
	.hack/devnet/cleanup.sh

plugins:
	@echo "Building plugin archives..."
	@mkdir -p bin/plugins
	@for dir in plugins/*/; do \
		if [ -d "$$dir" ]; then \
			plugin_name=$$(basename "$$dir"); \
			echo "Building plugin: $$plugin_name"; \
			echo "name: $$plugin_name" > "plugins/$$plugin_name/plugin.yaml"; \
			echo "build_time: $(BUILDTIME)" >> "plugins/$$plugin_name/plugin.yaml"; \
			echo "git_version: $(VERSION)" >> "plugins/$$plugin_name/plugin.yaml"; \
			tar -czf "plugins/$$plugin_name.tar.gz" -C "plugins/$$plugin_name" .; \
			rm "plugins/$$plugin_name/plugin.yaml"; \
		fi \
	done
	@echo "Plugins built in plugins/"
