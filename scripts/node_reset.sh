#! /bin/bash

export CI_VALIDATOR_KEY="EHNcQ+4ej1ODR6JT6AD7ooQ44yKZFCcvPTT95r92dSim6OSlW11vkulfOOdHg71VkncPsB5fZm/6ieR9yJltBQ=="
export CI_PUBKEY="nU3F8neWPfNkNwlHYfodV3Ar7aSFG/uauq7SfFuEhBc="
export CI_CHAIN_ID="ibc2"
export CI_NODE_DOMAIN="localhost"
export CI_STATEDB_HOST="127.0.0.1"
export CI_STATEDB_PORT=5001
export CI_STATEDB_TLS="false"
export CI_TOKENNAME="stake"
export CI_HOME="scripts/testdata/ibc2"

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
bash $DIR/reset_base.sh