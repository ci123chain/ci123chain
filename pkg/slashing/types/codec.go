package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
)

var SlashingCodec *codec.Codec

func init(){
	SlashingCodec = codec.New()
	RegisterCodec(SlashingCodec)
	codec.RegisterCrypto(SlashingCodec)
	SlashingCodec.Seal()
}

// RegisterLegacyAminoCodec registers concrete types on LegacyAmino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&MsgUnjail{}, "ci123chain/MsgUnjail", nil)
}