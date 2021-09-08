package module

import (
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/module"
	"github.com/ci123chain/ci123chain/pkg/account"
	"github.com/ci123chain/ci123chain/pkg/vm/evmtypes"
	"github.com/ci123chain/ci123chain/pkg/vm/module/basic"
	"github.com/ci123chain/ci123chain/pkg/vm/moduletypes"
	abci "github.com/tendermint/tendermint/abci/types"
)

type AppModule struct {
	basic.AppModuleBasic
	Keeper  moduletypes.KeeperI
	AccountKeeper  account.AccountKeeper
}

func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	//WasmInitGenesis(ctx, am.Keeper)
	var genesisState evmtypes.GenesisState
	//wasm.WasmCodec.MustUnmarshalJSON(data, &genesisState)
	_ = json.Unmarshal(data, &genesisState)
	EvmInitGenesis(ctx, am.Keeper, am.AccountKeeper, genesisState)
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

func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	//return wasm.WasmCodec.MustMarshalJSON(ExportGenesis(ctx, am.Keeper, am.AccountKeeper))
	by,_ := json.Marshal(ExportGenesis(ctx, am.Keeper, am.AccountKeeper))
	return by
}