package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
)

var ModuleCdc *codec.Codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&MsgUpgrade{}, "ci123chain/MsgUpgrade", nil)
	cdc.RegisterConcrete(&OrderBook{}, "ci123chain/OrderBook", nil)
}

func init()  {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
