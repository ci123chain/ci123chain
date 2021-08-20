#! /bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

rm -rf $CI_HOME

$DIR/../build/cid init --home=$CI_HOME --chain_id=$CI_CHAIN_ID --validator_key=$CI_VALIDATOR_KEY --denom=$CI_TOKENNAME

$DIR/../build/cid add-genesis-account 0x2d82ada163ee7c9287f07f6871c190cca8562e85 10000000000000000000000000000 --home $CI_HOME
#2b452434ac4f7cf9c5d61d62f23834f34e851fb6efdb8d4a8c6e214a8bc93d70

$DIR/../build/cid add-genesis-validator 0x2d82ada163ee7c9287f07f6871c190cca8562e85 8000000 $CI_PUBKEY 1 40 5 --home $CI_HOME
