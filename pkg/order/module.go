package order

import (
	"encoding/json"
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/order/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
)

var OrderCdc *codec.Codec
var ModuleCdc = OrderCdc
const ModuleName  = "order"

type AppModule struct {
	AppModuleBasic

	OrderKeeper	*keeper.OrderKeeper
}

func (am AppModule) BeginBlocker(ctx types.Context, req abci.RequestBeginBlock) {
	//do you want to do
	am.OrderKeeper.WaitForReady(ctx.ChainID(),ctx.BlockHeight())
}

func (am AppModule) Committer(ctx types.Context) {
	//do you want to do
	if am.OrderKeeper.IsDeal {
		rev, orderBook := am.OrderKeeper.GetOrderBook()
		err := am.OrderKeeper.UpdateOrderBook(orderBook.TotalShards, rev, orderBook.ShardNow, orderBook.HeightNow, keeper.StateDone)
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
	return nil
}

func (am AppModuleBasic) Name() string {
	return ModuleName
}

func RegisterCodec(cdc *codec.Codec)  {

}

func init()  {
	OrderCdc = codec.New()
	RegisterCodec(OrderCdc)
	codec.RegisterCrypto(OrderCdc)
	OrderCdc.Seal()
}

func (am AppModule) InitGenesis(ctx types.Context, data json.RawMessage)  {
	//do something
}