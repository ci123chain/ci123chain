package node

import (
	"encoding/json"
	acc_types "github.com/ci123chain/ci123chain/pkg/account/types"
	app_types "github.com/ci123chain/ci123chain/pkg/app/types"
	gravity_types "github.com/ci123chain/ci123chain/pkg/gravity/types"
	"github.com/tendermint/tendermint/types"
	"testing"
)

func TestFixFile(t *testing.T) {
	cdc := app_types.GetCodec()
	exportFile, err := types.GenesisDocFromFile("./exportFile.json")
	if err != nil {
		t.Log(err)
	}
	type GenesisState map[string]json.RawMessage
	var genesisState GenesisState
	cdc.MustUnmarshalJSON(exportFile.AppState, &genesisState)

	//fix accounts
	var accGenesisState, newAccGenesisState acc_types.GenesisState
	cdc.MustUnmarshalJSON(genesisState["accounts"], &accGenesisState)
	for i := range accGenesisState {
		if accGenesisState[i].GetContractType() != "" {
			newAccGenesisState = append(newAccGenesisState, accGenesisState[i])
		}
	}
	newAccGenesisStateBz := cdc.MustMarshalJSON(newAccGenesisState)
	genesisState["accounts"] = newAccGenesisStateBz

	//fix gravity
	var graGenesisState, newGraGenesisState gravity_types.GenesisState
	cdc.MustUnmarshalJSON(genesisState["gravity"], &graGenesisState)
	newGraGenesisState = *gravity_types.DefaultGenesisState()
	newGraGenesisState.DelegateKeys = graGenesisState.DelegateKeys
	newGraGenesisStateBz := cdc.MustMarshalJSON(newGraGenesisState)
	genesisState["gravity"] = newGraGenesisStateBz

	//write new file
	newGenesisStateBz := cdc.MustMarshalJSON(genesisState)
	exportFile.AppState = newGenesisStateBz
	exportFile.SaveAs("./exportFile3.json")
}