FROM harbor.oneitfarm.com/zhirenyun/baseimage:bionic-1.0.0

WORKDIR /opt/ci123chain

COPY ./docker/node/build/cid-linux /opt/cid-linux
COPY ./docker/node/build/cli-linux /opt/cli-linux
COPY ./docker/node/exportFile.json /opt/exportFile.json

COPY ./wasmer-go@v1.0.3-rc2 /go/pkg/mod/github.com/ci123chain/wasmer-go@v1.0.3-rc2
# For single package
COPY ./wasmer-go@v1.0.3-rc2 /opt/wasmer-go@v1.0.3-rc2

COPY ./migrate.sh /opt


ENV GOPATH /go

COPY ./docker/node/2start.sh /etc/service/ci123chain/run

WORKDIR /opt
RUN chmod +x /etc/service/ci123chain/run