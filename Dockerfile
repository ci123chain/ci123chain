FROM harbor.oneitfarm.com/zhirenyun/baseimage:bionic-1.0.0

WORKDIR /opt/ci123chain

COPY ./docker/node/build/cid-linux /opt/cid-linux
COPY ./docker/node/build/cli-linux /opt/cli-linux
COPY ./docker/node/build/tcptest /opt/tcptest
COPY ./docker/node/exportFile.json /opt/exportFile.json

COPY ./wasmer-go@v1.0.3-rc2 /go/pkg/mod/github.com/ci123chain/wasmer-go@v1.0.3-rc2
# For single package
COPY ./wasmer-go@v1.0.3-rc2 /opt/wasmer-go@v1.0.3-rc2

COPY ./migrate.sh /opt

ENV GOPATH /go

COPY ./docker/node/start-cid.sh /etc/service/cid/run
#COPY ./docker/node/start-cli.sh /etc/service/cli/run

RUN apt-get update -y && apt-get install jq -y

WORKDIR /opt
RUN chmod +x /etc/service/cid/run
#RUN chmod +x /etc/service/cli/run