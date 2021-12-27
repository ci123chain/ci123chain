package types

type ConfigFiles struct {
	GenesisFile []byte `json:"genesis_file"`
	NodeID 		string `json:"node_id"`
	ETHChainID  uint64  `json:"eth_chain_id"`
}