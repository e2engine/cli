ROOT_PATH := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
COVERAGE_PATH := $(ROOT_PATH).coverage/
BIN_PATH ?= $(ROOT_PATH)bin

# Use the locally installed Go toolchain.
export GOTOOLCHAIN=local

# Ensure our local tools are preferred
export PATH := $(ROOT_PATH)tools/bin:$(PATH)

include $(ROOT_PATH)tools/tools.mk

.PHONY: mockgen
mockgen: install-mockgen
	@echo "Generating mocks..."
	@go generate ./...

.PHONY: lint
lint: install-golangci-lint
	@echo "Running Go linter..."
	$(GOLANGCI_LINT) run

.PHONY: test
test:
	@echo "Running tests..."
	@go clean -testcache
	@go test ./... -count=1 -timeout=600s

.PHONY: test-cov
test-cov:
	@echo "Running tests with coverage..."
	@go clean -testcache
	@rm -rf $(COVERAGE_PATH)
	@mkdir -p $(COVERAGE_PATH)
	@go test -v -coverpkg=./... ./... -coverprofile $(COVERAGE_PATH)coverage.txt -count=1 -timeout=600s -tags=debug
	@go tool cover -func=$(COVERAGE_PATH)coverage.txt -o $(COVERAGE_PATH)functions.txt
	@go tool cover -html=$(COVERAGE_PATH)coverage.txt -o $(COVERAGE_PATH)coverage.html

.PHONY: test-race
test-race:
	@echo "Running tests with race detector..."
	@go clean -testcache
	@go test ./... -race -count=1 -timeout=600s

.PHONY: verify
verify: mockgen lint test-race
	@echo "All verifications passed successfully."

PKG := github.com/e2engine/cli/internal
VERSION_PKG := github.com/e2engine/cli/internal/command/version

VERSION ?= $(shell git describe --tags --exact-match 2>/dev/null)

BUILD_TIME ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT ?= $(shell git rev-parse HEAD)
GIT_TREE_STATE ?= $(shell if git diff --quiet --ignore-submodules HEAD; then echo "clean"; else echo "dirty"; fi)

GOOS ?= $(shell uname -s | tr '[:upper:]' '[:lower:]')
GOARCH ?= $(shell uname -m)
ifeq ($(GOARCH), x86_64)
	override GOARCH = amd64
endif

.PHONY: build-cli
build-cli:
	@echo "Building E2Engine CLI..."
	@mkdir -p $(BIN_PATH)
	@rm -f $(BIN_PATH)/e2engine-*
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build \
		-ldflags "-s -w \
		-X $(VERSION_PKG).version=$(VERSION) \
		-X $(VERSION_PKG).buildTime=$(BUILD_TIME) \
		-X $(VERSION_PKG).gitCommit=$(GIT_COMMIT) \
		-X $(VERSION_PKG).gitTreeState=$(GIT_TREE_STATE)" \
		-o $(BIN_PATH)/e2engine-$(GOOS)-$(GOARCH)$(shell [ "$(GOOS)" = windows ] && echo .exe) \
		./cmd/cli

.PHONY: build-image
build-image:
	@echo "Building E2Engine Docker image..."
	@rm -rf $(BIN_PATH)/linux
	@mkdir -p $(BIN_PATH)/linux/$(GOARCH)
	CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) go build \
		-ldflags "-s -w \
		-X $(VERSION_PKG).version=$(VERSION) \
		-X $(VERSION_PKG).buildTime=$(BUILD_TIME) \
		-X $(VERSION_PKG).gitCommit=$(GIT_COMMIT) \
		-X $(VERSION_PKG).gitTreeState=$(GIT_TREE_STATE)" \
		-o $(BIN_PATH)/linux/$(GOARCH)/e2engine \
		./cmd/cli
	@docker buildx build \
		--platform linux/$(GOARCH) \
		--load \
		-t e2engine:local \
		-f Dockerfile \
		$(BIN_PATH)
