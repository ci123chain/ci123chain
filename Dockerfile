FROM harbor.oneitfarm.com/zhirenyun/baseimage:bionic-1.0.0

WORKDIR /opt/ci123chain

COPY ./docker/node/build/cid-linux /opt/cid-linux
COPY ./docker/node/build/cli-linux /opt/cli-linux

COPY /usr/local/go/pkg/mod/github.com/wasmerio /usr/local/go/pkg/mod/github.com/wasmerio

ENV GOPATH /go

COPY ./docker/node/2start.sh /etc/service/ci123chain/run

WORKDIR /opt
RUN chmod +x /etc/service/ci123chain/run