From harbor.oneitfarm.com/library/golang:1.12

COPY ./build/cproxy-linux /opt

WORKDIR /opt

ENV CI_PORT=80

RUN chmod +x cproxy-linux

ENTRYPOINT ./cproxy-linux
