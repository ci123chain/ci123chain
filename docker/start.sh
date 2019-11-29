#!/bin/bash

# genesis file
if [ ! -d "$HOME/.ci123" ]; then
    ./cid-linux init
fi
./cid-linux add-genesis-account 0x204bCC42559Faf6DFE1485208F7951aaD800B313 100000
./cid-linux add-genesis-account 0x505A74675dc9C71eF3CB5DF309256952917E801e 100000
./cid-linux add-genesis-account 0xD1a14962627fAc768Fe885Eeb9FF072706B54c19 100000

# start
nohup ./cid-linux start > cid-output 2>&1 &

./cli-linux rest-server > rest-output 2>&1




