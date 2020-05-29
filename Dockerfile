From harbor.oneitfarm.com/library/golang:1.12

COPY ./build/cproxy-linux /opt

WORKDIR /opt

RUN chmod +x cproxy-linux

ENTRYPOINT ./cproxy-linux
