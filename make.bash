#!/bin/bash

function help() {
    echo "---- Project: Ptt-backend ----"
	  echo " Usage: ./make.bash COMMAND"
    echo
    echo " Management Commands:"
    echo "  build              Build project"
    echo "  deps               Ensures fresh go.mod and go.sum for dependencies"
    echo "  format             Formats Go code"
    echo "  lint               Run golangci-lint check"
    echo "  test-unit          Run all unit tests"
    echo "  test-integration   Run all integration and unit tests"
    echo
}

function build(){
  VERSION=$(git describe --tags $(git rev-list --tags --max-count=1) 2>/dev/null)
  BUILDTIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
  GITSHA=$(git rev-parse --short HEAD 2>/dev/null)
  if [ -z "$VERSION" ]; then \
    VERSION="git-$GITSHA"
  fi

  GOFLAGS="-trimpath"
  LDFLAGS="-X main/version.version=$VERSION -X main/version.commit=$GITSHA -X main/version.buildTime=$BUILDTIME"
  mkdir -p "bin"
  echo "binary file output into ./bin"
  go build "$GOFLAGS" -ldflags "$LDFLAGS" -o ./bin ./...
}

function format() {
    files=$(find . -path ./vendor -prune -o -name "*.go" -print)
	  gofmt -s -w $files
}

function lint() {
    GOBIN=$(go env GOBIN)
	  if [ ! -f "$GOBIN/golangci-lint" ]; then \
		  curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$GOBIN" "$GOLANGCI_LINT_VERSION"; \
		  echo "download golangci-lint into $GOBIN" ;\
   	fi;
   	go vet ./...
	  echo "golangci-lint checking..."
	  "$GOBIN"/golangci-lint run --deadline=30m --enable=misspell --enable=gosec --enable=gofmt --enable=goimports --enable=golint ./cmd/... ./...
}
# no arguments
if [ $# -lt 1 ];then
    help
    exit 0
# number of arguments greater than 1
elif [ $# -gt 1 ];then
    echo "invalid args, please check command"
    help
    exit 0
fi

case "$1" in
help)
    help
    ;;
# build: Build project
build)
    build
    ;;
# deps: Ensures fresh go.mod and go.sum for dependencies
deps)
	  go mod tidy
	  go mod verify
    ;;
# format: Formats Go code
format)
    format
    ;;
# lint: Run golangci-lint check
lint)
    lint
    ;;
# test-unit: Run all unit tests
test-unit)
	  go test -v -cover . ./...
    ;;
# test-integration: Run all integration and unit tests
test-integration)
    echo 'mode: atomic' > coverage.out
    go list ./... | xargs -n1 -I{} sh -c 'go test -race -tags=integration -covermode=atomic -coverprofile=coverage.tmp -coverpkg $(go list ./... | tr "\n" ",") {} && tail -n +2 coverage.tmp >> coverage.out || exit 255'
    rm coverage.tmp
    ;;
*)
    echo "invalid args, please check command"
    help
esac

