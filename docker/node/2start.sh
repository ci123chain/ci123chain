#!/bin/bash

if [ -z $CI_HOME ];
then
   CI_HOME="/opt/ci123chain"
fi

if [ -z $CI_ETH_CHAIN_ID ];
then
   CI_ETH_CHAIN_ID=7880
fi

if [ -z $CI_TOKENNAME ];
then
   CI_TOKENNAME="stake"
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
if [ ! -f $CI_HOME/config/genesis.json ]; then # first create

    if [ -z $CI_MASTER_DOMAIN ]; then # self create, first node
        echo "---Not found genesis file, Creating----"

        #./cid-linux init
        /opt/cid-linux init --home=$CI_HOME --chain_id=$CI_CHAIN_ID --denom=$CI_TOKENNAME

        CI_VALIDATOR_KEY=$(cat $CI_HOME/config/priv_validator_key.json | jq -r '.priv_key.value')
        CI_PUBKEY=$(cat $CI_HOME/config/priv_validator_key.json | jq -r '.pub_key.value')

        /opt/cid-linux add-genesis-account 0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c 10000000000000000000000000000 --home=$CI_HOME
        # 2b452434ac4f7cf9c5d61d62f23834f34e851fb6efdb8d4a8c6e214a8bc93d70


        if [ $CI_VALIDATOR_ADDRESS ]; then
          if [ -z $CI_GENESIS_AMOUNT ];then
            CI_GENESIS_AMOUNT=10000000000000000000000000000
          fi
          /opt/cid-linux add-genesis-account $CI_VALIDATOR_ADDRESS $CI_GENESIS_AMOUNT --home=$CI_HOME
          /opt/cid-linux add-genesis-validator $CI_VALIDATOR_ADDRESS 800000000 $CI_PUBKEY 1 40 5 --home=$CI_HOME --moniker=$CI_NODE_NAME
        else
          /opt/cid-linux add-genesis-validator 0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c 800000000 $CI_PUBKEY 1 40 5 --home=$CI_HOME --moniker=$CI_NODE_NAME
        fi

        if [ $GENESIS_SHARED ]; then
            /opt/cid-linux add-genesis-shard "$GENESIS_SHARED"
            #./cid-linux add-genesis-shard "ci0:0;ci1:0;ci2:0"
        fi
    else # second node
        echo "---Create Validator----"
        /opt/cid-linux gen-validator --home=$CI_HOME
        CI_VALIDATOR_KEY=$(cat $CI_HOME/config/priv_validator_key.json | jq -r '.priv_key.value')
        CI_PUBKEY=$(cat $CI_HOME/config/priv_validator_key.json | jq -r '.pub_key.value')
        echo "export CI_VALIDATOR_KEY=$CI_VALIDATOR_KEY" >> /etc/profile
        echo "export CI_PUBKEY=$CI_PUBKEY" >> /etc/profile
        source /etc/profile
    fi
else
    echo "---Found genesis file----"
    echo "----------"
fi

if [ -f $CI_HOME/config/config.toml ]; then
    sed -i "s/max_subscriptions_per_client = 5/max_subscriptions_per_client = 20/" $CI_HOME/config/config.toml
fi

CI_VALIDATOR_KEY=$(cat $CI_HOME/config/priv_validator_key.json | jq -r '.priv_key.value')
CI_PUBKEY=$(cat $CI_HOME/config/priv_validator_key.json | jq -r '.pub_key.value')
echo "export CI_VALIDATOR_KEY=$CI_VALIDATOR_KEY" >> /etc/profile
echo "export CI_PUBKEY=$CI_PUBKEY" >> /etc/profile
echo "export CI_ETH_CHAIN_ID=$CI_ETH_CHAIN_ID" >> /etc/profile
echo "export CI_HOME=$CI_HOME" >> /etc/profile
source /etc/profile

echo "---Start cid---"
/opt/cid-linux start --laddr=tcp://0.0.0.0:80 --home=$CI_HOME


