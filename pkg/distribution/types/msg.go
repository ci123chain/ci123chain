package types

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/transaction"
	"github.com/ci123chain/ci123chain/pkg/util"
)


type FundCommunityPoolTx struct {
	transaction.CommonTx
	Amount       sdk.Coin          `json:"amount"`
	Depositor    sdk.AccAddress    `json:"depositor"`
}


func NewMsgFundCommunityPool(from sdk.AccAddress,amount sdk.Coin, gas, nonce uint64, depositor sdk.AccAddress) FundCommunityPoolTx {
	return FundCommunityPoolTx{
		CommonTx:transaction.CommonTx{
			From:      from,
			Nonce:     nonce,
			Gas:       gas,
		},
		Amount:    amount,
		Depositor: depositor,
	}
}

// Route returns the MsgFundCommunityPool message route.
func (msg *FundCommunityPoolTx) Route() string { return RouteKey }

// GetSigners returns the signer addresses that are expected to sign the result
// of GetSignBytes.
func (msg *FundCommunityPoolTx) Bytes() []byte{
	bytes, err := DistributionCdc.MarshalBinaryLengthPrefixed(msg)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (msg *FundCommunityPoolTx) SetSignature(sig []byte) {
	msg.Signature = sig
}

func (msg *FundCommunityPoolTx) SetPubKey(pub []byte) {
	msg.PubKey = pub
}

// GetSignBytes returns the raw bytes for a MsgFundCommunityPool message that
// the expected signer needs to sign.
func (msg *FundCommunityPoolTx) GetSignBytes() []byte {
	tmsg := *msg
	tmsg.Signature = nil
	return util.TxHash(tmsg.Bytes())
}

// ValidateBasic performs basic MsgFundCommunityPool message validation.
func (msg *FundCommunityPoolTx) ValidateBasic() sdk.Error {

	err := msg.VerifySignature(msg.GetSignBytes(), false)
	if err != nil {
		return ErrInvalidSignature(DefaultCodespace, "invalid signature")
	}
	if !msg.Amount.IsValid() {
		return ErrInvalidCoin(DefaultCodespace, msg.Amount.String())
	}
	if msg.Depositor.Empty() {
		return ErrInvalidAddress(DefaultCodespace, msg.Depositor.String())
	}
	if !msg.From.Equal(msg.Depositor) {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("expected %s, got %s", msg.From.String(), msg.Depositor.String()))
	}

	return nil
}

func (msg *FundCommunityPoolTx) GetNonce() uint64 { return msg.Nonce}
func (msg *FundCommunityPoolTx) GetGas() uint64 { return msg.Gas}
func (msg *FundCommunityPoolTx) GetFromAddress() sdk.AccAddress { return msg.From}



type SetWithdrawAddressTx struct {
	transaction.CommonTx
	DelegatorAddress     sdk.AccAddress   `json:"delegator_address"`
	WithdrawAddress      sdk.AccAddress    `json:"withdraw_address"`
}

func NewSetWithdrawAddressTx(from, withdraw, del sdk.AccAddress, gas, nonce uint64) SetWithdrawAddressTx{
	return SetWithdrawAddressTx{
		CommonTx: transaction.CommonTx{
			From:      from,
			Nonce:     nonce,
			Gas:       gas,
		},
		DelegatorAddress:del,
		WithdrawAddress:withdraw,
	}
}

func (tx *SetWithdrawAddressTx) Route() string { return RouteKey}

func (tx *SetWithdrawAddressTx) ValidateBasic() sdk.Error {
	err := tx.VerifySignature(tx.GetSignBytes(), false)
	if err != nil {
		return ErrInvalidSignature(DefaultCodespace, "invalid signature")
	}
	if tx.DelegatorAddress.Empty() {
		return ErrInvalidAddress(DefaultCodespace, tx.DelegatorAddress.String())
	}
	if tx.WithdrawAddress.Empty() {
		return ErrInvalidAddress(DefaultCodespace, tx.WithdrawAddress.String())
	}
	if tx.From.Empty() {
		return ErrInvalidAddress(DefaultCodespace, tx.From.String())
	}
	//keep delegator address and from address the same.
	if !tx.From.Equal(tx.DelegatorAddress) {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("expected %s, got %s", tx.From.String(), tx.DelegatorAddress.String()))
	}
	return nil
}

func (tx *SetWithdrawAddressTx) Bytes() []byte {
	bytes, err := DistributionCdc.MarshalBinaryLengthPrefixed(tx)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (tx *SetWithdrawAddressTx) SetSignature(sig []byte) {
	tx.Signature = sig
}

func (tx *SetWithdrawAddressTx) SetPubKey(pubKey []byte) {
	tx.PubKey = pubKey
}

func (tx *SetWithdrawAddressTx) GetSignBytes() []byte {
	tmsg := *tx
	tmsg.Signature = nil
	return util.TxHash(tmsg.Bytes())
}

func (tx *SetWithdrawAddressTx) GetNonce() uint64 {
	return tx.Nonce
}

func (tx *SetWithdrawAddressTx) GetGas() uint64 {
	return tx.Gas
}

func (tx *SetWithdrawAddressTx) GetFromAddress() sdk.AccAddress { return tx.From}


type WithdrawDelegatorRewardTx struct {
	transaction.CommonTx
	DelegatorAddress     sdk.AccAddress    `json:"delegator_address"`
	ValidatorAddress     sdk.AccAddress    `json:"validator_address"`
}

func NewWithdrawDelegatorRewardTx(from, val, del sdk.AccAddress, gas, nonce uint64) WithdrawDelegatorRewardTx {
	return WithdrawDelegatorRewardTx{
		CommonTx:transaction.CommonTx{
			From:      from,
			Nonce:     nonce,
			Gas:       gas,
		},
		DelegatorAddress:del,
		ValidatorAddress:val,
	}
}

func (tx *WithdrawDelegatorRewardTx) Route() string { return RouteKey}

func (tx *WithdrawDelegatorRewardTx) ValidateBasic() sdk.Error {
	err := tx.VerifySignature(tx.GetSignBytes(), false)
	if err != nil {
		return ErrInvalidSignature(DefaultCodespace, "invalid signature")
	}
	if tx.ValidatorAddress.Empty() {
		return ErrInvalidAddress(DefaultCodespace, tx.ValidatorAddress.String())
	}
	if tx.DelegatorAddress.Empty() {
		return ErrInvalidAddress(DefaultCodespace, tx.DelegatorAddress.String())
	}
	if tx.From.Empty() {
		return ErrInvalidAddress(DefaultCodespace, tx.From.String())
	}
	if !tx.From.Equal(tx.DelegatorAddress) {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("expected %s, got %s", tx.From.String(), tx.DelegatorAddress.String()))
	}

	return nil
}

func (tx *WithdrawDelegatorRewardTx) Bytes() []byte {
	bytes, err := DistributionCdc.MarshalBinaryLengthPrefixed(tx)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (tx *WithdrawDelegatorRewardTx) SetSignature(sig []byte) {
	tx.Signature = sig
}

func (tx *WithdrawDelegatorRewardTx) SetPubKey(pubKey []byte) {
	tx.PubKey = pubKey
}

func (tx *WithdrawDelegatorRewardTx) GetSignBytes() []byte {
	tmsg := *tx
	tmsg.Signature = nil
	return util.TxHash(tmsg.Bytes())
}

func (tx *WithdrawDelegatorRewardTx) GetNonce() uint64 {
	return tx.Nonce
}

func (tx *WithdrawDelegatorRewardTx) GetGas() uint64 {
	return tx.Gas
}

func (tx *WithdrawDelegatorRewardTx) GetFromAddress() sdk.AccAddress { return tx.From}


type WithdrawValidatorCommissionTx struct {
	transaction.CommonTx
	ValidatorAddress    sdk.AccAddress    `json:"validator_address"`
}

func NewWithdrawValidatorCommissionTx(from, val sdk.AccAddress, gas, nonce uint64) WithdrawValidatorCommissionTx {
	return WithdrawValidatorCommissionTx{
		CommonTx:transaction.CommonTx{
			From:      from,
			Nonce:     nonce,
			Gas:       gas,
		},
		ValidatorAddress:val,
	}
}

func (tx *WithdrawValidatorCommissionTx) Route() string { return RouteKey}

func (tx *WithdrawValidatorCommissionTx) ValidateBasic() sdk.Error {
	err := tx.VerifySignature(tx.GetSignBytes(), false)
	if err != nil {
		return ErrInvalidSignature(DefaultCodespace, "invalid signature")
	}
	if tx.ValidatorAddress.Empty() {
		return ErrInvalidAddress(DefaultCodespace, tx.ValidatorAddress.String())
	}
	if tx.From.Empty() {
		return ErrInvalidAddress(DefaultCodespace, tx.From.String())
	}
	if !tx.From.Equal(tx.ValidatorAddress) {
		return ErrInvalidAddress(DefaultCodespace, fmt.Sprintf("expected %s, got %s", tx.From.String(), tx.ValidatorAddress.String()))
	}

	return nil
}

func (tx *WithdrawValidatorCommissionTx) Bytes() []byte {
	bytes, err := DistributionCdc.MarshalBinaryLengthPrefixed(tx)
	if err != nil {
		panic(err)
	}

	return bytes
}

func (tx *WithdrawValidatorCommissionTx) SetSignature(sig []byte) {
	tx.Signature = sig
}

func (tx *WithdrawValidatorCommissionTx) SetPubKey(pubKey []byte) {
	tx.PubKey = pubKey
}

func (tx *WithdrawValidatorCommissionTx) GetSignBytes() []byte {
	tmsg := *tx
	tmsg.Signature = nil
	return util.TxHash(tmsg.Bytes())
}

func (tx *WithdrawValidatorCommissionTx) GetNonce() uint64 {
	return tx.Nonce
}

func (tx *WithdrawValidatorCommissionTx) GetGas() uint64 {
	return tx.Gas
}

func (tx *WithdrawValidatorCommissionTx) GetFromAddress() sdk.AccAddress { return tx.From}