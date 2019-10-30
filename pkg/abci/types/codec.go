package types

import "gitlab.oneitfarm.com/blockchain/ci123chain/pkg/abci/codec"

// Register the sdk message type
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*Tx)(nil), nil)
}
