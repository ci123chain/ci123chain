package types

type GenesisState struct {
	Data   []StoredContent   `json:"data"`
}



func DefaultGenesisState() GenesisState {
	return GenesisState{}
}