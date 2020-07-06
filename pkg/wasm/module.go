package wasm

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	wasm "github.com/ci123chain/ci123chain/pkg/wasm/keeper"
	"github.com/ci123chain/ci123chain/pkg/wasm/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type AppModule struct {
	AppModuleBasic
	WasmKeeper  wasm.Keeper
}

func (am AppModule)InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
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
	types.RegisterCodec(codec)
}

func (am AppModuleBasic) DefaultGenesis(_ []tmtypes.GenesisValidator, _ []string) json.RawMessage {
	return nil
}


func (am AppModuleBasic) Name() string {
	return ModuleName
}