#!/bin/bash
#echo "Building cichain"
#docker-compose up --no-start raftleveldb ci0 &>/dev/null
#
#echo "Start cichain"
#docker-compose start raftleveldb ci0 &>/dev/null

#echo "Building ethereum"
#docker-compose build ethereum
#
#echo "Starting ethereum"
#docker-compose up --no-start ethereum &>/dev/null
#docker-compose start ethereum &>/dev/null

#echo "Applying contracts"
#docker-compose build contract_deployer
contractAddress=$(docker-compose up contract_deployer | grep "Gravity deployed at Address" | grep -Eow '0x[0-9a-fA-F]{40}')
echo "Contract address: $contractAddress"