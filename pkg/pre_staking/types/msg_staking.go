package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)

type MsgStaking struct {
	FromAddress    sdk.AccAddress   `json:"from_address"`
	Delegator      sdk.AccAddress   `json:"delegator"`
	Validator      sdk.AccAddress   `json:"validator"`
	Amount         sdk.Coin         `json:"amount"`
}

func NewMsgStaking(from sdk.AccAddress, delegatorAddr sdk.AccAddress, validatorAddr sdk.AccAddress,
	amount sdk.Coin) *MsgStaking {
		return &MsgStaking{
			FromAddress: from,
			Delegator:   delegatorAddr,
			Validator:   validatorAddr,
			Amount:      amount,
		}
}

func (msg *MsgStaking) ValidateBasic() error {
	if msg.Delegator.Empty() {
		return ErrInvalidDelegatorAddress
	}
	if msg.Validator.Empty() {
		return ErrInvalidValidatorAddress
	}
	if !msg.Amount.Amount.IsPositive() {
		return ErrInvalidAmount
	}
	if !msg.FromAddress.Equal(msg.Delegator) {
		return ErrFromNotEqualDelegator
	}
	return nil
}

func (msg *MsgStaking) Bytes() []byte {
	bytes, err := PreStakingCodec.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *MsgStaking) Route() string {return ModuleName}
func (msg *MsgStaking) MsgType() string {return "delegate"}
func (msg *MsgStaking) GetFromAddress() sdk.AccAddress { return msg.FromAddress}