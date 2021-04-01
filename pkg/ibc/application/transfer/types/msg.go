package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
)

// msg types
const (
	TypeMsgTransfer = "transfer"
)

var _ sdk.Msg = &MsgTransfer{}
// NewMsgTransfer creates a new MsgTransfer instance
//nolint:interfacer
func NewMsgTransfer(
	sourcePort, sourceChannel string,
	token sdk.Coin, sender sdk.AccAddress, receiver string,
	timeoutHeight clienttypes.Height, timeoutTimestamp uint64,
) *MsgTransfer {
	return &MsgTransfer{
		SourcePort:       sourcePort,
		SourceChannel:    sourceChannel,
		Token:            token,
		Sender:           sender.String(),
		Receiver:         receiver,
		TimeoutHeight:    timeoutHeight,
		TimeoutTimestamp: timeoutTimestamp,
	}
}
// Route implements sdk.Msg
func (MsgTransfer) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (MsgTransfer) MsgType() string {
	return TypeMsgTransfer
}

func (t MsgTransfer) ValidateBasic() sdk.Error {
	//if err := host.PortIdentifierValidator(msg.SourcePort); err != nil {
	//	return sdkerrors.Wrap(err, "invalid source port ID")
	//}
	//if err := host.ChannelIdentifierValidator(msg.SourceChannel); err != nil {
	//	return sdkerrors.Wrap(err, "invalid source channel ID")
	//}
	//if !msg.Token.IsValid() {
	//	return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Token.String())
	//}
	//if !msg.Token.IsPositive() {
	//	return sdkerrors.Wrap(sdkerrors.ErrInsufficientFunds, msg.Token.String())
	//}
	//// NOTE: sender format must be validated as it is required by the GetSigners function.
	//_, err := sdk.AccAddressFromBech32(msg.Sender)
	//if err != nil {
	//	return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "string could not be parsed as address: %v", err)
	//}
	//if strings.TrimSpace(msg.Receiver) == "" {
	//	return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing recipient address")
	//}
	//return ValidateIBCDenom(msg.Token.Denom)

	return nil
}

func (t MsgTransfer) GetFromAddress() sdk.AccAddress {
	return sdk.HexToAddress(t.Sender)
}

func (t MsgTransfer) Bytes() []byte {
	bytes, err := IBCTransferCdc.MarshalBinaryLengthPrefixed(t)
	if err != nil {
		panic(err)
	}
	return bytes
}
