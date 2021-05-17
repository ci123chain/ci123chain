package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

// MsgServer is the server API for Msg service.
type MsgServer interface {
	// Transfer defines a rpc handler method for MsgTransfer.
	Transfer(sdk.Context, *MsgTransfer) (*MsgTransferResponse, error)
}
//
//
//type MsgTransfer struct {
//	// the port on which the packet will be sent
//	SourcePort string `protobuf:"bytes,1,opt,name=source_port,json=sourcePort,proto3" json:"source_port,omitempty" yaml:"source_port"`
//	// the channel by which the packet will be sent
//	SourceChannel string `protobuf:"bytes,2,opt,name=source_channel,json=sourceChannel,proto3" json:"source_channel,omitempty" yaml:"source_channel"`
//	// the tokens to be transferred
//	Token sdk.Coin `protobuf:"bytes,3,opt,name=token,proto3" json:"token"`
//	// the sender address
//	Sender string `protobuf:"bytes,4,opt,name=sender,proto3" json:"sender,omitempty"`
//	// the recipient address on the destination chain
//	Receiver string `protobuf:"bytes,5,opt,name=receiver,proto3" json:"receiver,omitempty"`
//	// Timeout height relative to the current block height.
//	// The timeout is disabled when set to 0.
//	TimeoutHeight clienttypes.Height `protobuf:"bytes,6,opt,name=timeout_height,json=timeoutHeight,proto3" json:"timeout_height" yaml:"timeout_height"`
//	// Timeout timestamp (in nanoseconds) relative to the current block timestamp.
//	// The timeout is disabled when set to 0.
//	TimeoutTimestamp uint64 `protobuf:"varint,7,opt,name=timeout_timestamp,json=timeoutTimestamp,proto3" json:"timeout_timestamp,omitempty" yaml:"timeout_timestamp"`
//}
//
//// MsgTransferResponse defines the Msg/Transfer response type.
//type MsgTransferResponse struct {
//}
//
//// FungibleTokenPacketData defines a struct for the packet payload
//// See FungibleTokenPacketData spec:
//// https://github.com/cosmos/ics/tree/master/spec/ics-020-fungible-token-transfer#data-structures
//type FungibleTokenPacketData struct {
//	// the token denomination to be transferred
//	Denom string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
//	// the token amount to be transferred
//	Amount uint64 `protobuf:"varint,2,opt,name=amount,proto3" json:"amount,omitempty"`
//	// the sender address
//	Sender string `protobuf:"bytes,3,opt,name=sender,proto3" json:"sender,omitempty"`
//	// the recipient address on the destination chain
//	Receiver string `protobuf:"bytes,4,opt,name=receiver,proto3" json:"receiver,omitempty"`
//}
