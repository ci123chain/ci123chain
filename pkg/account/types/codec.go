package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/account/exported"
	"github.com/ci123chain/ci123chain/pkg/util"
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
	cdc.RegisterConcrete(&BaseAccount{}, "ci123chain/Account", nil)
	cdc.RegisterConcrete(&util.HeightUpdate{}, "ci123chain/HeightUpdate", nil)
	cdc.RegisterConcrete(&util.Heights{}, "ci123chain/Heights", nil)
	cdc.RegisterConcrete(&types.Coins{}, "ci123chain/Coins", nil)
}