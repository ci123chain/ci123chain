#! /bin/bash

mkdir node0 node1 node2 gateway
docker-compose -f part1.yaml up -d
sleep 15
docker network create testdocker_couchDB
docker-compose -f part22.yaml up -d