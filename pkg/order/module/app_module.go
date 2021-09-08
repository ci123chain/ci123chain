package module

import (
	"encoding/json"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/abci/types/module"
	"github.com/ci123chain/ci123chain/pkg/order/keeper"
	"github.com/ci123chain/ci123chain/pkg/order/module/basic"
	"github.com/ci123chain/ci123chain/pkg/order/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type AppModule struct {
	basic.AppModuleBasic

	OrderKeeper	*keeper.OrderKeeper
}

func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {

	return nil
}

func (am AppModule) BeginBlocker(ctx  sdk.Context, req abci.RequestBeginBlock) {
	//do you want to do
	am.OrderKeeper.WaitForReady(ctx)
}

func (am AppModule) Committer(ctx sdk.Context) {
	//do you want to do
}


func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate  {
	if am.OrderKeeper.ExistOrderBook(ctx) {
		return nil
	}
	//am.OrderKeeper.Cdb.ResetDB()
	var genesisState types.GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	shardID := ctx.ChainID()
	if genesisState.Params.OrderBook.Lists != nil {
		name := genesisState.Params.OrderBook.Lists[0].Name
		if name == "" {
			genesisState.Params.OrderBook.Lists[0].Name = shardID
			am.OrderKeeper.SetOrderBook(ctx, genesisState.Params.OrderBook)
		}else if name == shardID {
			am.OrderKeeper.SetOrderBook(ctx, genesisState.Params.OrderBook)
		}
	}
	return nil
}

func (am AppModule) RegisterServices(cfg module.Configurator) {
}

func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	ob, _ := am.OrderKeeper.GetOrderBook(ctx)
	gs := types.NewGenesisState(types.Params{OrderBook:ob})
	return types.ModuleCdc.MustMarshalJSON(gs)
}