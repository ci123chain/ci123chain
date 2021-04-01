package types

import "github.com/ci123chain/ci123chain/pkg/abci/codec"

var IBCTransferCdc *codec.Codec

func init()  {
	IBCTransferCdc = codec.New()
	RegisterCodec(IBCTransferCdc)
	codec.RegisterCrypto(IBCTransferCdc)
	IBCTransferCdc.Seal()
}

func RegisterCodec(cdc *codec.Codec)  {
	cdc.RegisterConcrete(&MsgTransfer{}, "ibcTransfer/msgTransfer", nil)
	cdc.RegisterConcrete(&MsgTransferResponse{}, "ibcTransfer/msgTransferResponse", nil)
	cdc.RegisterConcrete(&FungibleTokenPacketData{}, "ibcTransfer/fungibleTokenPacketData", nil)
}