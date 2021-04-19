#!/bin/bash
docker-compose -f bootstrap-docker/single-node.yaml down
docker-compose -f bootstrap-docker/single-node-db.yaml down
./node_reset.sh

docker-compose -f bootstrap-docker/single-node.yaml up -d
docker-compose -f bootstrap-docker/single-node-db.yaml up -d

