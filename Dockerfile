From golang

COPY ./build/cli-linux /opt
COPY ./build/cid-linux /opt
COPY ./docker/start.sh /opt

WORKDIR /opt

RUN chmod +x cli-linux
RUN chmod +x cid-linux
RUN chmod +x start.sh

ENTRYPOINT ./start.sh


