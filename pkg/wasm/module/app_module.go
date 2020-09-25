package module

import (
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	wasm_types "github.com/ci123chain/ci123chain/pkg/wasm/types"
	"github.com/ci123chain/ci123chain/pkg/wasm/module/basic"
	abci "github.com/tendermint/tendermint/abci/types"
)

type AppModule struct {
	basic.AppModuleBasic
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

