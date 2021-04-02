package types

import "github.com/pkg/errors"

var (
	ErrChannelNotFound = errors.New("channel not found")
	ErrInvalidChannelState = errors.New("channel state invalid")
	ErrInvalidChannelCapability = errors.New("channel capability invalid")
	ErrInvalidPacket = errors.New("invalid packet")
	ErrSequenceAckNotFound = errors.New("sequense ack not found")
	ErrInvalidSequence = errors.New("sequense invalid")
	ErrPacketTimeout = errors.New("packet timeout !")
	ErrSequenceSendNotFound = errors.New("sequence send not found !")
	ErrChannelCapabilityNotFound = errors.New("channel capability not found !")
	ErrSequenceReceiveNotFound = errors.New("sequence receive not found !")
	ErrAcknowledgementExists = errors.New("acknowledgement already exists !")
	ErrInvalidAcknowledgement = errors.New("invalid acknowledgement !")
	ErrInvalidChannelOrdering = errors.New("invalid channel ordering !")

)