#! /bin/bash

export CI_VALIDATOR_KEY="ZZyS+fsb1zlVDfTTcTXfCzEJ9vKzBERFaCj/jQ3xdOmYo73TTRHfeFLFJQc2uJWRH+x+7pzD6OaXzxRLvC9vKA=="
export CI_PUBKEY="mKO9000R33hSxSUHNriVkR/sfu6cw+jml88US7wvbyg="
export CI_CHAIN_ID="ibc0"
export CI_NODE_DOMAIN="localhost"
export CI_STATEDB_HOST="127.0.0.1"
export CI_STATEDB_PORT=5002
export CI_STATEDB_TLS="false"
export CI_TOKENNAME="stake"
export CI_HOME="scripts/testdata/ibc0"

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
bash $DIR/reset_base.sh