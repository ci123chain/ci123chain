package types

import (
	"crypto/sha256"
	"github.com/ci123chain/ci123chain/pkg/abci/codec"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	ibcerrors "github.com/ci123chain/ci123chain/pkg/ibc/core/errors"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	"github.com/pkg/errors"
)


// CommitPacket returns the packet commitment bytes. The commitment consists of:
// sha256_hash(timeout_timestamp + timeout_height.RevisionNumber + timeout_height.RevisionHeight + sha256_hash(data))
// from a given packet. This results in a fixed length preimage.
// NOTE: sdk.Uint64ToBigEndian sets the uint64 to a slice of length 8.
func CommitPacket(cdc *codec.Codec, packet exported.PacketI) []byte {
	timeoutHeight := packet.GetTimeoutHeight()

	buf := sdk.Uint64ToBigEndian(packet.GetTimeoutTimestamp())

	revisionNumber := sdk.Uint64ToBigEndian(timeoutHeight.GetRevisionNumber())
	buf = append(buf, revisionNumber...)

	revisionHeight := sdk.Uint64ToBigEndian(timeoutHeight.GetRevisionHeight())
	buf = append(buf, revisionHeight...)

	dataHash := sha256.Sum256(packet.GetData())
	buf = append(buf, dataHash[:]...)

	hash := sha256.Sum256(buf)
	return hash[:]
}


var _ exported.PacketI = (*Packet)(nil)

type Packet struct {
	// number corresponds to the order of sends and receives, where a Packet
	// with an earlier sequence number must be sent and received before a Packet
	// with a later sequence number.
	Sequence uint64 `protobuf:"varint,1,opt,name=sequence,proto3" json:"sequence,omitempty"`
	// identifies the port on the sending chain.
	SourcePort string `protobuf:"bytes,2,opt,name=source_port,json=sourcePort,proto3" json:"source_port,omitempty" yaml:"source_port"`
	// identifies the channel end on the sending chain.
	SourceChannel string `protobuf:"bytes,3,opt,name=source_channel,json=sourceChannel,proto3" json:"source_channel,omitempty" yaml:"source_channel"`
	// identifies the port on the receiving chain.
	DestinationPort string `protobuf:"bytes,4,opt,name=destination_port,json=destinationPort,proto3" json:"destination_port,omitempty" yaml:"destination_port"`
	// identifies the channel end on the receiving chain.
	DestinationChannel string `protobuf:"bytes,5,opt,name=destination_channel,json=destinationChannel,proto3" json:"destination_channel,omitempty" yaml:"destination_channel"`
	// actual opaque bytes transferred directly to the application module
	Data []byte `protobuf:"bytes,6,opt,name=data,proto3" json:"data,omitempty"`
	// block height after which the packet times out
	TimeoutHeight clienttypes.Height `protobuf:"bytes,7,opt,name=timeout_height,json=timeoutHeight,proto3" json:"timeout_height" yaml:"timeout_height"`
	// block timestamp (in nanoseconds) after which the packet times out
	TimeoutTimestamp uint64 `protobuf:"varint,8,opt,name=timeout_timestamp,json=timeoutTimestamp,proto3" json:"timeout_timestamp,omitempty" yaml:"timeout_timestamp"`
}
// NewPacket creates a new Packet instance. It panics if the provided
// packet data interface is not registered.
func NewPacket(
	data []byte,
	sequence uint64, sourcePort, sourceChannel,
	destinationPort, destinationChannel string,
	timeoutHeight clienttypes.Height, timeoutTimestamp uint64,
) Packet {
	return Packet{
		Data:               data,
		Sequence:           sequence,
		SourcePort:         sourcePort,
		SourceChannel:      sourceChannel,
		DestinationPort:    destinationPort,
		DestinationChannel: destinationChannel,
		TimeoutHeight:      timeoutHeight,
		TimeoutTimestamp:   timeoutTimestamp,
	}
}


// GetSequence implements PacketI interface
func (p Packet) GetSequence() uint64 { return p.Sequence }

// GetSourcePort implements PacketI interface
func (p Packet) GetSourcePort() string { return p.SourcePort }

// GetSourceChannel implements PacketI interface
func (p Packet) GetSourceChannel() string { return p.SourceChannel }

// GetDestPort implements PacketI interface
func (p Packet) GetDestPort() string { return p.DestinationPort }

// GetDestChannel implements PacketI interface
func (p Packet) GetDestChannel() string { return p.DestinationChannel }

// GetData implements PacketI interface
func (p Packet) GetData() []byte { return p.Data }

// GetTimeoutHeight implements PacketI interface
func (p Packet) GetTimeoutHeight() exported.Height { return p.TimeoutHeight }

// GetTimeoutTimestamp implements PacketI interface
func (p Packet) GetTimeoutTimestamp() uint64 { return p.TimeoutTimestamp }

// ValidateBasic implements PacketI interface
func (p Packet) ValidateBasic() error {
	if err := host.PortIdentifierValidator(p.SourcePort); err != nil {
		return errors.Wrap(err, "invalid source port ID")
	}
	if err := host.PortIdentifierValidator(p.DestinationPort); err != nil {
		return errors.Wrap(err, "invalid destination port ID")
	}
	if err := host.ChannelIdentifierValidator(p.SourceChannel); err != nil {
		return errors.Wrap(err, "invalid source channel ID")
	}
	if err := host.ChannelIdentifierValidator(p.DestinationChannel); err != nil {
		return errors.Wrap(err, "invalid destination channel ID")
	}
	if p.Sequence == 0 {
		return errors.Wrap(ibcerrors.ErrInvalidPacket, "packet sequence cannot be 0")
	}
	if p.TimeoutHeight.IsZero() && p.TimeoutTimestamp == 0 {
		return errors.Wrap(ibcerrors.ErrInvalidPacket, "packet timeout height and packet timeout timestamp cannot both be 0")
	}
	if len(p.Data) == 0 {
		return errors.Wrap(ibcerrors.ErrInvalidPacket, "packet data bytes cannot be empty")
	}
	return nil
}

// CommitAcknowledgement returns the hash of commitment bytes
func CommitAcknowledgement(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}