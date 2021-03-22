package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
)

var ModuleCdc *codec.Codec

func init()  {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}

func RegisterCodec(cdc *codec.Codec)  {
	clienttypes.RegisterCodec(cdc)
}