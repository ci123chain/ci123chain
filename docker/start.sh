#!/bin/bash

if [ ! -d "$HOME/.ci123" ]; then
    ./cid-linux init --address=0x204bCC42559Faf6DFE1485208F7951aaD800B313
fi

nohup ./cid-linux start > cid-output 2>&1 &

./cli-linux rest-server > rest-output 2>&1



