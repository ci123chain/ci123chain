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


.PHONY: build-linux
build-linux: server-linux client-linux
server-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_BUILD_CMD) -o $(CID)-linux ./cmd/cid
client-linux: 
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_BUILD_CMD) -o $(CLI)-linux ./cmd/cicli



.PHONY: build-docker
build-docker: build-linux build-image
build-image:
	docker build -t cichain:v0.0.1 .
	docker run --name ci123-container-v1 -p 1318:1317 -d cichain:v0.0.1

docker-clean:
	docker ps -a | grep "ci123-container-" | awk '{print $$1}' | xargs docker rm -f
	docker images | grep "cichain" | awk '{print $$3}' | xargs docker rmi

docker-stop:
	docker ps -a | grep "ci123-container-" | awk '{print $$1}' | xargs docker stop

docker-restart:
	docker ps -a | grep "ci123-container-" | awk '{print $$1}' | xargs docker start