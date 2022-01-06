package node

import (
	"encoding/json"
	acc_types "github.com/ci123chain/ci123chain/pkg/account/types"
	app_types "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/tendermint/tendermint/types"
	"io/ioutil"
	"testing"
)

func TestFixFile(t *testing.T) {
	cdc := app_types.GetCodec()
	var exportFile types.GenesisDoc
	exportFileBz, err := ioutil.ReadFile("./exportFile.json")
	if err != nil {
		t.Log(err)
	}
	json.Unmarshal(exportFileBz, &exportFile)
	type GenesisState map[string]json.RawMessage
	var genesisState GenesisState
	cdc.MustUnmarshalJSON(exportFile.AppState, &genesisState)
	var accGenesisState, newAccGenesisState acc_types.GenesisState
	cdc.MustUnmarshalJSON(genesisState["accounts"], &accGenesisState)
	for i := range accGenesisState {
		if accGenesisState[i].GetContractType() != "" {
			newAccGenesisState = append(newAccGenesisState, accGenesisState[i])
		}
	}
	newAccGenesisStateBz := cdc.MustMarshalJSON(newAccGenesisState)
	genesisState["accounts"] = newAccGenesisStateBz
	newGenesisStateBz := cdc.MustMarshalJSON(genesisState)
	exportFile.AppState = newGenesisStateBz
	exportFile.SaveAs("./exportFile2.json")
}