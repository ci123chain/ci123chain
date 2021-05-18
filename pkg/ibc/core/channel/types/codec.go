package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	codectypes "github.com/ci123chain/ci123chain/pkg/abci/codec/types"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
)

var ChannelCdc *codec.Codec

func init()  {
	ChannelCdc = codec.New()
	RegisterCodec(ChannelCdc)
	codec.RegisterCrypto(ChannelCdc)
	ChannelCdc.Seal()
}

var SubModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())

func RegisterCodec(cdc *codec.Codec)  {
	//cdc.RegisterInterface((*IsAcknowledgement_Response)(nil), nil)
	cdc.RegisterInterface((*isAcknowledgement_Response)(nil), nil)
	//
	cdc.RegisterConcrete(&Acknowledgement{}, "ibcChannel/Acknowledgement", nil)
	cdc.RegisterConcrete(&Acknowledgement_Error{}, "ibcChannel/AcknowledgementError", nil)
	cdc.RegisterConcrete(&Acknowledgement_Result{}, "ibcChannel/AcknowledgementResult", nil)
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

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterInterface(
		"ibc.core.channel.v1.ChannelI",
		(*exported.ChannelI)(nil),
	)
	registry.RegisterInterface(
		"ibc.core.channel.v1.CounterpartyChannelI",
		(*exported.CounterpartyChannelI)(nil),
	)
	registry.RegisterInterface(
		"ibc.core.channel.v1.PacketI",
		(*exported.PacketI)(nil),
	)
	registry.RegisterImplementations(
		(*exported.ChannelI)(nil),
		&Channel{},
	)
	registry.RegisterImplementations(
		(*exported.CounterpartyChannelI)(nil),
		&Counterparty{},
	)
	registry.RegisterImplementations(
		(*exported.PacketI)(nil),
		&Packet{},
	)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgChannelOpenInit{},
		&MsgChannelOpenTry{},
		&MsgChannelOpenAck{},
		&MsgChannelOpenConfirm{},
		//&MsgChannelCloseInit{},
		//&MsgChannelCloseConfirm{},
		&MsgRecvPacket{},
		&MsgAcknowledgement{},
		&MsgTimeout{},
		//&MsgTimeoutOnClose{},
	)

}