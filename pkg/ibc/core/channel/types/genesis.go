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
