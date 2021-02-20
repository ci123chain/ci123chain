FROM harbor.oneitfarm.com/zhirenyun/baseimage:bionic-1.0.0

WORKDIR /opt/ci123chain

COPY . /opt/ci123chain/

COPY --from=0 ./docker/node/build/cid-linux /opt/cid-linux
COPY --from=0 ./docker/node/build/cli-linux /opt/cli-linux

COPY --from=0 /go/pkg/mod/github.com/wasmerio /go/pkg/mod/github.com/wasmerio

ENV GOPATH /go

COPY ./docker/node/2start.sh /etc/service/ci123chain/run

WORKDIR /opt
RUN chmod +x /etc/service/ci123chain/run