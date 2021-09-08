package module

import (
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/module"
	ak "github.com/ci123chain/ci123chain/pkg/account/keeper"
	"github.com/ci123chain/ci123chain/pkg/staking"
	k "github.com/ci123chain/ci123chain/pkg/staking/keeper"
	"github.com/ci123chain/ci123chain/pkg/staking/module/basic"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
	sk "github.com/ci123chain/ci123chain/pkg/supply/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
)

type AppModule struct {
	basic.AppModuleBasic
	StakingKeeper k.StakingKeeper
	AccountKeeper ak.AccountKeeper
	SupplyKeeper  sk.Keeper
}

func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	//
	var genesisState types.GenesisState
	err := types.StakingCodec.UnmarshalJSON(data, &genesisState)
	if err != nil {
		panic(err)
	}
	return staking.InitGenesis(ctx, am.StakingKeeper, am.AccountKeeper, am.SupplyKeeper, genesisState)
}

func (am AppModule) BeginBlocker(ctx sdk.Context, _ abci.RequestBeginBlock) {
	staking.BeginBlock(ctx, am.StakingKeeper)
}

func (am AppModule) Committer(ctx sdk.Context) {
	//
}

func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	return staking.EndBlock(ctx, am.StakingKeeper)
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
}

func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return types.StakingCodec.MustMarshalJSON(staking.ExportGenesis(ctx, am.StakingKeeper))
}