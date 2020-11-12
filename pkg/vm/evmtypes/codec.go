package evmtypes

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
)

// ModuleCdc defines the evm module's codec
var ModuleCdc *codec.Codec

// RegisterCodec registers all the necessary types and interfaces for the
// evm module
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&MsgEvmTx{}, "ethermint/MsgEvmTx", nil)
	cdc.RegisterConcrete(&TxData{}, "ethermint/TxData", nil)
	cdc.RegisterConcrete(&ChainConfig{}, "ethermint/ChainConfig", nil)
}

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
