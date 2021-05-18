package types

import (
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
)

// IBC channel sentinel errors
var (
	ErrChannelExists             = sdkerrors.Register(SubModuleName, 2012, "channel already exists")
	ErrChannelNotFound           = sdkerrors.Register(SubModuleName, 2013, "channel not found")
	ErrInvalidChannel            = sdkerrors.Register(SubModuleName, 2014, "invalid channel")
	ErrInvalidChannelState       = sdkerrors.Register(SubModuleName, 2015, "invalid channel state")
	ErrInvalidChannelOrdering    = sdkerrors.Register(SubModuleName, 2016, "invalid channel ordering")
	ErrInvalidCounterparty       = sdkerrors.Register(SubModuleName, 2017, "invalid counterparty channel")
	ErrInvalidChannelCapability  = sdkerrors.Register(SubModuleName, 2018, "invalid channel capability")
	ErrChannelCapabilityNotFound = sdkerrors.Register(SubModuleName, 2019, "channel capability not found")
	ErrSequenceSendNotFound      = sdkerrors.Register(SubModuleName, 2010, "sequence send not found")
	ErrSequenceReceiveNotFound   = sdkerrors.Register(SubModuleName, 2021, "sequence receive not found")
	ErrSequenceAckNotFound       = sdkerrors.Register(SubModuleName, 2022, "sequence acknowledgement not found")
	ErrInvalidPacket             = sdkerrors.Register(SubModuleName, 2023, "invalid packet")
	ErrPacketTimeout             = sdkerrors.Register(SubModuleName, 2024, "packet timeout")
	ErrTooManyConnectionHops     = sdkerrors.Register(SubModuleName, 2025, "too many connection hops")
	ErrInvalidAcknowledgement    = sdkerrors.Register(SubModuleName, 2026, "invalid acknowledgement")
	ErrPacketCommitmentNotFound  = sdkerrors.Register(SubModuleName, 2027, "packet commitment not found")
	ErrPacketReceived            = sdkerrors.Register(SubModuleName, 2028, "packet already received")
	ErrAcknowledgementExists     = sdkerrors.Register(SubModuleName, 2029, "acknowledgement for packet already exists")
	ErrInvalidChannelIdentifier  = sdkerrors.Register(SubModuleName, 2030, "invalid channel identifier")
	ErrInvalidParam2   			 = sdkerrors.Register(SubModuleName, 2031, "param invalid")
)


func ErrInvalidParam(desc string) error {
	return sdkerrors.Wrapf(ErrInvalidParam2, desc)
}