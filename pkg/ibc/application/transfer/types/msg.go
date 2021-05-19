package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	clienttypes "github.com/ci123chain/ci123chain/pkg/ibc/core/clients/types"
	"github.com/ci123chain/ci123chain/pkg/ibc/core/host"
	cosmosSdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

// msg types
const (
	TypeMsgTransfer = "transfer"
)

var _ sdk.Msg = &MsgTransfer{}
var _ cosmosSdk.Msg = &MsgTransfer{}
// NewMsgTransfer creates a new MsgTransfer instance
//nolint:interfacer
func NewMsgTransfer(
	sourcePort, sourceChannel string,
	token sdk.Coin, sender string, receiver string,
	timeoutHeight clienttypes.Height, timeoutTimestamp uint64,
) *MsgTransfer {
	return &MsgTransfer{
		SourcePort:       sourcePort,
		SourceChannel:    sourceChannel,
		Token:            token,
		Sender:           sender,
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

func (msg MsgTransfer) ValidateBasic() error {
	if err := host.PortIdentifierValidator(msg.SourcePort); err != nil {
		return ErrInvalidPortID
	}
	if err := host.ChannelIdentifierValidator(msg.SourceChannel); err != nil {
		return ErrInvalidChannelID
	}
	if !msg.Token.IsValid() {
		return ErrInvalidToken
	}
	if !msg.Token.IsPositive() {
		return ErrInvalidToken
	}
	// NOTE: sender format must be validated as it is required by the GetSigners function.
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return ErrInvalidSender
	}
	if strings.TrimSpace(msg.Receiver) == "" {
		return ErrInvalidReceiver
	}
	return ValidateIBCDenom(msg.Token.Denom)
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

func (t MsgTransfer) GetSignBytes() []byte {
	return sdk.MustSortJSON(IBCTransferCdc.MustMarshalJSON(&t))
}

func (t MsgTransfer) GetSigners() []cosmosSdk.AccAddress {
	return []cosmosSdk.AccAddress{sdk.HexToAddress(t.Sender).Bytes()}
}

func (t MsgTransfer) Type() string {
	return "channel_open_init"
}