package types

import "CI123Chain/pkg/abci/codec"

// Register the sdk message type
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*Tx)(nil), nil)
}
