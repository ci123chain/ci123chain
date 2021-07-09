package module

import (
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/module"
	"github.com/ci123chain/ci123chain/pkg/vm/evmtypes"
	"github.com/ci123chain/ci123chain/pkg/vm/module/basic"
	"github.com/ci123chain/ci123chain/pkg/vm/moduletypes"
	wasm "github.com/ci123chain/ci123chain/pkg/vm/wasmtypes"
	abci "github.com/tendermint/tendermint/abci/types"
)

type AppModule struct {
	basic.AppModuleBasic
	Keeper  moduletypes.KeeperI
}

func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	WasmInitGenesis(ctx, am.Keeper)
	var genesisState evmtypes.GenesisState
	wasm.WasmCodec.MustUnmarshalJSON(data, &genesisState)
	EvmInitGenesis(ctx, am.Keeper, genesisState)
	return nil
}

// BeginBlock function for module at start of each block
func (am AppModule) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) {
	am.Keeper.BeginBlock(ctx, req)
}

// EndBlock function for module at end of block
func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return am.Keeper.EndBlock(ctx, req)
}

func (am AppModule) Committer(ctx sdk.Context) {
	//
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
}