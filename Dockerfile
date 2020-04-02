#From harbor.oneitfarm.com/library/golang:1.12
From ubuntu
COPY ./docker/node/build/cli-linux /opt
COPY ./docker/node/build/cid-linux /opt
COPY ./docker/node/2start.sh /opt


WORKDIR /opt

RUN chmod +x cli-linux
RUN chmod +x cid-linux
RUN chmod +x 2start.sh

ENTRYPOINT ./2start.sh


