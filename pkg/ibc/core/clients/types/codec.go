package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
)

func RegisterCodec(cdc *codec.Codec)  {
	cdc.RegisterInterface((*exported.ClientState)(nil), nil)
	cdc.RegisterInterface((*exported.ConsensusState)(nil), nil)
}