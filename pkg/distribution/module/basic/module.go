package basic

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/distribution"
	dtypes "github.com/ci123chain/ci123chain/pkg/distribution/types"
	tmtypes "github.com/tendermint/tendermint/types"
)


type AppModuleBasic struct {
}

func (am AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	dtypes.RegisterCodec(codec)
}


func (am AppModuleBasic) DefaultGenesis(validators []tmtypes.GenesisValidator) json.RawMessage {
	return distribution.ModuleCdc.MustMarshalJSON(dtypes.DefaultGenesisState(validators))
}

func (am AppModuleBasic) Name() string {
	return distribution.ModuleName
}

