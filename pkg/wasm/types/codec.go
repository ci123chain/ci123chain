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
	cdc.RegisterConcrete(&InstantiateContractTx{}, "ci123Chain/InstantiateContractTx", nil)
	cdc.RegisterConcrete(&ExecuteContractTx{}, "ci123Chain/ExecuteContractTx", nil)
	cdc.RegisterConcrete(&MigrateContractTx{}, "ci123Chain/MigrateContractTx", nil)
}