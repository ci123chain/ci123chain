module chain_proxy

go 1.13

require github.com/ci123chain/ci123chain v1.4.20

replace (
    github.com/tendermint/tendermint => github.com/ci123chain/tendermint v0.32.7-rc41
    github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
)
