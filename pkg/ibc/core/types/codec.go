package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	channeltypes "github.com/ci123chain/ci123chain/pkg/ibc/core/channel/types"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	connectiontypes "github.com/ci123chain/ci123chain/pkg/ibc/core/connection/types"
	lightclienttypes "github.com/ci123chain/ci123chain/pkg/ibc/light-clients/07-tendermint/types"
)

var ModuleCdc *codec.Codec

func init()  {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}

func RegisterCodec(cdc *codec.Codec)  {
	cdc.RegisterConcrete(&GenesisState{}, "ibccore/GenesisState", nil)

	clienttypes.RegisterCodec(cdc)
	connectiontypes.RegisterCodec(cdc)
	channeltypes.RegisterCodec(cdc)
	lightclienttypes.RegisterCodec(cdc)
}