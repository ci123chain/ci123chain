package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/order/keeper"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&MsgUpgrade{}, "ci123chain/MsgUpgrade", nil)
	cdc.RegisterConcrete(&keeper.OrderBook{}, "ci123chain/OrderBook", nil)
}

func init()  {
	keeper.ModuleCdc = codec.New()
	RegisterCodec(keeper.ModuleCdc)
	codec.RegisterCrypto(keeper.ModuleCdc)
	keeper.ModuleCdc.Seal()
}
