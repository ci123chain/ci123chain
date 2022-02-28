GO_BIN?=go

Tag?=v0.0.1

BUILD_DIR?=./build

CID?=$(BUILD_DIR)/cid

CLI?=$(BUILD_DIR)/cli

LB?=$(BUILD_DIR)/cproxy

GO_BUILD_CMD=$(GO_BIN) build

VERSION := $(shell echo $(shell git describe --always --match "v*") | sed 's/^v//')
ldflags = -X github.com/ci123chain/ci123chain/pkg/abci/version.Name=ci123chain \
        -X github.com/ci123chain/ci123chain/pkg/abci/version.Version=$(VERSION)
BUILD_FLAGS := -ldflags '$(ldflags)'

PROXY=https://goproxy.cn,direct

.PHONY: build
build: server
server:
	GOPROXY=$(PROXY) $(GO_BUILD_CMD) $(BUILD_FLAGS) -o $(CID) ./cmd/cid

local-start:
	docker-compose -f bootstrap-docker/single-node.yaml up -d
local-stop:
	docker-compose -f bootstrap-docker/single-node.yaml down

#.PHONY:release, build all
release:
	$(GO_BUILD_CMD) $(BUILD_FLAGS) -o ./docker/node/build/cid-linux ./cmd/cid
	$(GO_BUILD_CMD) $(BUILD_FLAGS) -o ./docker/node/build/tcptest ./cmd/test
	mv /go/pkg/mod/github.com/ci123chain/wasmer-go@v1.0.3-rc2 ./wasmer-go@v1.0.3-rc2
	cd ./cmd/upgrade && $(GO_BUILD_CMD) $(BUILD_FLAGS) -o ../../docker/node/build/upgrade .

release-build:
	docker build -f Dockerfile_local -t cichain:$(Tag) .

upgrade-build:
	docker build -f Dockerfile_upgrade -t cichainupgrade:v0.0.1 .
	docker run -v $(PWD)/bin:/opt/temp cichainupgrade:v0.0.1 bash -c "cp /opt/cid-linux /opt/temp"