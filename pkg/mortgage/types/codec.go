package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/transaction"
)

var MortgageCdc *codec.Codec

func RegisterCodec(cdc *codec.Codec) {

	cdc.RegisterConcrete(&MsgMortgage{}, "ci123chain/MsgMortgage", nil)
	cdc.RegisterConcrete(&MsgMortgageDone{}, "ci123chain/MsgMortgageDone", nil)
	cdc.RegisterConcrete(&MsgMortgageCancel{}, "ci123chain/MsgMortgageCancel", nil)
	cdc.RegisterConcrete(&Mortgage{}, "ci123chain/Mortgage", nil)
}

func init()  {
	MortgageCdc = codec.New()
	transaction.RegisterCodec(MortgageCdc)
	RegisterCodec(MortgageCdc)
	MortgageCdc.Seal()
}
