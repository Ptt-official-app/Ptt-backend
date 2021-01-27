FILES_TO_FMT      ?= $(shell find . -path ./vendor -prune -o -name "*.go" -print)

GOFLAGS   :=
LDFLAGS   :=

BUILDTIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GITSHA := $(shell git rev-parse --short HEAD 2>/dev/null)

ifndef VERSION
	VERSION := git-$(GITSHA)
endif

GOFLAGS += -trimpath

LDFLAGS += -X $(PKG)/version.version=$(VERSION)
LDFLAGS += -X $(PKG)/version.commit=$(GITSHA)
LDFLAGS += -X $(PKG)/version.buildTime=$(BUILDTIME)

## help: Show makefile commands
.PHONY: help
help: Makefile
	@echo "---- Project: Ptt-backend ----"
	@echo " Usage: make COMMAND"
	@echo
	@echo " Management Commands:"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## build: Build project
.PHONY: build
build:
	@mkdir -p "bin"
	@echo "binary file output into ./bin"
	@go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o ./bin ./...

## deps: Ensures fresh go.mod and go.sum for dependencies
.PHONY: deps
deps:
	@go mod tidy
	@go mod verify

## format: Formats Go code
.PHONY: format
format:
	@echo ">> formatting code"
	@gofmt -s -w $(FILES_TO_FMT)

## lint: Run golangci-lint check
.PHONY: lint
lint:
	@if [ ! -f $(GOBIN)/golangci-lint ]; then \
		curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) $(GOLANGCI_LINT_VERSION); \
		echo "download golangci-lint into $(GOBIN)" ;\
	fi;
	@echo "golangci-lint checking..."
	@$(GOBIN)/golangci-lint run --deadline=30m --enable=misspell --enable=gosec --enable=gofmt --enable=goimports --enable=golint ./cmd/... ./...
	@go vet ./...

## test-unit: Run all unit tests
.PHONY: test-unit
test-unit:
	@go test -v -cover . ./...

## test-integration: Run all integration and unit tests
.PHONY: test-integration
test-integration:
	echo 'mode: atomic' > coverage.out
	go list ./... | xargs -n1 -I{} sh -c 'go test -race -tags=integration -covermode=atomic -coverprofile=coverage.tmp -coverpkg $(go list ./... | tr "\n" ",") {} && tail -n +2 coverage.tmp >> coverage.out || exit 255'
	rm coverage.tmp
