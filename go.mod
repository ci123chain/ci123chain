module chain_proxy

go 1.13

require github.com/ci123chain/ci123chain v1.5.3-0.20210621100026-7c17d913c8d3

replace (
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/tendermint/tendermint => github.com/ci123chain/tendermint v0.32.7-rc44
)
