#!/bin/bash

HOME_DIR="$HOME"
if [ $CI123_HOME ];then
    HOME_DIR=$CI123_HOME
fi

CLI_HOME="${HOME_DIR}/cli"
CID_HOME="${HOME_DIR}/cid"

# genesis file
if [ ! -f ${CID_HOME}/config/genesis.json ]; then
    ./cid-linux init --home=$CID_HOME --chain-id=${ShardID}
    ./cid-linux add-genesis-account 0x204bCC42559Faf6DFE1485208F7951aaD800B313 10000000000 --home=$CID_HOME
    # a78a8a281d160847f1ed7881e5497e1a98ccd4fe6ba9ce918630f93a44e09793
    ./cid-linux add-genesis-account 0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c 10000000000 --home=$CID_HOME
    # 2b452434ac4f7cf9c5d61d62f23834f34e851fb6efdb8d4a8c6e214a8bc93d70
fi

# start
nohup ./cid-linux start --home=$CID_HOME --statedb=couchdb://couchdb-service:5984 > cid-output 2>&1 &

./cli-linux rest-server --laddr=tcp://0.0.0.0:80 --home=$CLI_HOME > rest-output 2>&1