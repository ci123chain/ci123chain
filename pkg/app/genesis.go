package app

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func AppGenStateJSON() (json.RawMessage, error) {
	appState := ModuleBasics.DefaultGenesis()
	stateBytes, err := json.Marshal(appState)
	if err != nil {
		return nil, err
	}
	return stateBytes, nil
}


type GenesisState map[string]json.RawMessage

func (c *Chain)InitChainer (ctx types.Context, req abci.RequestInitChain) abci.ResponseInitChain {

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
