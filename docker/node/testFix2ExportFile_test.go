package node

import (
	"encoding/json"
	acc_types "github.com/ci123chain/ci123chain/pkg/account/types"
	app_types "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/vm/evmtypes"
	"github.com/tendermint/tendermint/types"
	"testing"
)

func TestFix2File(t *testing.T) {
	cdc := app_types.GetCodec()
	rawFile, err := types.GenesisDocFromFile("./genesis.json")

	exportFile, err := types.GenesisDocFromFile("/Users/tanhui/Downloads/genesis_173117.json")
	if err != nil {
		t.Log(err)
	}

	//saveAddressesStr, err := ioutil.ReadFile("./address.txt")
	//var address []string
	//err = json.Unmarshal(saveAddressesStr, &address)
	//if err != nil {
	//	panic(err)
	//}


	type GenesisState map[string]json.RawMessage
	var genesisStateRaw GenesisState
	cdc.MustUnmarshalJSON(rawFile.AppState, &genesisStateRaw)

	var genesisState GenesisState
	cdc.MustUnmarshalJSON(exportFile.AppState, &genesisState)




	var vmGenesisState, newVmGenesisState evmtypes.GenesisState
	evmtypes.ModuleCdc.MustUnmarshalJSON(genesisState["vm"], &vmGenesisState)
	evmtypes.ModuleCdc.MustUnmarshalJSON(genesisState["vm"], &newVmGenesisState)
	newVmGenesisState.Accounts = []evmtypes.GenesisAccount{}
	for _, vmAcc := range vmGenesisState.Accounts {
		newVmGenesisState.Accounts = append(newVmGenesisState.Accounts, vmAcc)
	}
	newVMStateBz, _ := json.Marshal(newVmGenesisState)
	genesisStateRaw["vm"] = newVMStateBz


	//fix accounts
	var accGenesisState, newAccGenesisState acc_types.GenesisState
	cdc.MustUnmarshalJSON(genesisState["accounts"], &accGenesisState)
	for _, gs := range accGenesisState {
		newAccGenesisState = append(newAccGenesisState, gs)
	}
	newAccGenesisStateBz := cdc.MustMarshalJSON(newAccGenesisState)

	genesisStateRaw["accounts"] = newAccGenesisStateBz
	//fix gravity
	//var graGenesisState, newGraGenesisState gravity_types.GenesisState
	//cdc.MustUnmarshalJSON(genesisState["gravity"], &graGenesisState)
	//newGraGenesisState = *gravity_types.DefaultGenesisState()
	//newGraGenesisState.DelegateKeys = graGenesisState.DelegateKeys
	//newGraGenesisStateBz := cdc.MustMarshalJSON(newGraGenesisState)
	//genesisState["gravity"] = newGraGenesisStateBz


	//write new file
	newGenesisStateBz := cdc.MustMarshalJSON(genesisStateRaw)
	rawFile.AppState = newGenesisStateBz
	rawFile.SaveAs("./exportFileNew.json")
}

