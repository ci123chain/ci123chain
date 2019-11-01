package app

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/account"
	"github.com/tanhuiya/ci123chain/pkg/auth"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	genesisBalance = 100
)


func AppGenStateJSON(cdc *amino.Codec, appGenTxs []json.RawMessage) (json.RawMessage, error) {
	appState := make(map[string]json.RawMessage)

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


type GenesisState map[string]json.RawMessage 

func GetInitChainer(accModule account.AppModule) func(types.Context, abci.RequestInitChain) abci.ResponseInitChain {
	return func(ctx types.Context, req abci.RequestInitChain) abci.ResponseInitChain {

		stateJSON := req.AppStateBytes
		var genesisState GenesisState
		err := json.Unmarshal(stateJSON, &genesisState)
		if err != nil {
			panic(err)
			// return sdk.ErrGenesisParse("").TraceCause(err, "")
		}

		// 单独设置 account module
		accModule.InitGenesis(ctx, genesisState[account.ModuleName])
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
