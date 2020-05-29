FROM pretty66/cn-alpine:3.9

COPY ./build/cproxy-linux /opt

WORKDIR /opt

RUN chmod +x cproxy-linux

ENTRYPOINT ./cproxy-linux
