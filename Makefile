GO_BIN?=go

BUILD_DIR?=./build

CID?=$(BUILD_DIR)/cid

CLI?=$(BUILD_DIR)/cli

GO_BUILD_CMD=$(GO_BIN) build

.PHONY: build
build: server cli

server:
	$(GO_BUILD_CMD) -o $(CID) ./cmd/cid

cli:
	$(GO_BUILD_CMD) -o $(CLI) ./cmd/cicli