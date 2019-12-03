GO_BIN?=go

BUILD_DIR?=./build

CID?=$(BUILD_DIR)/cid

CLI?=$(BUILD_DIR)/cli

GO_BUILD_CMD=$(GO_BIN) build

MOD?=vendor

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
	GOPROXY=$(PROXY) CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_BUILD_CMD) -o $(CID)-linux ./cmd/cid
client-linux: 
	GOPROXY=$(PROXY) CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_BUILD_CMD) -o $(CLI)-linux ./cmd/cicli

build-image:
	docker build -t cichain:v0.0.1 .

.PHONY: build-docker
build-docker: build-linux build-image

.PHONY: docker-start
docker-start: build-docker start
start:
	docker run --name ci123-container-v1 -p 1318:1317 -e CI123_HOME=/opt -d cichain:v0.0.1

docker-clean:
	docker ps -a | grep "ci123-container-" | awk '{print $$1}' | xargs docker rm -f
	docker images | grep "cichain" | awk '{print $$3}' | xargs docker rmi

docker-stop:
	docker ps -a | grep "ci123-container-" | awk '{print $$1}' | xargs docker stop

docker-restart:
	docker ps -a | grep "ci123-container-" | awk '{print $$1}' | xargs docker start

.PHONY:release
release: build-linux

