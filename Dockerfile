FROM harbor.oneitfarm.com/zhirenyun/baseimage:bionic-1.0.0


ENV CI_PORT=80

COPY ./build/cproxy-linux /etc/service/cproxy-linux/run

RUN chmod +x /etc/service/cproxy-linux/run

