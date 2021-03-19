GO_BIN?=go

Tag?=v0.0.1

BUILD_DIR?=./build

CID?=$(BUILD_DIR)/cid

CLI?=$(BUILD_DIR)/cli

LB?=$(BUILD_DIR)/cproxy

GO_BUILD_CMD=$(GO_BIN) build


PROXY=https://goproxy.cn,direct

.PHONY: build
build: server cli
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
	GOPROXY=$(PROXY) GOSUMDB=off $(GO_BUILD_CMD) -o ./docker/node/build/cid-linux ./cmd/cid
	GOPROXY=$(PROXY) GOSUMDB=off $(GO_BUILD_CMD) -o ./docker/node/build/cli-linux ./cmd/cicli
	mv /go/pkg/mod/github.com/wasmerio ./wasmerio

release-build:
	docker build -f Dockerfile_local -t cichain:$(Tag) .

