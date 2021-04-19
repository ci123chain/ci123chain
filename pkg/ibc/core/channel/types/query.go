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



// QueryChannelResponse is the response type for the Query/Channel RPC method.
// Besides the Channel end, it includes a proof and the height from which the
// proof was retrieved.
type QueryChannelResponse struct {
	// channel associated with the request identifiers
	Channel *Channel `protobuf:"bytes,1,opt,name=channel,proto3" json:"channel,omitempty"`
	// merkle proof of existence
	Proof []byte `protobuf:"bytes,2,opt,name=proof,proto3" json:"proof,omitempty"`
	// height at which the proof was retrieved
	ProofHeight clienttypes.Height `protobuf:"bytes,3,opt,name=proof_height,json=proofHeight,proto3" json:"proof_height"`
}


// NewQueryChannelResponse creates a new QueryChannelResponse instance
func NewQueryChannelResponse(channel Channel, proof []byte, height clienttypes.Height) *QueryChannelResponse {
	return &QueryChannelResponse{
		Channel:     &channel,
		Proof:       proof,
		ProofHeight: height,
	}
}


// QueryChannelsRequest is the request type for the Query/Channels RPC method
type QueryChannelsRequest struct {
	// pagination request
	Pagination *pagination.PageRequest `protobuf:"bytes,1,opt,name=pagination,proto3" json:"pagination,omitempty"`
}
// QueryChannelsResponse is the response type for the Query/Channels RPC method.
type QueryChannelsResponse struct {
	// list of stored channels of the chain.
	Channels []*IdentifiedChannel `protobuf:"bytes,1,rep,name=channels,proto3" json:"channels,omitempty"`
	// pagination response
	Pagination *pagination.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
	// query block height
	Height clienttypes.Height `protobuf:"bytes,3,opt,name=height,proto3" json:"height"`
}


// QueryPacketAcknowledgementRequest is the request type for the
// Query/PacketAcknowledgement RPC method
type QueryPacketAcknowledgementRequest struct {
	// port unique identifier
	PortId string `protobuf:"bytes,1,opt,name=port_id,json=portId,proto3" json:"port_id,omitempty"`
	// channel unique identifier
	ChannelId string `protobuf:"bytes,2,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	// packet sequence
	Sequence uint64 `protobuf:"varint,3,opt,name=sequence,proto3" json:"sequence,omitempty"`
}

// QueryPacketAcknowledgementResponse defines the client query response for a
// packet which also includes a proof and the height from which the
// proof was retrieved
type QueryPacketAcknowledgementResponse struct {
	// packet associated with the request fields
	Acknowledgement []byte `protobuf:"bytes,1,opt,name=acknowledgement,proto3" json:"acknowledgement,omitempty"`
	// merkle proof of existence
	Proof []byte `protobuf:"bytes,2,opt,name=proof,proto3" json:"proof,omitempty"`
	// height at which the proof was retrieved
	ProofHeight clienttypes.Height `protobuf:"bytes,3,opt,name=proof_height,json=proofHeight,proto3" json:"proof_height"`
}

// QueryPacketAcknowledgementsRequest is the request type for the
// Query/QueryPacketCommitments RPC method
type QueryPacketAcknowledgementsRequest struct {
	// port unique identifier
	PortId string `protobuf:"bytes,1,opt,name=port_id,json=portId,proto3" json:"port_id,omitempty"`
	// channel unique identifier
	ChannelId string `protobuf:"bytes,2,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	// pagination request
	Pagination *pagination.PageRequest `protobuf:"bytes,3,opt,name=pagination,proto3" json:"pagination,omitempty"`
}
// QueryPacketAcknowledgemetsResponse is the request type for the
// Query/QueryPacketAcknowledgements RPC method
type QueryPacketAcknowledgementsResponse struct {
	Acknowledgements []*PacketState `protobuf:"bytes,1,rep,name=acknowledgements,proto3" json:"acknowledgements,omitempty"`
	// pagination response
	Pagination *pagination.PageResponse `protobuf:"bytes,2,opt,name=pagination,proto3" json:"pagination,omitempty"`
	// query block height
	Height clienttypes.Height `protobuf:"bytes,3,opt,name=height,proto3" json:"height"`
}
// NewQueryPacketAcknowledgementResponse creates a new QueryPacketAcknowledgementResponse instance
func NewQueryPacketAcknowledgementResponse(
	acknowledgement []byte, proof []byte, height clienttypes.Height,
) *QueryPacketAcknowledgementResponse {
	return &QueryPacketAcknowledgementResponse{
		Acknowledgement: acknowledgement,
		Proof:           proof,
		ProofHeight:     height,
	}
}


// QueryPacketReceiptRequest is the request type for the
// Query/PacketReceipt RPC method
type QueryPacketReceiptRequest struct {
	// port unique identifier
	PortId string `protobuf:"bytes,1,opt,name=port_id,json=portId,proto3" json:"port_id,omitempty"`
	// channel unique identifier
	ChannelId string `protobuf:"bytes,2,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	// packet sequence
	Sequence uint64 `protobuf:"varint,3,opt,name=sequence,proto3" json:"sequence,omitempty"`
}
// QueryPacketReceiptResponse defines the client query response for a packet receipt
// which also includes a proof, and the height from which the proof was
// retrieved
type QueryPacketReceiptResponse struct {
	// success flag for if receipt exists
	Received bool `protobuf:"varint,2,opt,name=received,proto3" json:"received,omitempty"`
	// merkle proof of existence
	Proof []byte `protobuf:"bytes,3,opt,name=proof,proto3" json:"proof,omitempty"`
	// height at which the proof was retrieved
	ProofHeight clienttypes.Height `protobuf:"bytes,4,opt,name=proof_height,json=proofHeight,proto3" json:"proof_height"`
}


// NewQueryPacketReceiptResponse creates a new QueryPacketReceiptResponse instance
func NewQueryPacketReceiptResponse(
	recvd bool, proof []byte, height clienttypes.Height,
) *QueryPacketReceiptResponse {
	return &QueryPacketReceiptResponse{
		Received:    recvd,
		Proof:       proof,
		ProofHeight: height,
	}
}
