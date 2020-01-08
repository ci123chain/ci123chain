GO_BIN?=go

Tag?=v0.0.1

BUILD_DIR?=./build

CID?=$(BUILD_DIR)/cid

CLI?=$(BUILD_DIR)/cli

LB?=$(BUILD_DIR)/lb

GO_BUILD_CMD=$(GO_BIN) build


PROXY=https://goproxy.io

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
	docker build -t cichain:$(Tag) ./docker/node

.PHONY: build-docker
build-docker: build-linux build-image

.PHONY: node-start
node-start: build-docker simple-start
simple-start:
	docker run --name ci123-chain-v1 -p 1318:80 -p 26676:26656 -d cichain:$(Tag)

clean-node:
	docker ps -a | grep "ci123-chain-" | awk '{print $$1}' | xargs docker rm -f
	docker images | grep "cichain" | awk '{print $$3}' | xargs docker rmi

node-stop:
	docker ps -a | grep "ci123-chain-" | awk '{print $$1}' | xargs docker stop

node-restart:
	docker ps -a | grep "ci123-chain-" | awk '{print $$1}' | xargs docker start

.PHONY:release
release: build-linux


.PHONY: build-lb
build-lb:
	GOPROXY=$(PROXY) $(GO_BUILD_CMD) -o $(LB) ./cmd/lb

build-lb-linux:
	GOPROXY=$(PROXY) CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_BUILD_CMD) -o ./docker/lb/build/lb-linux ./cmd/lb

build-lb-image: build-lb-linux
	docker build -t lbservice:$(Tag) ./docker/lb
start-lb: build-lb-image
	docker run --name ci123-lb-v1 -p 3030:3030 -d lbservice:$(Tag)
clean-lb:
	docker ps -a | grep "ci123-lb-" | awk '{print $$1}' | xargs docker rm -f
	docker images | grep "lbservice" | awk '{print $$3}' | xargs docker rmi
