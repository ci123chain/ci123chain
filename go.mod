module chain_proxy

go 1.13

require github.com/ci123chain/ci123chain v1.5.0-ibc-beta2.0.20210610014127-8947dade64b3

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/tendermint/tendermint => github.com/ci123chain/tendermint v0.32.7-rc41
)
