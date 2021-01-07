FROM harbor.oneitfarm.com/zhirenyun/baseimage:bionic-1.0.0

COPY ./build/cproxy-linux /opt

WORKDIR /opt

ENV CI_PORT=80

RUN chmod +x cproxy-linux

ENTRYPOINT ./cproxy-linux
