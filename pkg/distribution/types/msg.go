package types

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/util"
)


type MsgFundCommunityPool struct {
	FromAddress		  	sdk.AccAddress	 `json:"from_address"`
	Signature 		  	[]byte   		 `json:"signature"`
	PubKey			  	[]byte			 `json:"pub_key"`

	Amount       sdk.Coin          `json:"amount"`
	Depositor    sdk.AccAddress    `json:"depositor"`
}

func NewMsgFundCommunityPool(from sdk.AccAddress,amount sdk.Coin, gas, nonce uint64, depositor sdk.AccAddress) MsgFundCommunityPool {
	return MsgFundCommunityPool{
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

func (msg *MsgFundCommunityPool) SetSignature(sig []byte) {
	msg.Signature = sig
}

func (msg *MsgFundCommunityPool) SetPubKey(pub []byte) {
	msg.PubKey = pub
}

// GetSignBytes returns the raw bytes for a MsgFundCommunityPool message that
// the expected signer needs to sign.
func (msg *MsgFundCommunityPool) GetSignBytes() []byte {
	tmsg := *msg
	tmsg.Signature = nil
	return util.TxHash(tmsg.Bytes())
}

// ValidateBasic performs basic MsgFundCommunityPool message validation.
func (msg *MsgFundCommunityPool) ValidateBasic() sdk.Error {
	if !msg.Amount.IsValid() {
		return ErrInvalidCoin(DefaultCodespace, msg.Amount.String())
	}
	if msg.Depositor.Empty() {
		return ErrInvalidAddress(DefaultCodespace, msg.Depositor.String())
	}
	if !msg.FromAddress.Equal(msg.Depositor) {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("expected %s, got %s", msg.FromAddress.String(), msg.Depositor.String()))
	}

	return nil
}

func (msg *MsgFundCommunityPool) GetFromAddress() sdk.AccAddress { return msg.FromAddress}

func (msg *MsgFundCommunityPool) GetSignature() []byte { return msg.Signature}

type MsgSetWithdrawAddress struct {
	FromAddress		  	sdk.AccAddress	 `json:"from_address"`
	Signature 		  	[]byte   		 `json:"signature"`
	PubKey			  	[]byte			 `json:"pub_key"`

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

func (msg *MsgSetWithdrawAddress) ValidateBasic() sdk.Error {
	if msg.DelegatorAddress.Empty() {
		return ErrInvalidAddress(DefaultCodespace, msg.DelegatorAddress.String())
	}
	if msg.WithdrawAddress.Empty() {
		return ErrInvalidAddress(DefaultCodespace, msg.WithdrawAddress.String())
	}
	if msg.FromAddress.Empty() {
		return ErrInvalidAddress(DefaultCodespace, msg.FromAddress.String())
	}
	//keep delegator address and from address the same.
	if !msg.FromAddress.Equal(msg.DelegatorAddress) {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("expected %s, got %s", msg.FromAddress.String(), msg.DelegatorAddress.String()))
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

func (msg *MsgSetWithdrawAddress) SetSignature(sig []byte) {
	msg.Signature = sig
}

func (msg *MsgSetWithdrawAddress) SetPubKey(pubKey []byte) {
	msg.PubKey = pubKey
}

func (msg *MsgSetWithdrawAddress) GetSignBytes() []byte {
	tmsg := *msg
	tmsg.Signature = nil
	return util.TxHash(tmsg.Bytes())
}

func (msg *MsgSetWithdrawAddress) GetFromAddress() sdk.AccAddress { return msg.FromAddress}

func (msg *MsgSetWithdrawAddress) GetSignature() []byte { return msg.Signature}


type MsgWithdrawDelegatorReward struct {
	FromAddress		  	sdk.AccAddress	 `json:"from_address"`
	Signature 		  	[]byte   		 `json:"signature"`
	PubKey			  	[]byte			 `json:"pub_key"`

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

func (msg *MsgWithdrawDelegatorReward) ValidateBasic() sdk.Error {
	if msg.ValidatorAddress.Empty() {
		return ErrInvalidAddress(DefaultCodespace, msg.ValidatorAddress.String())
	}
	if msg.DelegatorAddress.Empty() {
		return ErrInvalidAddress(DefaultCodespace, msg.DelegatorAddress.String())
	}
	if msg.FromAddress.Empty() {
		return ErrInvalidAddress(DefaultCodespace, msg.FromAddress.String())
	}
	if !msg.FromAddress.Equal(msg.DelegatorAddress) {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("expected %s, got %s", msg.FromAddress.String(), msg.DelegatorAddress.String()))
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

func (msg *MsgWithdrawDelegatorReward) SetSignature(sig []byte) {
	msg.Signature = sig
}

func (msg *MsgWithdrawDelegatorReward) SetPubKey(pubKey []byte) {
	msg.PubKey = pubKey
}

func (msg *MsgWithdrawDelegatorReward) GetSignBytes() []byte {
	tmsg := *msg
	tmsg.Signature = nil
	return util.TxHash(tmsg.Bytes())
}

func (msg *MsgWithdrawDelegatorReward) GetSignature() []byte {
	return msg.Signature
}

func (msg *MsgWithdrawDelegatorReward) GetFromAddress() sdk.AccAddress { return msg.FromAddress}


type MsgWithdrawValidatorCommission struct {
	FromAddress		  	sdk.AccAddress	 `json:"from_address"`
	Signature 		  	[]byte   		 `json:"signature"`
	PubKey			  	[]byte			 `json:"pub_key"`

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

func (msg *MsgWithdrawValidatorCommission) ValidateBasic() sdk.Error {
	if msg.ValidatorAddress.Empty() {
		return ErrInvalidAddress(DefaultCodespace, msg.ValidatorAddress.String())
	}
	if msg.FromAddress.Empty() {
		return ErrInvalidAddress(DefaultCodespace, msg.FromAddress.String())
	}
	if !msg.FromAddress.Equal(msg.ValidatorAddress) {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("expected %s, got %s", msg.FromAddress.String(), msg.ValidatorAddress.String()))
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

func (msg *MsgWithdrawValidatorCommission) SetSignature(sig []byte) {
	msg.Signature = sig
}

func (msg *MsgWithdrawValidatorCommission) SetPubKey(pubKey []byte) {
	msg.PubKey = pubKey
}

func (msg *MsgWithdrawValidatorCommission) GetSignBytes() []byte {
	tmsg := *msg
	tmsg.Signature = nil
	return util.TxHash(tmsg.Bytes())
}

func (msg *MsgWithdrawValidatorCommission) GetSignature() []byte {
	return msg.Signature
}

func (msg *MsgWithdrawValidatorCommission) GetFromAddress() sdk.AccAddress { return msg.FromAddress}