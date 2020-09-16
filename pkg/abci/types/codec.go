package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
)

// Register the sdk message type
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*Tx)(nil), nil)
	cdc.RegisterInterface((*Msg)(nil), nil)
}
