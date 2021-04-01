#! /bin/bash

docker-compose -f master.yaml up -d
sleep 10
docker-compose -f slave.yaml up -d