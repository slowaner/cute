ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

LINT_IMAGE ?= golangci/golangci-lint:v2.6.2-alpine

ifeq ($(shell uname),Darwin)
  GO_BUILD_CACHE ?= $(HOME)/Library/Caches/go-build
else
  GO_BUILD_CACHE ?= $(HOME)/.cache/go-build
endif
GO_MOD_CACHE ?= $(HOME)/go/pkg/mod
GOLANGCI_LINT_CACHE ?= $(HOME)/.cache/golangci-lint

export GO111MODULE=on
export GOSUMDB=off
LOCAL_BIN:=$(CURDIR)/bin

.PHONY: install
install:
	go mod tidy && go mod download

# run full lint like in pipeline
.PHONY: lint
lint: ### Run linter using Docker.
	@echo app version $(BUILD_VERSION)
	@echo run $(LINT_IMAGE)
	@docker run -it --rm \
		--platform linux/amd64 \
		-v $(GO_BUILD_CACHE):/.cache/go-build -e GOCACHE=/.cache/go-build \
		-v $(GO_MOD_CACHE):/.cache/mod -e GOMODCACHE=/.cache/mod \
		-v $(GOLANGCI_LINT_CACHE):/.cache/golangci-lint -e GOLANGCI_LINT_CACHE=/.cache/golangci-lint \
		-v $(ROOT_DIR):/app \
		-w /app \
		$(LINT_IMAGE) golangci-lint run -v --new-from-rev=origin/master --build-tags=examples,allure_go,provider

.PHONY: cover
cover:
	go test -v -coverprofile=coverage.out ./...  && go tool cover -html=coverage.out

.PHONY: example
example:
	go test ./... -tags example

.PHONY: test
test:
	go test ./...
