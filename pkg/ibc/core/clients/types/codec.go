package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
)


var IBCClientCodec *codec.Codec

func init(){
	IBCClientCodec = codec.New()
	RegisterCodec(IBCClientCodec)
	codec.RegisterCrypto(IBCClientCodec)
	IBCClientCodec.Seal()
}
func RegisterCodec(cdc *codec.Codec)  {
	cdc.RegisterInterface((*exported.ClientState)(nil), nil)
	cdc.RegisterInterface((*exported.ConsensusState)(nil), nil)

}

