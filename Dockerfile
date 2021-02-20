FROM harbor.oneitfarm.com/zhirenyun/go:1.15.6

WORKDIR /opt/ci123chain

COPY . /opt/ci123chain/

RUN GOPROXY=https://goproxy.cn,direct GOSUMDB=off go build -o /opt/cid-linux ./cmd/cid
RUN GOPROXY=https://goproxy.cn,direct GOSUMDB=off go build -o /opt/cli-linux ./cmd/cicli
RUN GOPROXY=https://goproxy.cn,direct GOSUMDB=off go build -o /opt/cproxy-linux ./cmd/gateway

FROM harbor.oneitfarm.com/zhirenyun/baseimage:bionic-1.0.0

COPY --from=0 /opt/cid-linux /opt/cid-linux
COPY --from=0 /opt/cli-linux /opt/cli-linux
COPY --from=0 /opt/cproxy-linux /opt/cproxy-linux

COPY --from=0 /go/pkg/mod/github.com/wasmerio /go/pkg/mod/github.com/wasmerio

ENV GOPATH /go

COPY ./docker/node/2start.sh /etc/service/ci123chain/run

WORKDIR /opt
RUN chmod +x /etc/service/ci123chain/run