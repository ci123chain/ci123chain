package types

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
)



type GenesisState GenesisAccounts

func NewGenesisState(accounts GenesisAccounts) GenesisState {
	return GenesisState(accounts)
}

func SetGenesisStateInAppState(cdc *codec.Codec, appState map[string]json.RawMessage,
	genesisState GenesisState)  map[string]json.RawMessage  {

	genesisSteteBz := cdc.MustMarshalJSON(genesisState)
	appState[ModuleName] = genesisSteteBz
	return appState
}
