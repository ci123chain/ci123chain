From harbor.oneitfarm.com/library/golang:1.12

WORKDIR /opt/ci123chain

COPY . /opt/ci123chain/

RUN GOPROXY=https://goproxy.io go build -o /opt/cid-linux ./cmd/cid
RUN GOPROXY=https://goproxy.io go build -o /opt/cli-linux ./cmd/cicli
RUN GOPROXY=https://goproxy.io go build -o /opt/cproxy-linux ./cmd/gateway

COPY ./docker/node/2start.sh /opt

WORKDIR /opt
RUN chmod +x 2start.sh
ENTRYPOINT ./2start.sh


