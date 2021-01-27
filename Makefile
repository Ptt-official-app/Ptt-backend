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
	@echo "output binary file into ./bin"
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
