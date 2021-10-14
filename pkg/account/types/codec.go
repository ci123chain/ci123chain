package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account/exported"
)

var ModuleCdc *codec.Codec

func init()  {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}

func RegisterCodec(cdc *codec.Codec)  {
	cdc.RegisterInterface((*exported.Account)(nil), nil)
	cdc.RegisterInterface((*exported.ModuleAccountI)(nil), nil)
	cdc.RegisterConcrete(&BaseAccount{}, "ci123chain/Account", nil)
	//cdc.RegisterConcrete(&types2.ModuleAccount{}, "ci123chain/ModuleAccount", nil)
	//cdc.RegisterConcrete(&util.HeightUpdate{}, "ci123chain/HeightUpdate", nil)
	//cdc.RegisterConcrete(&util.Heights{}, "ci123chain/Heights", nil)
	cdc.RegisterConcrete(&types.Coins{}, "ci123chain/Coins", nil)
}