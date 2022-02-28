FROM harbor.oneitfarm.com/zhirenyun/baseimage:bionic-1.0.0

WORKDIR /opt/ci123chain

COPY ./docker/node/build/cid-linux /opt/cid-linux
COPY ./docker/node/build/tcptest /opt/tcptest
COPY ./docker/node/build/upgrade /opt/upgrade
COPY ./docker/node/exportFile.json /opt/exportFile.json

COPY ./wasmer-go@v1.0.3-rc2 /go/pkg/mod/github.com/ci123chain/wasmer-go@v1.0.3-rc2
# For single package
COPY ./wasmer-go@v1.0.3-rc2 /opt/wasmer-go@v1.0.3-rc2

COPY ./migrate.sh /opt

ENV GOPATH /go
ENV DAEMON_NAME cid-linux

COPY ./docker/node/start-cid.sh /etc/service/cid/run
COPY ./jq-linux64 /usr/bin/jq
COPY ./serve_linux_amd64 /usr/bin/file_server
RUN chmod +x /usr/bin/jq
RUN chmod +x /usr/bin/file_server

WORKDIR /opt
RUN chmod +x /etc/service/cid/run