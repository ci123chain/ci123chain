#!/bin/bash
docker-compose -f ../bootstrap-docker/three-node-db.yaml down

export CI_VALIDATOR_KEY="4wttMiieaewLiRYu+y05j0uslBDOX5IA3k4TY9GtQzSdTcXyd5Y982Q3CUdh+h1XcCvtpIUb+5q6rtJ8W4SEFw=="
export CI_PUBKEY="nU3F8neWPfNkNwlHYfodV3Ar7aSFG/uauq7SfFuEhBc="
export CI_NODE_DOMAIN="localhost"
export CI_STATEDB_HOST="127.0.0.1"
export CI_STATEDB_TLS="false"
export CI_CHAIN_ID="ibc0"
export CI_STATEDB_PORT=5001
export CI_TOKENNAME="stack0"
export CI_HOME="testdata/ibc0"

./reset_base.sh
sed -i '' 's/timeout_broadcast_tx_commit = "10s"/timeout_broadcast_tx_commit = "60s"/g' $CI_HOME/config/config.toml


export CI_CHAIN_ID="ibc1"
export CI_STATEDB_PORT=5002
export CI_TOKENNAME="stack1"
export CI_HOME="testdata/ibc1"

./reset_base.sh

sed -i '' 's#"tcp://0.0.0.0:26657"#"tcp://0.0.0.0:26557"#g' $CI_HOME/config/config.toml
sed -i '' 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:26556"#g' $CI_HOME/config/config.toml
sed -i '' 's/timeout_broadcast_tx_commit = "10s"/timeout_broadcast_tx_commit = "60s"/g' $CI_HOME/config/config.toml



export CI_CHAIN_ID="ibc2"
export CI_STATEDB_PORT=5003
export CI_TOKENNAME="stack2"
export CI_HOME="testdata/ibc2"
rm -rf $CI_HOME

./reset_base.sh

sed -i '' 's#"tcp://0.0.0.0:26657"#"tcp://0.0.0.0:26457"#g' $CI_HOME/config/config.toml
sed -i '' 's#"tcp://0.0.0.0:26656"#"tcp://0.0.0.0:26456"#g' $CI_HOME/config/config.toml
sed -i '' 's/timeout_broadcast_tx_commit = "10s"/timeout_broadcast_tx_commit = "60s"/g' $CI_HOME/config/config.toml


#docker-compose -f bootstrap-docker/single-node.yaml up -d
docker-compose -f ../bootstrap-docker/three-node-db.yaml up -d

