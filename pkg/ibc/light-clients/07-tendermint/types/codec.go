package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
)

var LightClientCodec *codec.Codec

func init(){
	LightClientCodec = codec.New()
	RegisterCodec(LightClientCodec)
	codec.RegisterCrypto(LightClientCodec)
	LightClientCodec.Seal()
}
func RegisterCodec(cdc *codec.Codec)  {
	cdc.RegisterConcrete(&ClientState{}, "ibclightclient/ClientState", nil)
	cdc.RegisterConcrete(&ConsensusState{}, "ibclightclient/ConsensusState", nil)

}
