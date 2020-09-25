package app

import (
	"encoding/json"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/app/module"
	"github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ethereum/go-ethereum/common"
	tmabci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func AppGenStateJSON(validators []tmtypes.GenesisValidator) (json.RawMessage, error) {
	appState := module.ModuleBasics.DefaultGenesis(validators)
	stateBytes, err := json.Marshal(appState)
	if err != nil {
		return nil, abci.ErrInternal("Marshal failed")
	}
	return stateBytes, nil
}


type GenesisState map[string]json.RawMessage

func (c *Chain) InitChainer (ctx sdk.Context, req tmabci.RequestInitChain) tmabci.ResponseInitChain {
	var genesisState GenesisState
	c.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	return c.mm.InitGenesis(ctx, genesisState)
}


func GenesisStateFromGenFile(cdc *codec.Codec, genFile string) (genesisState map[string]json.RawMessage, genDoc *tmtypes.GenesisDoc, err error)  {
	if !common.FileExist(genFile) {
		return genesisState, genDoc, fmt.Errorf("%s does not exist, run `init` first", genFile)
	}
	genDoc, err = tmtypes.GenesisDocFromFile(genFile)
	if err != nil {
		return genesisState, genDoc, types.ErrGenesisFile(types.DefaultCodespace, err)
	}
	genesisState, err = GenesisStateFromGenDoc(cdc, *genDoc)
	return genesisState, genDoc, nil
}

func GenesisStateFromGenDoc(cdc *codec.Codec, genDoc tmtypes.GenesisDoc,
) (genesisState map[string]json.RawMessage, err error) {
	if err = cdc.UnmarshalJSON(genDoc.AppState, &genesisState); err != nil {
		return genesisState, abci.ErrInternal("Unmarshal failed")
	}
	return genesisState, nil
}

func ExportGenesisFile(genDoc *tmtypes.GenesisDoc, genFile string) error {
	if err := genDoc.ValidateAndComplete(); err != nil {
		return types.ErrGenesisFile(types.DefaultCodespace, err)
	}
	return genDoc.SaveAs(genFile)
}
