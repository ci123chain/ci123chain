package order

import (
	"encoding/json"
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/order/keeper"
	"github.com/tanhuiya/ci123chain/pkg/order/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type AppModule struct {
	AppModuleBasic

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

type AppModuleBasic struct {

}

func (am AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	types.RegisterCodec(codec)
}

func (am AppModuleBasic) DefaultGenesis() json.RawMessage {
	return keeper.ModuleCdc.MustMarshalJSON(types.DefaultGenesisState())
}

func (am AppModuleBasic) Name() string {
	return ModuleName
}
/*
func RegisterCodec(cdc *codec.Codec)  {

}
*/

func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate  {
	if am.OrderKeeper.ExistOrderBook(ctx) {
		return nil
	}

	var genesisState types.GenesisState
	keeper.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
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