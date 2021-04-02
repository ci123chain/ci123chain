package types

import "github.com/ci123chain/ci123chain/pkg/abci/codec"

var channelCdc *codec.Codec

func init()  {
	channelCdc = codec.New()
	RegisterCodec(channelCdc)
	codec.RegisterCrypto(channelCdc)
	channelCdc.Seal()
}

func RegisterCodec(cdc *codec.Codec)  {
	cdc.RegisterConcrete(&Acknowledgement{}, "ibcChannel/acknowledgement", nil)
}
