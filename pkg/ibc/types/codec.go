package types

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
)

var IbcCdc *codec.Codec

func RegisterCodec(cdc *codec.Codec) {

	cdc.RegisterConcrete(&IBCTransfer{}, "ci123chain/IBCTransfer", nil)
	cdc.RegisterConcrete(&IBCInfo{}, "ci123chain/IBCInfo", nil)
	cdc.RegisterConcrete(&ApplyIBCTx{}, "ci123chain/ApplyIBCTx", nil)
	cdc.RegisterConcrete(&ApplyReceipt{}, "ci123chain/ApplyReceipt", nil)
	cdc.RegisterConcrete(&IBCMsgBankSend{}, "ci123chain/IBCMsgBankSend", nil)
	cdc.RegisterConcrete(&IBCReceiveReceiptMsg{}, "ci123chain/IBCReceiveReceiptMsg", nil)
}

func init()  {
	IbcCdc = codec.New()
	IbcCdc.RegisterInterface((*transaction.Transaction)(nil), nil)
	IbcCdc.RegisterConcrete(&transaction.CommonTx{}, "ci123chain/commontx", nil)
	RegisterCodec(IbcCdc)
	IbcCdc.Seal()
}
