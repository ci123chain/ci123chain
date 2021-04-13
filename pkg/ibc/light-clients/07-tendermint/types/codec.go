package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	types2 "github.com/tendermint/tendermint/proto/tendermint/types"
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
	cdc.RegisterConcrete(&types2.ValidatorSet{}, "ibclightclient/tendermint.ValidatorSet", nil)
	cdc.RegisterConcrete(&Header{}, "ibclightclient/Header", nil)
}
