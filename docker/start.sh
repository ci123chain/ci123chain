#!/bin/bash

HOME_DIR="$HOME"
if [ $CI123_HOME ];then
    HOME_DIR=$CI123_HOME
fi

CLI_HOME="${HOME_DIR}/cli"
CID_HOME="${HOME_DIR}/cid"

# genesis file
if [ ! -f ${CID_HOME}/config/genesis.json ]; then
    ./cid-linux init --home=$CID_HOME
    ./cid-linux add-genesis-account 0x204bCC42559Faf6DFE1485208F7951aaD800B313 100000 --home=$CID_HOME
    ./cid-linux add-genesis-account 0x505A74675dc9C71eF3CB5DF309256952917E801e 100000 --home=$CID_HOME
    ./cid-linux add-genesis-account 0xD1a14962627fAc768Fe885Eeb9FF072706B54c19 100000 --home=$CID_HOME
fi

# start
nohup ./cid-linux start --home=$CID_HOME > cid-output 2>&1 &

./cli-linux rest-server --laddr=tcp://0.0.0.0:80 --home=$CLI_HOME > rest-output 2>&1