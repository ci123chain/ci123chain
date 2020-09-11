package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
)

var WasmCodec  *codec.Codec

func init() {
	WasmCodec = codec.New()
	RegisterCodec(WasmCodec)
	codec.RegisterCrypto(WasmCodec)
	WasmCodec.Seal()
}

func RegisterCodec(cdc *codec.Codec) {
	//
	cdc.RegisterConcrete(&MsgInstantiateContract{}, "ci123Chain/MsgInstantiateContract", nil)
	cdc.RegisterConcrete(&MsgExecuteContract{}, "ci123Chain/MsgExecuteContract", nil)
	cdc.RegisterConcrete(&MsgMigrateContract{}, "ci123Chain/MsgMigrateContract", nil)
}