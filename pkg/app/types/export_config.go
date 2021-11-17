package types

type ConfigFiles struct {
	GenesisFile []byte `json:"genesis_file"`
	NodeID 		string `json:"node_id"`
}