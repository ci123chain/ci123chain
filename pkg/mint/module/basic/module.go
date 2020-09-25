package basic

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/mint"
	tmtypes "github.com/tendermint/tendermint/types"
)

type AppModuleBasic struct {}

func (am AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	mint.RegisterCodec(codec)
}


func (am AppModuleBasic) DefaultGenesis(_ []tmtypes.GenesisValidator) json.RawMessage {
	return mint.ModuleCdc.MustMarshalJSON(mint.DefaultGenesisState())
}

func (am AppModuleBasic) Name() string {
	return mint.ModuleName
}

