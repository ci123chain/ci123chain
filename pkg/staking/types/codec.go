package types

import "github.com/tanhuiya/ci123chain/pkg/abci/codec"

var StakingCodec *codec.Codec

func init(){
	StakingCodec = codec.New()
	RegisterCodec(StakingCodec)
	codec.RegisterCrypto(StakingCodec)
	StakingCodec.Seal()
}

func RegisterCodec(cdc *codec.Codec) {
	//
	cdc.RegisterConcrete(CreateValidatorTx{}, "ciChain/MsgCreateValidator", nil)
	cdc.RegisterConcrete(DelegateTx{}, "ciChain/MsgDelegate", nil)
	cdc.RegisterConcrete(UndelegateTx{}, "ciChain/MsgUndelegate", nil)
	cdc.RegisterConcrete(RedelegateTx{}, "ciChain/MsgBeginRedelegate", nil)
}