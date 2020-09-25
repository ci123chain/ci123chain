package basic

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	wasm_types "github.com/ci123chain/ci123chain/pkg/wasm/types"
	tmtypes "github.com/tendermint/tendermint/types"
)


type AppModuleBasic struct {}


func (am AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	wasm_types.RegisterCodec(codec)
}

func (am AppModuleBasic) DefaultGenesis(vals []tmtypes.GenesisValidator) json.RawMessage {
	return nil
}


func (am AppModuleBasic) Name() string {
	return wasm_types.ModuleName
}