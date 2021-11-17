package module

import (
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/module"
	ak "github.com/ci123chain/ci123chain/pkg/account/keeper"
	sk "github.com/ci123chain/ci123chain/pkg/supply/keeper"
	"github.com/ci123chain/ci123chain/pkg/upgrade"
	k "github.com/ci123chain/ci123chain/pkg/upgrade/keeper"
	"github.com/ci123chain/ci123chain/pkg/upgrade/module/basic"
	abci "github.com/tendermint/tendermint/abci/types"
)

type AppModule struct {
	basic.AppModuleBasic
	UpgradeKeeper k.Keeper
	AccountKeeper ak.AccountKeeper
	SupplyKeeper  sk.Keeper
}

func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	return nil
}

func (am AppModule) BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock) {
	upgrade.BeginBlock(am.UpgradeKeeper, ctx)
}

func (am AppModule) Committer(ctx sdk.Context) {
	//
}

func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return nil
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
}

func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return nil
}