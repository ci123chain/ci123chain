package types

import "github.com/ci123chain/ci123chain/pkg/abci/codec"

var StakingCodec *codec.Codec

func init(){
	StakingCodec = codec.New()
	RegisterCodec(StakingCodec)
	codec.RegisterCrypto(StakingCodec)
	StakingCodec.Seal()
}

func RegisterCodec(cdc *codec.Codec) {
	//
	cdc.RegisterConcrete(&CreateValidatorTx{}, "ci123chain/CreateValidatorTx", nil)
	cdc.RegisterConcrete(&EditValidatorTx{}, "ci123chain/EditValidatorTx", nil)
	cdc.RegisterConcrete(&DelegateTx{}, "ci123chain/DelegateTx", nil)
	cdc.RegisterConcrete(&UndelegateTx{}, "ci123chain/UndelegateTx", nil)
	cdc.RegisterConcrete(&RedelegateTx{}, "ci123chain/RedelegateTx", nil)
}