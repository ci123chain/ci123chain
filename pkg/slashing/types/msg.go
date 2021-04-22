package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

// slashing message types
const (
	TypeMsgUnjail = "unjail"
)

// verify interface at compile time
var _ sdk.Msg = &MsgUnjail{}

type MsgUnjail struct {
	ValidatorAddr string `json:"validator_addr"`
}

// NewMsgUnjail creates a new MsgUnjail instance
//nolint:interfacer
func NewMsgUnjail(validatorAddr sdk.AccAddress) *MsgUnjail {
	return &MsgUnjail{
		ValidatorAddr: validatorAddr.String(),
	}
}

func (msg MsgUnjail) Route() string { return RouterKey }
func (msg MsgUnjail) MsgType() string { return TypeMsgUnjail }
func (msg MsgUnjail) GetSigners() []sdk.AccAddress {
	valAddr, err := sdk.AccAddressFromBech32(msg.ValidatorAddr)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{valAddr}
}

func (msg MsgUnjail) GetFromAddress() sdk.AccAddress {
	valAddr, err := sdk.AccAddressFromBech32(msg.ValidatorAddr)
	if err != nil {
		panic(err)
	}
	return valAddr
}

func (msg MsgUnjail) Bytes() []byte {
	bz := SlashingCodec.MustMarshalJSON(&msg)
	return bz
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgUnjail) GetSignBytes() []byte {
	bz := SlashingCodec.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgUnjail) ValidateBasic() error {
	if msg.ValidatorAddr == "" {
		return ErrBadValidatorAddr
	}

	return nil
}
