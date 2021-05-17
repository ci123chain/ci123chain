#! /bin/bash

export CI_VALIDATOR_KEY="4wttMiieaewLiRYu+y05j0uslBDOX5IA3k4TY9GtQzSdTcXyd5Y982Q3CUdh+h1XcCvtpIUb+5q6rtJ8W4SEFw=="
export CI_PUBKEY="nU3F8neWPfNkNwlHYfodV3Ar7aSFG/uauq7SfFuEhBc="
export CI_CHAIN_ID="ibc-0"
export CI_NODE_DOMAIN="localhost"
export CI_STATEDB_HOST="127.0.0.1"
export CI_STATEDB_PORT=5002
export CI_STATEDB_TLS="false"
export CI_TOKENNAME="stack0"
export CI_HOME="scripts/testdata/ibc0"

rm -rf $CI_HOME

./build/cid init --home=$CI_HOME --chain_id=$CI_CHAIN_ID --validator_key=$CI_VALIDATOR_KEY --denom=$CI_TOKENNAME

./build/cid add-genesis-account 0xD1a14962627fAc768Fe885Eeb9FF072706B54c19 10000000000000000000000000000 --home $CI_HOME

./build/cid add-genesis-account 0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c 10000000000000000000000000000 --home $CI_HOME
#2b452434ac4f7cf9c5d61d62f23834f34e851fb6efdb8d4a8c6e214a8bc93d70

./build/cid add-genesis-account 0x204bCC42559Faf6DFE1485208F7951aaD800B313 10000000000000000000000000000 --home $CI_HOME

./build/cid add-genesis-validator 0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c 8000000 $CI_PUBKEY 1 40 5 --home $CI_HOME
