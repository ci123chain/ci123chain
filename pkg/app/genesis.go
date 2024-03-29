package app

import (
	"encoding/json"
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/app/module"
	"github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
	tmabci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	INITCHAINKEY  = "isInitChain"
	FIRSTSHARD    = "1"
)
func AppGenStateJSON(validators []tmtypes.GenesisValidator) (json.RawMessage, error) {
	appState := module.ModuleBasics.DefaultGenesis(validators)
	stateBytes, err := json.Marshal(appState)
	if err != nil {
		return nil, err
	}
	return stateBytes, nil
}


type GenesisState map[string]json.RawMessage

func (c *Chain) InitChainer (ctx sdk.Context, req tmabci.RequestInitChain) tmabci.ResponseInitChain {
	c.BaseApp.Logger.Info("Begin initChain")
	index := viper.GetString(flagShardIndex)
	if len(index) != 0 && index != FIRSTSHARD {
		return tmabci.ResponseInitChain{}
	} else {
		var genesisState GenesisState
		store := ctx.KVStore(keys[StoreKey])
		isInitChain := store.Get([]byte(INITCHAINKEY))
		if isInitChain != nil {
			return tmabci.ResponseInitChain{}
		} else {
			c.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
			store.Set([]byte(INITCHAINKEY), []byte("true"))
		}
		c.BaseApp.Logger.Info("initChain Finish")
		return c.mm.InitGenesis(ctx, genesisState)
	}
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
		return genesisState, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error())
	}
	return genesisState, nil
}

func ExportGenesisFile(genDoc *tmtypes.GenesisDoc, genFile string) error {
	if err := genDoc.ValidateAndComplete(); err != nil {
		return types.ErrGenesisFile(types.DefaultCodespace, err)
	}
	return genDoc.SaveAs(genFile)
}
