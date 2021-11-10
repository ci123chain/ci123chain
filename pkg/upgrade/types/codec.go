package types

import "github.com/ci123chain/ci123chain/pkg/abci/codec"

var UpgradeCodec *codec.Codec

func init(){
	UpgradeCodec = codec.New()
	RegisterCodec(UpgradeCodec)
	codec.RegisterCrypto(UpgradeCodec)
	UpgradeCodec.Seal()
}

func RegisterCodec(cdc *codec.Codec) {
	//
	cdc.RegisterConcrete(&Plan{}, "ci123chain/Plan", nil)
}
