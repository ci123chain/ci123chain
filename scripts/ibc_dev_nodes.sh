#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

docker-compose -f $DIR/../bootstrap-docker/two-node.yaml down
sleep 3s
docker-compose -f $DIR/../bootstrap-docker/two-node.yaml up -d