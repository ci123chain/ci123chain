package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/supply/exported"
)

var ModuleCdc *codec.Codec

// RegisterCodec registers the account types and interface
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*exported.ModuleAccountI)(nil), nil)
	cdc.RegisterInterface((*exported.SupplyI)(nil), nil)
	cdc.RegisterConcrete(&ModuleAccount{}, "ci123chain/ModuleAccount", nil)
	cdc.RegisterConcrete(&Supply{}, "ci123chain/Supply", nil)
}


func init()  {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
}