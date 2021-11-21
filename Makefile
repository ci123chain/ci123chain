GO_BIN?=go

Tag?=v0.0.1

BUILD_DIR?=./build

CID?=$(BUILD_DIR)/cid

CLI?=$(BUILD_DIR)/cli

LB?=$(BUILD_DIR)/cproxy

GO_BUILD_CMD=$(GO_BIN) build


PROXY=https://goproxy.cn,direct

.PHONY: build
build: server
server:
	GOPROXY=$(PROXY) $(GO_BUILD_CMD) -o $(CID) ./cmd/cid
cli:
	GOPROXY=$(PROXY) $(GO_BUILD_CMD) -o $(CLI) ./cmd/cicli

local-start:
	docker-compose -f bootstrap-docker/single-node.yaml up -d
local-stop:
	docker-compose -f bootstrap-docker/single-node.yaml down

#.PHONY:release, build all
release:
	$(GO_BUILD_CMD) -o ./docker/node/build/cid-linux ./cmd/cid
	$(GO_BUILD_CMD) -o ./docker/node/build/tcptest ./cmd/test
	mv /go/pkg/mod/github.com/ci123chain/wasmer-go@v1.0.3-rc2 ./wasmer-go@v1.0.3-rc2

release-build:
	docker build -f Dockerfile_local -t cichain:$(Tag) .

