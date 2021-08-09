
# Usage

### Local Dev

*  Compile binaries

> `make build`

* Init configs

> `./scripts/node_reset.sh`

* Start Raft LevelDB

> `docker-compose -f bootstrap-docker/single-node-db.yaml up -d`

* Run Node

>`./build/cid start --statedb_port=5002 --statedb_tls=false --home=scripts/testdata/ibc0`
>`./build/cli rest-server`


### Quick Start
`make release_build && docker-compose -f bootstrap-docker/single-node-db.yaml up -d`


This repository hosts Weelink, the implementation of the CI123Chain is a modification based on [Cosmos SDK](https://github.com/cosmos/cosmos-sdk).



