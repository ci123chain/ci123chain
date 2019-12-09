package types

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
)
var DistributionCdc *codec.Codec


func init()  {
	DistributionCdc = codec.New()
	RegisterCodec(DistributionCdc)
	codec.RegisterCrypto(DistributionCdc)
	DistributionCdc.Seal()
}

func RegisterCodec(cdc *codec.Codec)  {

}