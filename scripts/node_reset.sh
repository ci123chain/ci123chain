#! /bin/bash

export CI_VALIDATOR_KEY="qS4LtbZ9nxk/5HOpGBLQaMLWjzgrfd81VbrxoiQqgZkIvOle+S28kv+u8136PAfvqRDTRnuGVlIEmbFprIRRFg=="
export CI_PUBKEY="CLzpXvktvJL/rvNd+jwH76kQ00Z7hlZSBJmxaayEURY="
export CI_CHAIN_ID="weelink"
export CI_ETH_CHAIN_ID=444900
export CI_NODE_DOMAIN="localhost"
export CI_STATEDB_HOST="127.0.0.1"
export CI_STATEDB_PORT=5002
export CI_STATEDB_TLS="false"
export CI_TOKENNAME="stake"
export CI_HOME="scripts/testdata/ibc4"

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
bash $DIR/reset_base.sh