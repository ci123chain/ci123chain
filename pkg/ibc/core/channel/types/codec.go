package types

import "github.com/ci123chain/ci123chain/pkg/abci/codec"

var ChannelCdc *codec.Codec

func init()  {
	ChannelCdc = codec.New()
	RegisterCodec(ChannelCdc)
	codec.RegisterCrypto(ChannelCdc)
	ChannelCdc.Seal()
}

func RegisterCodec(cdc *codec.Codec)  {
	cdc.RegisterConcrete(&Acknowledgement{}, "ibcChannel/Acknowledgement", nil)
	cdc.RegisterConcrete(&MsgChannelOpenInit{}, "ibcChannel/MsgChannelOpenInit", nil)
	cdc.RegisterConcrete(&MsgChannelOpenTry{}, "ibcChannel/MsgChannelOpenTry", nil)
	cdc.RegisterConcrete(&MsgChannelOpenAck{}, "ibcChannel/MsgChannelOpenAck", nil)
	cdc.RegisterConcrete(&MsgChannelOpenConfirm{}, "ibcChannel/MsgChannelOpenConfirm", nil)
	cdc.RegisterConcrete(&MsgRecvPacket{}, "ibcChannel/MsgRecvPacket", nil)
	cdc.RegisterConcrete(&MsgAcknowledgement{}, "ibcChannel/MsgAcknowledgement", nil)
	cdc.RegisterConcrete(&MsgTimeout{}, "ibcChannel/MsgTimeout", nil)
	cdc.RegisterConcrete(&Channel{}, "ibcChannel/Channel", nil)
	cdc.RegisterConcrete(&IdentifiedChannel{}, "ibcChannel/IdentifiedChannel", nil)
}
