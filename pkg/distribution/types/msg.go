package types

import (
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
)


type MsgFundCommunityPool struct {
	FromAddress		  	sdk.AccAddress	 `json:"from_address"`
	Amount       sdk.Coin          `json:"amount"`
	Depositor    sdk.AccAddress    `json:"depositor"`
}

func NewMsgFundCommunityPool(from sdk.AccAddress,amount sdk.Coin, _, _ uint64, depositor sdk.AccAddress) *MsgFundCommunityPool {
	return &MsgFundCommunityPool{
		FromAddress: from,
		Amount:    amount,
		Depositor: depositor,
	}
}

// Route returns the MsgFundCommunityPool message route.
func (msg *MsgFundCommunityPool) Route() string { return RouteKey }

func (msg *MsgFundCommunityPool) MsgType() string { return "fund_community_pool" }

func (msg *MsgFundCommunityPool) Bytes() []byte{
	bytes, err := DistributionCdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

// ValidateBasic performs basic MsgFundCommunityPool message validation.
func (msg *MsgFundCommunityPool) ValidateBasic() error {
	if !msg.Amount.IsValid() {
		return ErrInvalidParams
	}
	if msg.Depositor.Empty() {
		return ErrInvalidParams
	}
	if !msg.FromAddress.Equal(msg.Depositor) {
		return ErrInvalidParams
	}

	return nil
}

func (msg *MsgFundCommunityPool) GetFromAddress() sdk.AccAddress { return msg.FromAddress}

type MsgSetWithdrawAddress struct {
	FromAddress		  	sdk.AccAddress	 `json:"from_address"`
	DelegatorAddress     sdk.AccAddress  `json:"delegator_address"`
	WithdrawAddress      sdk.AccAddress  `json:"withdraw_address"`
}

func NewMsgSetWithdrawAddress(from, withdraw, del sdk.AccAddress) *MsgSetWithdrawAddress{
	return &MsgSetWithdrawAddress{
		FromAddress:      from,
		DelegatorAddress: del,
		WithdrawAddress:  withdraw,
	}
}

func (msg *MsgSetWithdrawAddress) Route() string { return RouteKey}

func (msg *MsgSetWithdrawAddress) MsgType() string { return "set_withdraw_address"}

func (msg *MsgSetWithdrawAddress) ValidateBasic() error {
	if msg.DelegatorAddress.Empty() {
		return ErrInvalidParams
	}
	if msg.WithdrawAddress.Empty() {
		return ErrInvalidParams
	}
	if msg.FromAddress.Empty() {
		return ErrInvalidParams
	}
	//keep delegator address and from address the same.
	if !msg.FromAddress.Equal(msg.DelegatorAddress) {
		return ErrInvalidParams
	}
	return nil
}

func (msg *MsgSetWithdrawAddress) Bytes() []byte {
	bytes, err := DistributionCdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}


func (msg *MsgSetWithdrawAddress) GetFromAddress() sdk.AccAddress { return msg.FromAddress}

type MsgWithdrawDelegatorReward struct {
	FromAddress		  	sdk.AccAddress	 `json:"from_address"`
	DelegatorAddress     sdk.AccAddress    `json:"delegator_address"`
	ValidatorAddress     sdk.AccAddress    `json:"validator_address"`
}

func NewMsgWithdrawDelegatorReward(from, val, del sdk.AccAddress) *MsgWithdrawDelegatorReward {
	return &MsgWithdrawDelegatorReward{
		FromAddress: from,
		DelegatorAddress:del,
		ValidatorAddress:val,
	}
}

func (msg *MsgWithdrawDelegatorReward) Route() string { return RouteKey}
func (msg *MsgWithdrawDelegatorReward) MsgType() string { return "withdraw_delegator_reward"}

func (msg *MsgWithdrawDelegatorReward) ValidateBasic() error {
	if msg.ValidatorAddress.Empty() {
		return ErrInvalidParams
	}
	if msg.DelegatorAddress.Empty() {
		return ErrInvalidParams
	}
	if msg.FromAddress.Empty() {
		return ErrInvalidParams
	}
	if !msg.FromAddress.Equal(msg.DelegatorAddress) {
		return ErrInvalidParams
	}

	return nil
}

func (msg *MsgWithdrawDelegatorReward) Bytes() []byte {
	bytes, err := DistributionCdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *MsgWithdrawDelegatorReward) GetFromAddress() sdk.AccAddress { return msg.FromAddress}


type MsgWithdrawValidatorCommission struct {
	FromAddress		  	sdk.AccAddress	 `json:"from_address"`
	ValidatorAddress    sdk.AccAddress   `json:"validator_address"`
}

func NewMsgWithdrawValidatorCommission(from, val sdk.AccAddress) *MsgWithdrawValidatorCommission {
	return &MsgWithdrawValidatorCommission{
		FromAddress:      from,
		ValidatorAddress: val,
	}
}

func (msg *MsgWithdrawValidatorCommission) Route() string { return RouteKey}
func (msg *MsgWithdrawValidatorCommission) MsgType() string { return "withdraw_validator_commission"}

func (msg *MsgWithdrawValidatorCommission) ValidateBasic() error {
	if msg.ValidatorAddress.Empty() {
		return ErrInvalidParams
	}
	if msg.FromAddress.Empty() {
		return ErrInvalidParams
	}
	if !msg.FromAddress.Equal(msg.ValidatorAddress) {
		return ErrInvalidParams
	}

	return nil
}

func (msg *MsgWithdrawValidatorCommission) Bytes() []byte {
	bytes, err := DistributionCdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *MsgWithdrawValidatorCommission) GetFromAddress() sdk.AccAddress { return msg.FromAddress}