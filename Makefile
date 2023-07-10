#@IgnoreInspection BashAddShebang

export ROOT=$(realpath $(dir $(lastword $(MAKEFILE_LIST))))

export APP=webhook-proxy

export BUILD_INFO_PKG="github.com/snapp-incubator/jira-element-proxy/pkg/version"

export LDFLAGS="-w -s -X '$(BUILD_INFO_PKG).Date=$$(date)' -X '$(BUILD_INFO_PKG).BuildVersion=$$(git rev-parse HEAD | cut -c 1-8)' -X '$(BUILD_INFO_PKG).VCSRef=$$(git rev-parse --abbrev-ref HEAD)'"

all: format lint build

run-version:
	go run -ldflags $(LDFLAGS) ./cmd/webhook-proxy version

run-api:
	go run -ldflags $(LDFLAGS) ./cmd/webhook-proxy api

build:
	go build -ldflags $(LDFLAGS) ./cmd/webhook-proxy

install:
	go install -ldflags $(LDFLAGS) ./cmd/webhook-proxy

check-formatter:
	which goimports || GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports

format: check-formatter
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R goimports -w R
	find $(ROOT) -type f -name "*.go" -not -path "$(ROOT)/vendor/*" | xargs -n 1 -I R gofmt -s -w R

check-linter:
	which golangci-lint || GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint@v1.23.8

lint: check-linter
	golangci-lint run $(ROOT)/...

install-hook:
	git config --local core.hooksPath ./githooks

test:
	go test -ldflags $(LDFLAGS) -v -race -p 1 `go list ./... | grep -v integration`

ci-test:
	go test -ldflags $(LDFLAGS) -v -race -p 1 -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -func coverage.txt

integration-tests:
	go test -ldflags $(LDFLAGS) -v -race -p 1 `go list ./... | grep integration`

up:
	docker-compose up -d

down:
	docker-compose down
