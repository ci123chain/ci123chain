package types

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/codec"
	"github.com/tanhuiya/ci123chain/pkg/supply/exported"
)

var ModuleCdc *codec.Codec

// RegisterCodec registers the account types and interface
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*exported.ModuleAccountI)(nil), nil)
	cdc.RegisterConcrete(&ModuleAccount{}, "ci123chain/ModuleAccount", nil)
}


func init()  {
	ModuleCdc := codec.New()
	RegisterCodec(ModuleCdc)
}