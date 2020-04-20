package supply

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type AppModule struct {
	AppModuleBasic
}


func (am AppModule) InitGenesis(ctx types.Context, data json.RawMessage)  {
}


type AppModuleBasic struct {
}

func (am AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	RegisterCodec(codec)
}


func (am AppModuleBasic) DefaultGenesis(_ []tmtypes.GenesisValidator) json.RawMessage {
	return nil
}

func (am AppModuleBasic) Name() string {
	return ModuleName
}