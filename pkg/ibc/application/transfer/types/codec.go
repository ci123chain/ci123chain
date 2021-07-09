package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	codectypes "github.com/ci123chain/ci123chain/pkg/abci/codec/types"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

var IBCTransferCdc *codec.Codec

func init()  {
	IBCTransferCdc = codec.New()
	RegisterCodec(IBCTransferCdc)
	codec.RegisterCrypto(IBCTransferCdc)
	IBCTransferCdc.Seal()
}

var ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())

func RegisterCodec(cdc *codec.Codec)  {
	cdc.RegisterConcrete(&MsgTransfer{}, "ibcTransfer/msgTransfer", nil)
	cdc.RegisterConcrete(&MsgTransferResponse{}, "ibcTransfer/msgTransferResponse", nil)
	cdc.RegisterConcrete(&FungibleTokenPacketData{}, "ibcTransfer/fungibleTokenPacketData", nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgTransfer{})
}