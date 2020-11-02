package basic


import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/infrastructure"
	"github.com/ci123chain/ci123chain/pkg/infrastructure/types"
	tmtypes "github.com/tendermint/tendermint/types"
)


type AppModuleBasic struct {

}

func (am AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	types.RegisterCodec(codec)
}

func (am AppModuleBasic) DefaultGenesis(_ []tmtypes.GenesisValidator) json.RawMessage {
	return infrastructure.ModuleCdc.MustMarshalJSON(types.DefaultGenesisState())
}

func (am AppModuleBasic) Name() string {
	return infrastructure.ModuleName
}
