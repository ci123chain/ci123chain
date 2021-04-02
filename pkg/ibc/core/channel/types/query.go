package types

import (
	"github.com/ci123chain/ci123chain/pkg/abci/types/pagination"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
)

// QueryPacketCommitmentRequest is the request type for the
// Query/PacketCommitment RPC method
type QueryPacketCommitmentRequest struct {
	// port unique identifier
	PortId string `protobuf:"bytes,1,opt,name=port_id,json=portId,proto3" json:"port_id,omitempty"`
	// channel unique identifier
	ChannelId string `protobuf:"bytes,2,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	// packet sequence
	Sequence uint64 `protobuf:"varint,3,opt,name=sequence,proto3" json:"sequence,omitempty"`
}


// QueryPacketCommitmentResponse defines the client query response for a packet
// which also includes a proof and the height from which the proof was
// retrieved
type QueryPacketCommitmentResponse struct {
	// packet associated with the request fields
	Commitment []byte `protobuf:"bytes,1,opt,name=commitment,proto3" json:"commitment,omitempty"`
	// merkle proof of existence
	Proof []byte `protobuf:"bytes,2,opt,name=proof,proto3" json:"proof,omitempty"`
	// height at which the proof was retrieved
	ProofHeight clienttypes.Height `protobuf:"bytes,3,opt,name=proof_height,json=proofHeight,proto3" json:"proof_height"`
}


// QueryPacketCommitmentsRequest is the request type for the
// Query/QueryPacketCommitments RPC method
type QueryPacketCommitmentsRequest struct {
	// port unique identifier
	PortId string `protobuf:"bytes,1,opt,name=port_id,json=portId,proto3" json:"port_id,omitempty"`
	// channel unique identifier
	ChannelId string `protobuf:"bytes,2,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	// pagination request
	Pagination *pagination.PageRequest `protobuf:"bytes,3,opt,name=pagination,proto3" json:"pagination,omitempty"`
}

// QueryPacketCommitmentsResponse is the request type for the
// Query/QueryPacketCommitments RPC method
type QueryPacketCommitmentsResponse struct {
	Commitments []*PacketState `protobuf:"bytes,1,rep,name=commitments,proto3" json:"commitments,omitempty"`
	// pagination response
	Pagination *pagination.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
	// query block height
	Height clienttypes.Height `protobuf:"bytes,3,opt,name=height,proto3" json:"height"`
}


// PacketState defines the generic type necessary to retrieve and store
// packet commitments, acknowledgements, and receipts.
// Caller is responsible for knowing the context necessary to interpret this
// state as a commitment, acknowledgement, or a receipt.
type PacketState struct {
	// channel port identifier.
	PortId string `protobuf:"bytes,1,opt,name=port_id,json=portId,proto3" json:"port_id,omitempty" yaml:"port_id"`
	// channel unique identifier.
	ChannelId string `protobuf:"bytes,2,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty" yaml:"channel_id"`
	// packet sequence.
	Sequence uint64 `protobuf:"varint,3,opt,name=sequence,proto3" json:"sequence,omitempty"`
	// embedded data that represents packet state.
	Data []byte `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
}


// QueryUnreceivedPacketsRequest is the request type for the
// Query/UnreceivedPackets RPC method
type QueryUnreceivedPacketsRequest struct {
	// port unique identifier
	PortId string `protobuf:"bytes,1,opt,name=port_id,json=portId,proto3" json:"port_id,omitempty"`
	// channel unique identifier
	ChannelId string `protobuf:"bytes,2,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	// list of packet sequences
	PacketCommitmentSequences []uint64 `protobuf:"varint,3,rep,packed,name=packet_commitment_sequences,json=packetCommitmentSequences,proto3" json:"packet_commitment_sequences,omitempty"`
}

// QueryUnreceivedPacketsResponse is the response type for the
// Query/UnreceivedPacketCommitments RPC method
type QueryUnreceivedPacketsResponse struct {
	// list of unreceived packet sequences
	Sequences []uint64 `protobuf:"varint,1,rep,packed,name=sequences,proto3" json:"sequences,omitempty"`
	// query block height
	Height clienttypes.Height `protobuf:"bytes,2,opt,name=height,proto3" json:"height"`
}