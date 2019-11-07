package types

import "github.com/tanhuiya/ci123chain/pkg/abci/codec"

var MortgageCdc *codec.Codec

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgMortgage{}, "ci123chain/MsgMortgage", nil)
	cdc.RegisterConcrete(MsgMortgageDone{}, "ci123chain/MsgMortgageDone", nil)
	cdc.RegisterConcrete(MsgMortgageCancel{}, "ci123chain/MsgMortgageCancel", nil)
	cdc.RegisterConcrete(Mortgage{}, "ci123chain/Mortgage", nil)
}

func init()  {
	MortgageCdc = codec.New()
	RegisterCodec(MortgageCdc)
	MortgageCdc.Seal()
}
