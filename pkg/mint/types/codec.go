package types

import "github.com/ci123chain/ci123chain/pkg/abci/codec"

var MintCdc *codec.Codec


func init() {
	MintCdc = codec.New()
	codec.RegisterCrypto(MintCdc)
	MintCdc.Seal()
}

func RegisterCodec(cdc *codec.Codec) {
	//
}