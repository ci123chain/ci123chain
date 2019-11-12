package types

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
)

var IbcCdc *codec.Codec

func RegisterCodec(cdc *codec.Codec) {

	cdc.RegisterConcrete(&IBCTransfer{}, "ci123chain/IBCTransfer", nil)
	cdc.RegisterConcrete(&IBCMsg{}, "ci123chain/IBCMsg", nil)
	cdc.RegisterConcrete(&ApplyIBCTx{}, "ci123chain/ApplyIBCTx", nil)
	cdc.RegisterConcrete(&SignedIBCMsg{}, "ci123chain/SignedIBCMsg", nil)
}

func init()  {
	IbcCdc = codec.New()
	IbcCdc.RegisterInterface((*transaction.Transaction)(nil), nil)
	IbcCdc.RegisterConcrete(&transaction.CommonTx{}, "ci123chain/commontx", nil)
	RegisterCodec(IbcCdc)
	IbcCdc.Seal()
}
