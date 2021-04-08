package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
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
