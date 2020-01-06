package order

import (
	"encoding/json"
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/order/keeper"
	"github.com/tanhuiya/ci123chain/pkg/order/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const ModuleName  = "order"

type AppModule struct {
	AppModuleBasic

	OrderKeeper	*keeper.OrderKeeper
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
	RegisterCodec(codec)
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

func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage)  {
	store := ctx.KVStore(am.OrderKeeper.StoreKey)
	if store.Has([]byte(keeper.OrderBookKey)) {
		return
	}
	var genesisState types.GenesisState
	keeper.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	shardID := ctx.ChainID()
	if genesisState.Params.OrderBook.Lists != nil && genesisState.Params.OrderBook.Lists[0].Name == ""{
		genesisState.Params.OrderBook.Lists[0].Name = shardID
	}
	bz, err := keeper.ModuleCdc.MarshalBinaryLengthPrefixed(genesisState.Params.OrderBook)
	if err != nil {
		panic(err)
	}
	store.Set([]byte(keeper.OrderBookKey), bz)
}