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
   CI_TOKENNAME="WLK"
fi

CI_LOGDIR=$CI_HOME/logs
if [ ! -d $CI_LOGDIR ]; then
    mkdir -p $CI_LOGDIR
fi
# delay 3 second for cid
sleep 1s

source /etc/profile
if [ -z $CI_VALIDATOR_KEY ];
then
  echo "CI_VALIDATOR_KEY not exist"
  exit 0
fi

echo "---Loading CLI ENV---"
echo "export CI_ETH_CHAIN_ID=$CI_ETH_CHAIN_ID" >> /etc/profile
echo "export CI_HOME=$CI_HOME" >> /etc/profile
source /etc/profile

echo "---Start cli---"
/opt/cli-linux rest-server --laddr=tcp://0.0.0.0:80