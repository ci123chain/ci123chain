package order

import (
	"encoding/json"
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/order/keeper"
	"github.com/tanhuiya/ci123chain/pkg/order/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

var OrderCdc *codec.Codec
var ModuleCdc = OrderCdc
const ModuleName  = "order"

type AppModule struct {
	AppModuleBasic

	OrderKeeper	*keeper.OrderKeeper
}

func (am AppModule) BeginBlocker(ctx  sdk.Context, req abci.RequestBeginBlock) {
	//do you want to do
	am.OrderKeeper.WaitForReady(ctx.ChainID(),ctx.BlockHeight())
}

func (am AppModule) Committer(ctx sdk.Context) {
	//do you want to do
	if am.OrderKeeper.IsDeal {
		rev, orderBook := am.OrderKeeper.GetOrderBook()
		err := am.OrderKeeper.UpdateOrderBook(orderBook, rev, orderBook.Lists[orderBook.Current.Index].Name, orderBook.Lists[orderBook.Current.Index].Height, keeper.StateDone)
		if err != nil {
			panic(err)
		}
		am.OrderKeeper.IsDeal = false
	}
}

type AppModuleBasic struct {

}

func (am AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	RegisterCodec(codec)
}

func (am AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(types.DefaultGenesisState())
}

func (am AppModuleBasic) Name() string {
	return ModuleName
}
/*
func RegisterCodec(cdc *codec.Codec)  {

}
*/
func init()  {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
}

func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage)  {
	var genesisState types.GenesisState
	ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.OrderKeeper, genesisState)
}

func InitGenesis(ctx sdk.Context, ok *keeper.OrderKeeper, data types.GenesisState) {
	ok.SetOrderBook(data.Params.OrderBook)
}