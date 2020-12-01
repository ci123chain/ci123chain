package module

import (
	"encoding/json"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account"
	"github.com/ci123chain/ci123chain/pkg/distribution"
	"github.com/ci123chain/ci123chain/pkg/distribution/module/basic"
	"github.com/ci123chain/ci123chain/pkg/supply"

	k "github.com/ci123chain/ci123chain/pkg/distribution/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
)

type AppModule struct {
	basic.AppModuleBasic

	DistributionKeeper  k.DistrKeeper
	AccountKeeper      account.AccountKeeper
	SupplyKeeper       supply.Keeper
}

func (am AppModule) EndBlock(ctx types.Context, req abci.RequestEndBlock) ([]abci.ValidatorUpdate, []abci.Event) {
	return nil, nil
}

func (am AppModule) BeginBlocker(ctx types.Context, req abci.RequestBeginBlock) {
	distribution.BeginBlock(ctx, req, am.DistributionKeeper)
}

func (am AppModule) InitGenesis(ctx types.Context, data json.RawMessage) []abci.ValidatorUpdate {

	var genesisState distribution.GenesisState
	distribution.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	distribution.InitGenesis(ctx, am.AccountKeeper, am.SupplyKeeper, am.DistributionKeeper, genesisState)
	return nil
}
func (am AppModule) Committer(ctx types.Context) {

}
