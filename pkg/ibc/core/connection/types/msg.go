package types

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	commitmenttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/commitment/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/exported"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
)

var _ sdk.Msg = &MsgConnectionOpenInit{}
var _ sdk.Msg = &MsgConnectionOpenTry{}
var _ sdk.Msg = &MsgConnectionOpenAck{}
var _ sdk.Msg = &MsgConnectionOpenConfirm{}

func (msg MsgConnectionOpenInit) MsgType() string {
	return "connection_open_init"
}

func (msg MsgConnectionOpenInit) GetFromAddress() sdk.AccAddress {
	accAddr := sdk.HexToAddress(msg.Signer)
	return accAddr
}

func (msg MsgConnectionOpenInit) Bytes() []byte {
	panic("IBC Message doesnot implement amino")
}

// NewMsgConnectionOpenInit creates a new MsgConnectionOpenInit instance. It sets the
// counterparty connection identifier to be empty.
//nolint:interfacer
func NewMsgConnectionOpenInit(
	clientID, counterpartyClientID string,
	counterpartyPrefix commitmenttypes.MerklePrefix,
	version *Version, delayPeriod uint64, signer sdk.AccAddress,
) *MsgConnectionOpenInit {
	// counterparty must have the same delay period
	counterparty := NewCounterparty(counterpartyClientID, "", counterpartyPrefix)
	return &MsgConnectionOpenInit{
		ClientId:     clientID,
		Counterparty: counterparty,
		Version:      version,
		DelayPeriod:  delayPeriod,
		Signer:       signer.String(),
	}
}


// Route implements sdk.Msg
func (msg MsgConnectionOpenInit) Route() string {
	return host.RouterKey
}

// Type implements sdk.Msg
func (msg MsgConnectionOpenInit) Type() string {
	return "connection_open_init"
}

// ValidateBasic implements sdk.Msg.
func (msg MsgConnectionOpenInit) ValidateBasic() error {
	if err := host.ClientIdentifierValidator(msg.ClientId); err != nil {
		return ErrInvalidParam( "invalid client ID")
	}
	if msg.Counterparty.ConnectionId != "" {
		return ErrInvalidParam(  "counterparty connection identifier must be empty")
	}

	// NOTE: Version can be nil on MsgConnectionOpenInit
	if msg.Version != nil {
		if err := ValidateVersion(msg.Version); err != nil {
			return ErrInvalidParam( "basic validation of the provided version failed")
		}
	}
	_ = sdk.HexToAddress(msg.Signer)
	if err := msg.Counterparty.ValidateBasic(); err != nil {
		return ErrInvalidParam(  err.Error())
	}
	return nil
}

// GetSignBytes implements sdk.Msg. The function will panic since it is used
// for amino transaction verification which IBC does not support.
func (msg MsgConnectionOpenInit) GetSignBytes() []byte {
	panic("IBC messages do not support amino")
}

// GetSigners implements sdk.Msg
func (msg MsgConnectionOpenInit) GetSigners() []sdk.AccAddress {
	accAddr := sdk.HexToAddress(msg.Signer)
	return []sdk.AccAddress{accAddr}
}





// NewMsgConnectionOpenTry creates a new MsgConnectionOpenTry instance
//nolint:interfacer
func NewMsgConnectionOpenTry(
	previousConnectionID, clientID, counterpartyConnectionID,
	counterpartyClientID string, counterpartyClient exported.ClientState,
	counterpartyPrefix commitmenttypes.MerklePrefix,
	counterpartyVersions []*Version, delayPeriod uint64,
	proofInit, proofClient, proofConsensus []byte,
	proofHeight, consensusHeight clienttypes.Height, signer sdk.AccAddress,
) *MsgConnectionOpenTry {
	counterparty := NewCounterparty(counterpartyClientID, counterpartyConnectionID, counterpartyPrefix)
	return &MsgConnectionOpenTry{
		PreviousConnectionId: previousConnectionID,
		ClientId:             clientID,
		ClientState:          counterpartyClient,
		Counterparty:         counterparty,
		CounterpartyVersions: counterpartyVersions,
		DelayPeriod:          delayPeriod,
		ProofInit:            proofInit,
		ProofClient:          proofClient,
		ProofConsensus:       proofConsensus,
		ProofHeight:          proofHeight,
		ConsensusHeight:      consensusHeight,
		Signer:               signer.String(),
	}
}

// Route implements sdk.Msg
func (msg MsgConnectionOpenTry) Route() string {
	return host.RouterKey
}

// Type implements sdk.Msg
func (msg MsgConnectionOpenTry) Type() string {
	return "connection_open_try"
}

func (msg MsgConnectionOpenTry) MsgType() string {
	return "connection_open_try"
}
// ValidateBasic implements sdk.Msg
func (msg MsgConnectionOpenTry) ValidateBasic() error {
	// an empty connection identifier indicates that a connection identifier should be generated
	if msg.PreviousConnectionId != "" {
		if !IsValidConnectionID(msg.PreviousConnectionId) {
			return ErrInvalidParam( "invalid previous connection ID")
		}
	}
	if err := host.ClientIdentifierValidator(msg.ClientId); err != nil {
		return ErrInvalidParam( "invalid client ID")
	}
	// counterparty validate basic allows empty counterparty connection identifiers
	if err := host.ConnectionIdentifierValidator(msg.Counterparty.ConnectionId); err != nil {
		return ErrInvalidParam( "invalid counterparty connection ID")
	}
	if msg.ClientState == nil {
		return ErrInvalidParam(  "counterparty client is nil")
	}
	if err := msg.ClientState.Validate(); err != nil {
		return ErrInvalidParam( "counterparty client is invalid")
	}
	if len(msg.CounterpartyVersions) == 0 {
		return ErrInvalidParam( "empty counterparty versions")
	}
	for i, version := range msg.CounterpartyVersions {
		if err := ValidateVersion(version); err != nil {
			return ErrInvalidParam( fmt.Sprintf("basic validation failed on version with index %d", i))
		}
	}
	if len(msg.ProofInit) == 0 {
		return ErrInvalidParam( "cannot submit an empty proof init")
	}
	if len(msg.ProofClient) == 0 {
		return ErrInvalidParam( "cannot submit empty proof client")
	}
	if len(msg.ProofConsensus) == 0 {
		return ErrInvalidParam( "cannot submit an empty proof of consensus state")
	}
	if msg.ProofHeight.IsZero() {
		return ErrInvalidParam( "proof height must be non-zero")
	}
	if msg.ConsensusHeight.IsZero() {
		return ErrInvalidParam( "consensus height must be non-zero")
	}

	return msg.Counterparty.ValidateBasic()
}

// GetSigners implements sdk.Msg
func (msg MsgConnectionOpenTry) GetSigners() []sdk.AccAddress {
	accAddr := sdk.HexToAddress(msg.Signer)
	return []sdk.AccAddress{accAddr}
}



func (msg MsgConnectionOpenTry) GetFromAddress() sdk.AccAddress {
	return sdk.HexToAddress(msg.Signer)
}
// GetSignBytes implements sdk.Msg. The function will panic since it is used
// for amino transaction verification which IBC does not support.
func (msg MsgConnectionOpenTry) Bytes() []byte {
	panic("IBC messages do not support amino")
}


// ------------MsgConnectionOpenAck

// NewMsgConnectionOpenAck creates a new MsgConnectionOpenAck instance
//nolint:interfacer
func NewMsgConnectionOpenAck(
	connectionID, counterpartyConnectionID string, counterpartyClient exported.ClientState,
	proofTry, proofClient, proofConsensus []byte,
	proofHeight, consensusHeight clienttypes.Height,
	version *Version,
	signer sdk.AccAddress,
) *MsgConnectionOpenAck {
	return &MsgConnectionOpenAck{
		ConnectionId:             connectionID,
		CounterpartyConnectionId: counterpartyConnectionID,
		ClientState:              counterpartyClient,
		ProofTry:                 proofTry,
		ProofClient:              proofClient,
		ProofConsensus:           proofConsensus,
		ProofHeight:              proofHeight,
		ConsensusHeight:          consensusHeight,
		Version:                  version,
		Signer:                   signer.String(),
	}
}

func (m MsgConnectionOpenAck) Route() string {
	return host.RouterKey
}

func (m MsgConnectionOpenAck) MsgType() string {
	return "connection_open_ack"
}

func (m MsgConnectionOpenAck) Type() string {
	return "connection_open_ack"
}

func (msg MsgConnectionOpenAck) ValidateBasic() error {
	if !IsValidConnectionID(msg.ConnectionId) {
		return ErrInvalidConnectionIdentifier
	}
	if err := host.ConnectionIdentifierValidator(msg.CounterpartyConnectionId); err != nil {
		return ErrInvalidParam( "invalid counterparty connection ID")
	}
	if err := ValidateVersion(msg.Version); err != nil {
		return err
	}
	if msg.ClientState == nil {
		return ErrInvalidParam( "counterparty client is nil")
	}
	if err := msg.ClientState.Validate(); err != nil {
		return ErrInvalidParam( "counterparty client is invalid")
	}
	if len(msg.ProofTry) == 0 {
		return ErrInvalidParam( "cannot submit an empty proof try")
	}
	if len(msg.ProofClient) == 0 {
		return ErrInvalidParam( "cannot submit empty proof client")
	}
	if len(msg.ProofConsensus) == 0 {
		return ErrInvalidParam(  "cannot submit an empty proof of consensus state")
	}
	if msg.ProofHeight.IsZero() {
		return ErrInvalidParam( "proof height must be non-zero")
	}
	if msg.ConsensusHeight.IsZero() {
		return ErrInvalidParam(  "consensus height must be non-zero")
	}

	return nil
}

func (m MsgConnectionOpenAck) GetFromAddress() sdk.AccAddress {
	return sdk.HexToAddress(m.Signer)
}

// GetSigners implements sdk.Msg
func (msg MsgConnectionOpenAck) GetSigners() []sdk.AccAddress {
	accAddr := sdk.HexToAddress(msg.Signer)
	return []sdk.AccAddress{accAddr}
}

func (m MsgConnectionOpenAck) Bytes() []byte {
	panic("IBC messages do not support amino")
}


// ------------MsgConnectionOpenConfirm

// NewMsgConnectionOpenConfirm creates a new MsgConnectionOpenConfirm instance
//nolint:interfacer
func NewMsgConnectionOpenConfirm(
	connectionID string, proofAck []byte, proofHeight clienttypes.Height,
	signer sdk.AccAddress,
) *MsgConnectionOpenConfirm {
	return &MsgConnectionOpenConfirm{
		ConnectionId: connectionID,
		ProofAck:     proofAck,
		ProofHeight:  proofHeight,
		Signer:       signer.String(),
	}
}

func (m MsgConnectionOpenConfirm) Route() string {
	return host.RouterKey
}

func (m MsgConnectionOpenConfirm) MsgType() string {
	return "connection_open_confirm"
}

func (m MsgConnectionOpenConfirm) Type() string {
	return "connection_open_confirm"
}

func (msg MsgConnectionOpenConfirm) ValidateBasic() error {
	if !IsValidConnectionID(msg.ConnectionId) {
		return ErrInvalidConnectionIdentifier
	}
	if len(msg.ProofAck) == 0 {
		return ErrInvalidParam( "cannot submit an empty proof ack")
	}
	if msg.ProofHeight.IsZero() {
		return ErrInvalidParam( "proof height must be non-zero")
	}
	//_, err := sdk.AccAddressFromBech32(msg.Signer)
	//if err != nil {
	//	return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "string could not be parsed as address: %v", err)
	//}
	return nil
}

func (msg MsgConnectionOpenConfirm) GetFromAddress() sdk.AccAddress {
	return sdk.HexToAddress(msg.Signer)
}

func (m MsgConnectionOpenConfirm) Bytes() []byte {
	panic("IBC messages do not support amino")
}
// GetSigners implements sdk.Msg
func (msg MsgConnectionOpenConfirm) GetSigners() []sdk.AccAddress {
	accAddr := sdk.HexToAddress(msg.Signer)
	return []sdk.AccAddress{accAddr}
}
