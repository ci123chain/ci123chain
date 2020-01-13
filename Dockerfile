From ubuntu

COPY ./docker/node/build/cli-linux /opt
COPY ./docker/node/build/cid-linux /opt
COPY ./docker/gateway/build/gateway-linux /opt
COPY ./docker/node/start.sh /opt


WORKDIR /opt

RUN chmod +x cli-linux
RUN chmod +x cid-linux
RUN chmod +x gateway-linux
RUN chmod +x start.sh

ENTRYPOINT ./start.sh


