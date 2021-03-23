GO_BIN?=go

Tag?=v0.0.1

BUILD_DIR?=./build

LB?=$(BUILD_DIR)/lb

GO_BUILD_CMD=$(GO_BIN) build

PROXY=https://goproxy.io,direct

.PHONY: build-cproxy
build-cproxy:
	GOPROXY=$(PROXY) $(GO_BUILD_CMD) -o $(LB) ./cmd

build-cproxy-linux:
	GOPROXY=$(PROXY) CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_BUILD_CMD) -o ./build/cproxy-linux ./cmd

build-cproxy-image: build-cproxy-linux
	docker build -t cproxyservice:$(Tag) .
start-cproxy: build-cproxy-image
	docker run --name ci123-cproxy-v1 -p 3030:3030 -d cproxyservice:$(Tag)
clean-cproxy:
	docker ps -a | grep "ci123-$(Tag)-" | awk '{print $$1}' | xargs docker rm -f
	docker images | grep "$(Tag)service" | awk '{print $$3}' | xargs docker rmi

# .PHONY:release, build all
release:
    $(GO_BUILD_CMD) -o ./build/cproxy-linux ./cmd