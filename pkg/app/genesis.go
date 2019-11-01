package app

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/account"
	account_types "github.com/tanhuiya/ci123chain/pkg/account/types"
	"github.com/tanhuiya/ci123chain/pkg/auth"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	genesisBalance = 100
)

//// State to Unmarshal
//type GenesisState struct {
//	Accounts []account.Account `json:"accounts"`
//}

func AppGenStateJSON(cdc *amino.Codec, appGenTxs []json.RawMessage) (json.RawMessage, error) {
	appState := make(map[string]json.RawMessage)
	//err := AppGenStateAccount(cdc, appState, appGenTxs)
	//if err != nil {
	//	return nil, err
	//}
	//appState, err = json.Marshal(genesisState)

	AppGenstateAuth(cdc, appState)

	stateBytes, err := json.Marshal(appState)
	if err != nil {
		return nil, err
	}
	return stateBytes, nil
}

func AppGenstateAuth(cdc *amino.Codec, appState map[string]json.RawMessage)  {
	m := auth.AppModuleBasic{}
	appState[m.Name()] = m.DefaultGenesis()
}



func AppGenStateAccount(cdc *amino.Codec, appState map[string]json.RawMessage,
	genAccounts []account_types.GenesisAccount) (err error) {

	if len(genAccounts) == 0 {
		return
	}
	account.SetGenesisStateInAppState(cdc, appState, genAccounts)

	//
	//// get genesis flag account information
	//accounts := make([]account.Account, 0, len(appGenTxs))
	//accountm := make(map[common.Address]struct{})
	//
	//for _, appGenTx := range appGenTxs {
	//
	//	var genTx AppGenTx
	//	err = cdc.UnmarshalJSON(appGenTx, &genTx)
	//	if err != nil {
	//		return
	//	}
	//
	//	addr := common.HexToAddress(genTx.Address)
	//	if _, ok := accountm[addr]; !ok && genTx.Address != "" {
	//		accounts = append(accounts, account.Account{
	//			Address: addr,
	//			Amount:  genesisBalance,
	//		})
	//		accountm[addr] = struct{}{}
	//	}
	//}
	//
	//accountBytes, err := json.Marshal(accounts)
	//appState["accounts"] = accountBytes
	// create the final app state
	//genesisState = GenesisState{
	//	Accounts: accounts,
	//}
	return
}

func GetInitChainer(am account.AccountMapper) func(types.Context, abci.RequestInitChain) abci.ResponseInitChain {
	return func(ctx types.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		//stateJSON := req.AppStateBytes

		//var genesisState GenesisState
		//err := json.Unmarshal(stateJSON, &genesisState)
		//if err != nil {
		//	panic(err)
		//	// return sdk.ErrGenesisParse("").TraceCause(err, "")
		//}
		//
		//for _, acc := range genesisState.Accounts {
		//	if _, err := am.AddBalance(ctx, acc.Address, acc.Amount); err != nil {
		//		panic(err)
		//	}
		//	fmt.Printf("addr=%v amount=%v\n", acc.Address.Hex(), acc.Amount)
		//}
		return abci.ResponseInitChain{}
	}
}

func GenesisStateFromGenFile(cdc *codec.Codec, genFile string) (genesisState map[string]json.RawMessage, genDoc *tmtypes.GenesisDoc, err error)  {
	if !common.FileExist(genFile) {
		return genesisState, genDoc, fmt.Errorf("%s does not exist, run `init` first", genFile)
	}
	genDoc, err = tmtypes.GenesisDocFromFile(genFile)
	if err != nil {
		return genesisState, genDoc, err
	}
	genesisState, err = GenesisStateFromGenDoc(cdc, *genDoc)
	return genesisState, genDoc, nil
}

func GenesisStateFromGenDoc(cdc *codec.Codec, genDoc tmtypes.GenesisDoc,
) (genesisState map[string]json.RawMessage, err error) {
	if err = cdc.UnmarshalJSON(genDoc.AppState, &genesisState); err != nil {
		return genesisState, err
	}
	return genesisState, nil
}

func ExportGenesisFile(genDoc *tmtypes.GenesisDoc, genFile string) error {
	if err := genDoc.ValidateAndComplete(); err != nil {
		return err
	}
	return genDoc.SaveAs(genFile)
}
