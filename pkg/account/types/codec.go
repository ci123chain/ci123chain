package types

import "github.com/tanhuiya/ci123chain/pkg/abci/codec"

var ModuleCdc *codec.Codec

func init()  {
	ModuleCdc = codec.New()
	codec.RegisterCrypto(ModuleCdc)
}