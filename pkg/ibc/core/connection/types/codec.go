package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	codectypes "github.com/ci123chain/ci123chain/pkg/abci/codec/types"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
)

var IBCConnectionCodec *codec.Codec

func init(){
	IBCConnectionCodec = codec.New()
	RegisterCodec(IBCConnectionCodec)
	codec.RegisterCrypto(IBCConnectionCodec)
	IBCConnectionCodec.Seal()
}
func RegisterCodec(cdc *codec.Codec)  {
	cdc.RegisterConcrete(&MsgConnectionOpenInit{}, "ibcConnection/MsgConnectionOpenInit", nil)
	cdc.RegisterConcrete(&MsgConnectionOpenTry{}, "ibcConnection/MsgConnectionOpenTry", nil)
	cdc.RegisterConcrete(&MsgConnectionOpenAck{}, "ibcConnection/MsgConnectionOpenAck", nil)
	cdc.RegisterConcrete(&MsgConnectionOpenConfirm{}, "ibcConnection/MsgConnectionOpenConfirm", nil)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterInterface(
		"ibc.core.connection.v1.ConnectionI",
		(*exported.ConnectionI)(nil),
		&ConnectionEnd{},
	)
	registry.RegisterInterface(
		"ibc.core.connection.v1.CounterpartyConnectionI",
		(*exported.CounterpartyConnectionI)(nil),
		&Counterparty{},
	)
	registry.RegisterInterface(
		"ibc.core.connection.v1.Version",
		(*exported.Version)(nil),
		&Version{},
	)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgConnectionOpenInit{},
		&MsgConnectionOpenTry{},
		&MsgConnectionOpenAck{},
		&MsgConnectionOpenConfirm{},
	)

	//msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}