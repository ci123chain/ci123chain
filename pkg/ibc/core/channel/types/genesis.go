package types

// NewPacketState creates a new PacketState instance.
func NewPacketState(portID, channelID string, seq uint64, data []byte) PacketState {
	return PacketState{
		PortId:    portID,
		ChannelId: channelID,
		Sequence:  seq,
		Data:      data,
	}
}

// PacketSequence defines the genesis type necessary to retrieve and store
// next send and receive sequences.
type PacketSequence struct {
	PortId    string `protobuf:"bytes,1,opt,name=port_id,json=portId,proto3" json:"port_id,omitempty" yaml:"port_id"`
	ChannelId string `protobuf:"bytes,2,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty" yaml:"channel_id"`
	Sequence  uint64 `protobuf:"varint,3,opt,name=sequence,proto3" json:"sequence,omitempty"`
}

// GenesisState defines the ibc channel submodule's genesis state.
type GenesisState struct {
	Channels         []IdentifiedChannel `protobuf:"bytes,1,rep,name=channels,proto3,casttype=IdentifiedChannel" json:"channels"`
	Acknowledgements []PacketState       `protobuf:"bytes,2,rep,name=acknowledgements,proto3" json:"acknowledgements"`
	Commitments      []PacketState       `protobuf:"bytes,3,rep,name=commitments,proto3" json:"commitments"`
	Receipts         []PacketState       `protobuf:"bytes,4,rep,name=receipts,proto3" json:"receipts"`
	SendSequences    []PacketSequence    `protobuf:"bytes,5,rep,name=send_sequences,json=sendSequences,proto3" json:"send_sequences" yaml:"send_sequences"`
	RecvSequences    []PacketSequence    `protobuf:"bytes,6,rep,name=recv_sequences,json=recvSequences,proto3" json:"recv_sequences" yaml:"recv_sequences"`
	AckSequences     []PacketSequence    `protobuf:"bytes,7,rep,name=ack_sequences,json=ackSequences,proto3" json:"ack_sequences" yaml:"ack_sequences"`
	// the sequence for the next generated channel identifier
	NextChannelSequence uint64 `protobuf:"varint,8,opt,name=next_channel_sequence,json=nextChannelSequence,proto3" json:"next_channel_sequence,omitempty" yaml:"next_channel_sequence"`
}

// DefaultGenesisState returns the ibc channel submodule's default genesis state.
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Channels:            []IdentifiedChannel{},
		Acknowledgements:    []PacketState{},
		Receipts:            []PacketState{},
		Commitments:         []PacketState{},
		SendSequences:       []PacketSequence{},
		RecvSequences:       []PacketSequence{},
		AckSequences:        []PacketSequence{},
		NextChannelSequence: 0,
	}
}
