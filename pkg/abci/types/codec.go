package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	"github.com/ci123chain/ci123chain/pkg/abci/codec/types"
)

// Register the sdk message types
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*Tx)(nil), nil)
	cdc.RegisterInterface((*Msg)(nil), nil)
}

// RegisterInterfaces registers the sdk message type.
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterInterface("cosmos.base.v1beta1.Msg", (*Msg)(nil))
	registry.RegisterInterface("cosmos.base.v1beta1.Tx", (*Tx)(nil))
	// the interface name for MsgRequest is ServiceMsg because this is most useful for clients
	// to understand - it will be the way for clients to introspect on available Msg service methods
}