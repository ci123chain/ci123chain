#!/bin/bash

HOME_DIR="$HOME/.ci123"
if [ $CI123_HOME ];then
    HOME_DIR=$CI123_HOME
fi
# genesis file
echo $HOME_DIR/genesis.json
if [ ! -f $HOME_DIR/genesis.json ]; then
    ./cid-linux init --home=$HOME_DIR
fi
./cid-linux add-genesis-account 0x204bCC42559Faf6DFE1485208F7951aaD800B313 100000 --home=$HOME_DIR
./cid-linux add-genesis-account 0x505A74675dc9C71eF3CB5DF309256952917E801e 100000 --home=$HOME_DIR
./cid-linux add-genesis-account 0xD1a14962627fAc768Fe885Eeb9FF072706B54c19 100000 --home=$HOME_DIR

# start
nohup ./cid-linux start --home=$HOME_DIR > cid-output 2>&1 &

./cli-linux rest-server --home=$HOME_DIR > rest-output 2>&1