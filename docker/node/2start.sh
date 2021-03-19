#!/bin/bash
if [ $GATEWAY ]; then
    /opt/cproxy-linux
    exit 0
fi

if [ -z $CI_HOME ];
then
   CI_HOME="/root/cid"
fi

CI_LOGDIR=$CI_HOME/logs
if [ ! -d $CI_LOGDIR ]; then
    mkdir -p $CI_LOGDIR
fi


if [ $LITECLIENT ]; then
    if [ -z $CONNECT_NODE_ADDRESS ] || [ -z $CONNECT_CHAIN_ID ]; then
        echo "---Please Special FULL_NODE_ADDRESS and CONNECT_CHAIN_ID----"
        exit 1
    fi
    nohup /opt/cid-linux tendermint lite --node=$CONNECT_NODE_ADDRESS --chain-id=$CONNECT_CHAIN_ID --home-dir=$CI_HOME >> $CI_LOGDIR/liteclient-output.log 2>&1 &
    /opt/cli-linux rest-server --laddr=tcp://0.0.0.0:80 --node=tcp://0.0.0.0:8888 >> $CI_LOGDIR/cid-output.log 2>&1
    exit 0
fi

# genesis file
if [ ! -f $CI_HOME/config/genesis.json ]; then
    if [ -z $CI_MASTER_DOMAIN]; then
      if [ $CI_CONFIG ]; then
        echo "---Found CI_CONFIG env----"
      else
        echo "---Not found genesis file, Creating----"

        #./cid-linux init
        /opt/cid-linux init --home=$CI_HOME --chain_id=$CI_CHAIN_ID --validator_key=$CI_VALIDATOR_KEY

        /opt/cid-linux add-genesis-account 0x204bCC42559Faf6DFE1485208F7951aaD800B313 10000000000000000000000000000 --home=$CI_HOME
        # a78a8a281d160847f1ed7881e5497e1a98ccd4fe6ba9ce918630f93a44e09793
        /opt/cid-linux add-genesis-account 0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c 10000000000000000000000000000 --home=$CI_HOME
        # 2b452434ac4f7cf9c5d61d62f23834f34e851fb6efdb8d4a8c6e214a8bc93d70

        /opt/cid-linux add-genesis-validator 0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c 800000000 $CI_PUBKEY 1 40 5 --home=$CI_HOME

        if [ $GENESIS_SHARED ]; then
            /opt/cid-linux add-genesis-shard "$GENESIS_SHARED"
            #./cid-linux add-genesis-shard "ci0:0;ci1:0;ci2:0"
        fi
      fi
    fi
else
    echo "---Found genesis file----"
    cat $CI_HOME/config/genesis.json
    echo "----------"
fi


# start
nohup /opt/cid-linux start --home=$CI_HOME >> $CI_LOGDIR/cid-output.log 2>&1 &
sleep 10
/opt/cli-linux rest-server --laddr=tcp://0.0.0.0:80 >> $CI_LOGDIR/rest-output.log 2>&1