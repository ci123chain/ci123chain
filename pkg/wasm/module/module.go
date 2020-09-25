package module

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	wasm_types "github.com/ci123chain/ci123chain/pkg/wasm/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type AppModule struct {
	AppModuleBasic
	WasmKeeper  wasm_types.WasmKeeperI
}

func (am AppModule)InitGenesis(ctx sdk.Context, _ json.RawMessage) []abci.ValidatorUpdate {
	InitGenesis(ctx, am.WasmKeeper)
	return nil
}

func (am AppModule)BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) {
	//TODO
}

func (am AppModule)Committer(ctx sdk.Context) {
	//
}

func (am AppModule)EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	//
	return nil
}


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