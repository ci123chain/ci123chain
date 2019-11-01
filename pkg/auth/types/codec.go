package types

import "github.com/tanhuiya/ci123chain/pkg/abci/codec"

func RegisterCodec(cdc *codec.Codec)  {
	
}

var ModuleCdc *codec.Codec

func init()  {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
}