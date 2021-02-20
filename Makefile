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


.PHONY: build-linux
build-linux: server-linux client-linux
server-linux:
	GOPROXY=$(PROXY) CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_BUILD_CMD) -o ./docker/node/build/cid-linux ./cmd/cid
client-linux: 
	GOPROXY=$(PROXY) CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_BUILD_CMD) -o ./docker/node/build/cli-linux ./cmd/cicli

build-image:
	docker rmi cichain:$(Tag)
	docker build -t cichain:$(Tag) ./docker/node


.PHONY: docker-clean
docker-clean: clean-node
clean-node:
	docker ps -a | grep "ci123-chain-" | awk '{print $$1}' | xargs docker rm -f
	docker images | grep "cichain" | awk '{print $$3}' | xargs docker rmi

node-stop:
	docker ps -a | grep "ci123-chain-" | awk '{print $$1}' | xargs docker stop

node-restart:
	docker ps -a | grep "ci123-chain-" | awk '{print $$1}' | xargs docker start

.PHONY: build-cproxy
build-cproxy:
	GOPROXY=$(PROXY) $(GO_BUILD_CMD) -o $(LB) ./cmd/gateway

build-cproxy-linux:
	GOPROXY=$(PROXY) CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_BUILD_CMD) -o ./docker/gateway/build/cproxy-linux ./cmd/gateway

build-cproxy-image: build-cproxy-linux
	docker rmi cproxyservice:$(Tag)
	docker build -t cproxyservice:$(Tag) ./docker/gateway
start-cproxy: build-cproxy-image
	docker run --name ci123-cproxy-v1 -p 3030:3030 -d cproxyservice:$(Tag)
clean-cproxy:
	docker ps -a | grep "ci123-$(Tag)-" | awk '{print $$1}' | xargs docker rm -f
	docker images | grep "$(Tag)service" | awk '{print $$3}' | xargs docker rmi

#.PHONY:release, build all
release: build-linux

release-build:
	docker build -t cichain:$(Tag) .

