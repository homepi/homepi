export ROOT=$(realpath $(dir $(firstword $(MAKEFILE_LIST))))
export BIN=$(ROOT)/bin
export GOBIN?=$(BIN)
export GO=$(shell which go)
export BUILD=cd $(ROOT) && $(GO) install -v -ldflags "-s"
export CGO_ENABLED=1

# Linter configurations
export LINTER=$(GOBIN)/golangci-lint
export LINTERCMD=run --no-config -v \
	--print-linter-name \
	--skip-files ".*.gen.go" \
	--skip-files ".*_test.go" \
	--sort-results \
	--disable-all \
	--enable=structcheck \
	--enable=deadcode \
	--enable=gocyclo \
	--enable=ineffassign \
	--enable=revive \
	--enable=goimports \
	--enable=errcheck \
	--enable=varcheck \
	--enable=goconst \
	--enable=megacheck \
	--enable=misspell \
	--enable=unused \
	--enable=typecheck \
	--enable=staticcheck \
	--enable=govet \
	--enable=gosimple

# Build the project
all:
	$(BUILD) ./...

.PHONY: server
server: all
	$(ROOT)/bin/homepi server --config-file=$(ROOT)/config.yaml

# lint runs vet plus a number of other checkers, it is more comprehensive, but louder
.PHONY: lint
lint:
	@LINTER_BIN=$$(command -v $(LINTER)) || { echo "golangci-lint command not found! Installing..." && $(MAKE) install-metalinter; };
	@$(GO) list -f '{{.Dir}}' ./src/... | grep -v /vendor/ \
		| xargs $(LINTER) $(LINTERCMD) ./...; if [ $$? -eq 1 ]; then \
			echo ""; \
			echo "Lint found suspicious constructs. Please check the reported constructs"; \
			echo "and fix them if necessary before submitting the code for reviewal."; \
		fi

.PHONY: gox
gox:
	@gox -output="dist/homepi_{{.OS}}_{{.Arch}}"

# for ci jobs, runs lint against the changed packages in the commit
.PHONY: ci-lint
ci-lint:
	$(shell which golangci-lint) $(LINTERCMD) --deadline 10m ./...

# Check if golangci-lint not exists, then install it
.PHONY: install-metalinter
install-metalinter:
	$(GO) get -v github.com/golangci/golangci-lint/cmd/golangci-lint@v1.41.1
	$(GO) install -v github.com/golangci/golangci-lint/cmd/golangci-lint@v1.41.1

# Run tests
.PHONY: test
test:
	@$(GO) test ./src/... -v -race

install-gox:
	@$(GO) get -v github.com/mitchellh/gox@v1.0.1

.PHONY: build-linux
build-linux: install-gox
	@$(GOX) --arch=amd64 --os=linux --output="dist/homepi_{{.OS}}_{{.Arch}}"

.PHONY: build-windows
build-windows: install-gox
	@$(GOX) --arch=amd64 --os=windows --output="dist/homepi_{{.OS}}_{{.Arch}}"

.PHONY: build-macOS
build-macOS: install-gox
	@$(GOX) --arch=amd64 --os=darwin --output="dist/homepi_{{.OS}}_{{.Arch}}"

.PHONY: build-artifacts
build-artifacts:
	@$(MAKE) build-linux && $(MAKE) build-windows && $(MAKE) build-macOS
