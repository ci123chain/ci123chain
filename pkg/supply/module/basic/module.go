package basic

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/supply"

	types2 "github.com/ci123chain/ci123chain/pkg/supply/types"
	tmtypes "github.com/tendermint/tendermint/types"
)


type AppModuleBasic struct {
}

func (am AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	supply.RegisterCodec(codec)
}


func (am AppModuleBasic) DefaultGenesis(_ []tmtypes.GenesisValidator) json.RawMessage {
	var res = types2.DefaultGenesisState()
	return supply.ModuleCdc.MustMarshalJSON(res)

}

func (am AppModuleBasic) Name() string {
	return supply.ModuleName
}