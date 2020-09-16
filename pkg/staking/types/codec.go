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
	cdc.RegisterConcrete(&MsgCreateValidator{}, "ci123chain/MsgCreateValidator", nil)
	cdc.RegisterConcrete(&MsgEditValidator{}, "ci123chain/MsgEditValidator", nil)
	cdc.RegisterConcrete(&MsgDelegate{}, "ci123chain/MsgDelegate", nil)
	cdc.RegisterConcrete(&MsgUndelegate{}, "ci123chain/MsgUndelegate", nil)
	cdc.RegisterConcrete(&MsgRedelegate{}, "ci123chain/MsgRedelegate", nil)
}