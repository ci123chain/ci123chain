package types

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	cosmosSdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ cosmosSdk.Msg = &MsgChannelOpenInit{}
	_ cosmosSdk.Msg = &MsgChannelOpenTry{}
	_ cosmosSdk.Msg = &MsgChannelOpenAck{}
	_ cosmosSdk.Msg = &MsgChannelOpenConfirm{}
	_ cosmosSdk.Msg = &MsgAcknowledgement{}
	_ cosmosSdk.Msg = &MsgRecvPacket{}
	_ cosmosSdk.Msg = &MsgTimeout{}
)

var _ sdk.Msg = &MsgChannelOpenInit{}
// NewMsgChannelOpenInit creates a new MsgChannelOpenInit. It sets the counterparty channel
// identifier to be empty.
// nolint:interfacer
func NewMsgChannelOpenInit(
	portID, version string, channelOrder Order, connectionHops []string,
	counterpartyPortID string, signer string,
) *MsgChannelOpenInit {
	counterparty := NewCounterparty(counterpartyPortID, "")
	channel := NewChannel(INIT, channelOrder, counterparty, connectionHops, version)
	return &MsgChannelOpenInit{
		PortId:  portID,
		Channel: channel,
		Signer:  signer,
	}
}
func (m MsgChannelOpenInit) Route() string {
	return host.RouterKey
}

func (m MsgChannelOpenInit) MsgType() string {
	return "channel_open_init"
}

func (msg MsgChannelOpenInit) ValidateBasic() error {
	if err := host.PortIdentifierValidator(msg.PortId); err != nil {
		return ErrInvalidParam("invalid port ID")
	}
	if msg.Channel.State != INIT {
		return ErrInvalidParam(fmt.Sprintf(
			"channel state must be INIT in MsgChannelOpenInit. expected: %s, got: %s",
			INIT, msg.Channel.State,
		))
	}
	if msg.Channel.Counterparty.ChannelId != "" {
		return ErrInvalidParam( "counterparty channel identifier must be empty")
	}
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return ErrInvalidParam(fmt.Sprintf("invalid signer: %v", msg.Signer))
	}
	return msg.Channel.ValidateBasic()
}

func (m MsgChannelOpenInit) GetFromAddress() sdk.AccAddress {
	return sdk.HexToAddress(m.Signer)
}

func (m MsgChannelOpenInit) Bytes() []byte {
	return ChannelCdc.MustMarshalBinaryLengthPrefixed(m)
}

func (m MsgChannelOpenInit) GetSignBytes() []byte {
	panic("IBC messages do not support amino")
}

func (m MsgChannelOpenInit) GetSigners() []cosmosSdk.AccAddress {
	return []cosmosSdk.AccAddress{sdk.HexToAddress(m.Signer).Bytes()}
}

func (m MsgChannelOpenInit) Type() string {
	return "channel_open_init"
}

// -----------MsgChannelOpenTry
var _ sdk.Msg = &MsgChannelOpenTry{}

// NewMsgChannelOpenTry creates a new MsgChannelOpenTry instance
// nolint:interfacer
func NewMsgChannelOpenTry(
	portID, previousChannelID, version string, channelOrder Order, connectionHops []string,
	counterpartyPortID, counterpartyChannelID, counterpartyVersion string,
	proofInit []byte, proofHeight clienttypes.Height, signer string,
) *MsgChannelOpenTry {
	counterparty := NewCounterparty(counterpartyPortID, counterpartyChannelID)
	channel := NewChannel(TRYOPEN, channelOrder, counterparty, connectionHops, version)
	return &MsgChannelOpenTry{
		PortId:              portID,
		PreviousChannelId:   previousChannelID,
		Channel:             channel,
		CounterpartyVersion: counterpartyVersion,
		ProofInit:           proofInit,
		ProofHeight:         proofHeight,
		Signer:              signer,
	}
}


func (m MsgChannelOpenTry) Route() string {
	return host.RouterKey
}

func (m MsgChannelOpenTry) MsgType() string {
	return "channel_open_try"
}

func (msg MsgChannelOpenTry) ValidateBasic() error {
	if err := host.PortIdentifierValidator(msg.PortId); err != nil {
		return ErrInvalidParam("invalid port ID")
	}
	if msg.PreviousChannelId != "" {
		if !IsValidChannelID(msg.PreviousChannelId) {
			return ErrInvalidParam( "invalid previous channel ID")
		}
	}
	if len(msg.ProofInit) == 0 {
		return ErrInvalidParam( "cannot submit an empty proof init")
	}
	if msg.ProofHeight.IsZero() {
		return ErrInvalidParam("proof height must be non-zero")
	}
	if msg.Channel.State != TRYOPEN {
		return ErrInvalidParam(fmt.Sprintf(
			"channel state must be TRYOPEN in MsgChannelOpenTry. expected: %s, got: %s",
			TRYOPEN, msg.Channel.State,
		))
	}
	// counterparty validate basic allows empty counterparty channel identifiers
	if err := host.ChannelIdentifierValidator(msg.Channel.Counterparty.ChannelId); err != nil {
		return ErrInvalidParam("invalid counterparty channel ID")
	}

	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return ErrInvalidParam(fmt.Sprintf("invalid signer: %v", msg.Signer))
	}
	return msg.Channel.ValidateBasic()
}

func (m MsgChannelOpenTry) GetFromAddress() sdk.AccAddress {
	return sdk.HexToAddress(m.Signer)
}

func (m MsgChannelOpenTry) Bytes() []byte {
	return ChannelCdc.MustMarshalBinaryLengthPrefixed(m)
}

func (m MsgChannelOpenTry) GetSignBytes() []byte {
	panic("IBC messages do not support amino")
}

func (m MsgChannelOpenTry) GetSigners() []cosmosSdk.AccAddress {
	return []cosmosSdk.AccAddress{sdk.HexToAddress(m.Signer).Bytes()}
}

func (m MsgChannelOpenTry) Type() string {
	return "channel_open_try"
}

// -----------MsgChannelOpenAck
var _ sdk.Msg = &MsgChannelOpenAck{}
// NewMsgChannelOpenAck creates a new MsgChannelOpenAck instance
// nolint:interfacer
func NewMsgChannelOpenAck(
	portID, channelID, counterpartyChannelID string, cpv string, proofTry []byte, proofHeight clienttypes.Height,
	signer string,
) *MsgChannelOpenAck {
	return &MsgChannelOpenAck{
		PortId:                portID,
		ChannelId:             channelID,
		CounterpartyChannelId: counterpartyChannelID,
		CounterpartyVersion:   cpv,
		ProofTry:              proofTry,
		ProofHeight:           proofHeight,
		Signer:                signer,
	}
}

func (m MsgChannelOpenAck) Route() string {
	return host.RouterKey
}

func (m MsgChannelOpenAck) MsgType() string {
	return "channel_open_ack"
}

func (msg MsgChannelOpenAck) ValidateBasic() error {
	if err := host.PortIdentifierValidator(msg.PortId); err != nil {
		return ErrInvalidParam("invalid port ID")
	}
	if !IsValidChannelID(msg.ChannelId) {
		return ErrInvalidChannelIdentifier
	}
	if err := host.ChannelIdentifierValidator(msg.CounterpartyChannelId); err != nil {
		return ErrInvalidParam( "invalid counterparty channel ID")
	}
	if len(msg.ProofTry) == 0 {
		return ErrInvalidParam( "cannot submit an empty proof try")
	}
	if msg.ProofHeight.IsZero() {
		return ErrInvalidParam("proof height must be non-zero")
	}
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return ErrInvalidParam(fmt.Sprintf("invalid signer: %v", msg.Signer))
	}
	return nil
}

func (m MsgChannelOpenAck) GetFromAddress() sdk.AccAddress {
	return sdk.HexToAddress(m.Signer)
}

func (m MsgChannelOpenAck) Bytes() []byte {
	return ChannelCdc.MustMarshalBinaryLengthPrefixed(m)

}

func (m MsgChannelOpenAck) GetSignBytes() []byte {
	panic("IBC messages do not support amino")
}

func (m MsgChannelOpenAck) GetSigners() []cosmosSdk.AccAddress {
	return []cosmosSdk.AccAddress{sdk.HexToAddress(m.Signer).Bytes()}
}

func (m MsgChannelOpenAck) Type() string {
	return "channel_open_ack"
}

// -----------MsgChannelOpenConfirm
var _ sdk.Msg = &MsgChannelOpenConfirm{}

// NewMsgChannelOpenConfirm creates a new MsgChannelOpenConfirm instance
// nolint:interfacer
func NewMsgChannelOpenConfirm(
	portID, channelID string, proofAck []byte, proofHeight clienttypes.Height,
	signer string,
) *MsgChannelOpenConfirm {
	return &MsgChannelOpenConfirm{
		PortId:      portID,
		ChannelId:   channelID,
		ProofAck:    proofAck,
		ProofHeight: proofHeight,
		Signer:      signer,
	}
}

func (m MsgChannelOpenConfirm) Route() string {
	return host.RouterKey
}

func (m MsgChannelOpenConfirm) MsgType() string {
	return "channel_open_confirm"
}

func (msg MsgChannelOpenConfirm) ValidateBasic() error {
	if err := host.PortIdentifierValidator(msg.PortId); err != nil {
		return ErrInvalidParam( "invalid port ID")
	}
	if !IsValidChannelID(msg.ChannelId) {
		return ErrInvalidChannelIdentifier
	}
	if len(msg.ProofAck) == 0 {
		return ErrInvalidParam( "cannot submit an empty proof ack")
	}
	if msg.ProofHeight.IsZero() {
		return ErrInvalidParam( "proof height must be non-zero")
	}
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return ErrInvalidParam(fmt.Sprintf("invalid signer: %v", msg.Signer))
	}
	return nil
}

func (m MsgChannelOpenConfirm) GetFromAddress() sdk.AccAddress {
	return sdk.HexToAddress(m.Signer)
}

func (m MsgChannelOpenConfirm) Bytes() []byte {
	return ChannelCdc.MustMarshalBinaryLengthPrefixed(m)
}

func (m MsgChannelOpenConfirm) GetSignBytes() []byte {
	panic("IBC messages do not support amino")
}

func (m MsgChannelOpenConfirm) GetSigners() []cosmosSdk.AccAddress {
	return []cosmosSdk.AccAddress{sdk.HexToAddress(m.Signer).Bytes()}
}

func (m MsgChannelOpenConfirm) Type() string {
	return "channel_open_confirm"
}

var _ sdk.Msg = &MsgAcknowledgement{}

// NewMsgAcknowledgement constructs a new MsgAcknowledgement
// nolint:interfacer
func NewMsgAcknowledgement(
	packet Packet,
	ack, proofAcked []byte,
	proofHeight clienttypes.Height,
	signer string,
) *MsgAcknowledgement {
	return &MsgAcknowledgement{
		Packet:          packet,
		Acknowledgement: ack,
		ProofAcked:      proofAcked,
		ProofHeight:     proofHeight,
		Signer:          signer,
	}
}


func (m MsgAcknowledgement) Route() string {
	return host.RouterKey
}

func (m MsgAcknowledgement) MsgType() string {
	return "acknowledge_packet"
}

func (msg MsgAcknowledgement) ValidateBasic() error {
	if len(msg.ProofAcked) == 0 {
		return ErrInvalidParam("cannot submit an empty proof")
	}
	if msg.ProofHeight.IsZero() {
		return ErrInvalidParam( "proof height must be non-zero")
	}
	if len(msg.Acknowledgement) == 0 {
		return ErrInvalidParam( "ack bytes cannot be empty")
	}
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return ErrInvalidParam(fmt.Sprintf("invalid signer: %v", msg.Signer))
	}
	return msg.Packet.ValidateBasic()
}

func (m MsgAcknowledgement) GetFromAddress() sdk.AccAddress {
	return sdk.HexToAddress(m.Signer)
}

func (m MsgAcknowledgement) Bytes() []byte {
	return ChannelCdc.MustMarshalBinaryLengthPrefixed(m)
}

func (m MsgAcknowledgement) GetSignBytes() []byte {
	panic("IBC messages do not support amino")
}

func (m MsgAcknowledgement) GetSigners() []cosmosSdk.AccAddress {
	return []cosmosSdk.AccAddress{sdk.HexToAddress(m.Signer).Bytes()}
}

func (m MsgAcknowledgement) Type() string {
	return "acknowledge_packet"
}

var _ sdk.Msg = &MsgRecvPacket{}

// NewMsgRecvPacket constructs new MsgRecvPacket
// nolint:interfacer
func NewMsgRecvPacket(
	packet Packet, proofCommitment []byte, proofHeight clienttypes.Height,
	signer string,
) *MsgRecvPacket {
	return &MsgRecvPacket{
		Packet:          packet,
		ProofCommitment: proofCommitment,
		ProofHeight:     proofHeight,
		Signer:          signer,
	}
}


func (m MsgRecvPacket) Route() string {
	return host.RouterKey
}

func (m MsgRecvPacket) MsgType() string {
	return "receive_packet"
}

func (msg MsgRecvPacket) ValidateBasic() error {
	if len(msg.ProofCommitment) == 0 {
		return ErrInvalidParam( "cannot submit an empty proof")
	}
	if msg.ProofHeight.IsZero() {
		return ErrInvalidParam( "proof height must be non-zero")
	}
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return ErrInvalidParam(fmt.Sprintf("invalid signer: %v", msg.Signer))
	}
	return msg.Packet.ValidateBasic()
}

func (m MsgRecvPacket) GetFromAddress() sdk.AccAddress {
	return sdk.HexToAddress(m.Signer)
}

func (m MsgRecvPacket) Bytes() []byte {
	return ChannelCdc.MustMarshalBinaryLengthPrefixed(m)
}

func (m MsgRecvPacket) GetSignBytes() []byte {
	panic("IBC messages do not support amino")
}

func (m MsgRecvPacket) GetSigners() []cosmosSdk.AccAddress {
	return []cosmosSdk.AccAddress{sdk.HexToAddress(m.Signer).Bytes()}
}

func (m MsgRecvPacket) Type() string {
	return "receive_packet"
}

var _ sdk.Msg = &MsgTimeout{}

// NewMsgTimeout constructs new MsgTimeout
// nolint:interfacer
func NewMsgTimeout(
	packet Packet, nextSequenceRecv uint64, proofUnreceived []byte,
	proofHeight clienttypes.Height, signer string,
) *MsgTimeout {
	return &MsgTimeout{
		Packet:           packet,
		NextSequenceRecv: nextSequenceRecv,
		ProofUnreceived:  proofUnreceived,
		ProofHeight:      proofHeight,
		Signer:           signer,
	}
}


func (m MsgTimeout) Route() string {
	return host.RouterKey
}

func (m MsgTimeout) MsgType() string {
	return "timeout_packet"
}

func (msg MsgTimeout) ValidateBasic() error {
	if len(msg.ProofUnreceived) == 0 {
		return ErrInvalidParam( "cannot submit an empty unreceived proof")
	}
	if msg.ProofHeight.IsZero() {
		return ErrInvalidParam("proof height must be non-zero")
	}
	if msg.NextSequenceRecv == 0 {
		return ErrInvalidParam( "next sequence receive cannot be 0")
	}
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return ErrInvalidParam(fmt.Sprintf("invalid signer: %v", msg.Signer))
	}
	return msg.Packet.ValidateBasic()
}

func (m MsgTimeout) GetFromAddress() sdk.AccAddress {
	return sdk.HexToAddress(m.Signer)
}

func (m MsgTimeout) Bytes() []byte {
	return ChannelCdc.MustMarshalBinaryLengthPrefixed(m)
}

func (m MsgTimeout) GetSignBytes() []byte {
	panic("IBC messages do not support amino")
}

func (m MsgTimeout) GetSigners() []cosmosSdk.AccAddress {
	return []cosmosSdk.AccAddress{sdk.HexToAddress(m.Signer).Bytes()}
}

func (m MsgTimeout) Type() string {
	return "timeout_packet"
}

// NewQueryPacketCommitmentResponse creates a new QueryPacketCommitmentResponse instance
func NewQueryPacketCommitmentResponse(
	commitment []byte, proof []byte, height clienttypes.Height,
) *QueryPacketCommitmentResponse {
	return &QueryPacketCommitmentResponse{
		Commitment:  commitment,
		Proof:       proof,
		ProofHeight: height,
	}
}