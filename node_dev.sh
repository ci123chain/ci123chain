#!/bin/bash


#export CI_HOME="/Users/tanhui/Documents/work/golang/ci123chain/docker/node/test3"
export CI_VALIDATOR_KEY="BpiPMqeXTNJLBg6hRoWJjsjOLSHIQRIrkQlykLQ/0AE="
export CI_PUBKEY="Ap0lbWGAnzfqpc0D0GL081WCnWatdk2d5B21orPl30AS"
export CI_CHAIN_ID="testchain123"
export CI_STATEDB="couchdb://admin:123rewQAQtre56@193.112.144.129:5984/test1230"


if [ -z $CI_HOME ];
then
   CI_HOME="/root/cid"
fi

# genesis file
if [ ! -f $CI_HOME/config/genesis.json ]; then
    echo "---Not found genesis file, Creating----"

    #./cid-linux init
    ./cid-linux init --home=$CI_HOME --chain_id=$CI_CHAIN_ID --validator_key=$CI_VALIDATOR_KEY

    ./cid-linux add-genesis-account 0x204bCC42559Faf6DFE1485208F7951aaD800B313 10000000000 --home=$CI_HOME
    # a78a8a281d160847f1ed7881e5497e1a98ccd4fe6ba9ce918630f93a44e09793
    ./cid-linux add-genesis-account 0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c 10000000000 --home=$CI_HOME
    # 2b452434ac4f7cf9c5d61d62f23834f34e851fb6efdb8d4a8c6e214a8bc93d70

    ./cid-linux add-genesis-account 0xB6727FCbC60A03A6689AEE6E5fBC83a7FDc9beBf 10000000000 --home=$CI_HOME

    ./cid-linux add-genesis-validator 0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c 8000000 $CI_PUBKEY 1 40 5 --home=$CI_HOME

    if [ $GENESIS_SHARED ]; then
        ./cid-linux add-genesis-shard "$GENESIS_SHARED"
        #./cid-linux add-genesis-shard "ci0:0;ci1:0;ci2:0"
    fi
else
    echo "---Found genesis file----"
    cat $CI_HOME/config/genesis.json
    echo "----------"
fi

CI_LOGDIR=$CI_HOME/logs
if [ ! -d $CI_LOGDIR ]; then
    mkdir $CI_LOGDIR
fi

# start
nohup ./cli-linux rest-server --laddr=tcp://0.0.0.0:80 >> $CI_LOGDIR/rest-output.log 2>&1 &
./cid-linux start >> $CI_LOGDIR/cid-output.log 2>&1
