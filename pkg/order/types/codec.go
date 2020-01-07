package types

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	"github.com/tanhuiya/ci123chain/pkg/order/keeper"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&UpgradeTx{}, "ci123chain/UpgradeTx", nil)
	cdc.RegisterConcrete(&keeper.OrderBook{}, "ci123chain/OrderBook", nil)
}

func init()  {
	keeper.ModuleCdc = codec.New()
	RegisterCodec(keeper.ModuleCdc)
	codec.RegisterCrypto(keeper.ModuleCdc)
	keeper.ModuleCdc.Seal()
}
