package types

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
)
var ModuleCdc *codec.Codec


func init()  {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}

func RegisterCodec(cdc *codec.Codec)  {

}