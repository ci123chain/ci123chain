package supply

import (
	"encoding/json"
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
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


func (am AppModuleBasic) DefaultGenesis() json.RawMessage {
	return nil
}

func (am AppModuleBasic) Name() string {
	return ModuleName
}