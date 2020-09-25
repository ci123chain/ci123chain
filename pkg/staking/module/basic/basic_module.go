package basic

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
	tmtypes "github.com/tendermint/tendermint/types"
)


type AppModuleBasic struct {}

func (am AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	types.RegisterCodec(codec)
}

func (am AppModuleBasic) DefaultGenesis(validators []tmtypes.GenesisValidator) json.RawMessage {
	p := types.DefaultGenesisState(validators)
	b, err := types.StakingCodec.MarshalJSONIndent(p, "", "")
	if err != nil{
		panic(err)
	}
	return b
}

func (am AppModuleBasic) Name() string {
	return types.ModuleName
}