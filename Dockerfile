FROM harbor.oneitfarm.com/zhirenyun/baseimage:bionic-1.0.0

COPY ./build/cproxy-linux /opt

WORKDIR /opt

RUN chmod +x cproxy-linux

ENTRYPOINT ./cproxy-linux --port=80
