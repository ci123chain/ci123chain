package distribution

import (
	"encoding/hex"
	"errors"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/app"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"github.com/ci123chain/ci123chain/pkg/distribution/types"
)

var cdc = app.MakeCodec()

//todo
func SignFundCommunityPoolTx(from string, amount sdk.Coin, gas, nonce uint64, priv string) ([]byte, error) {
	//
	privateKey, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}
	//Amount := sdk.NewCoin(sdk.NewInt(amount))
	if amount.IsNegative() || amount.IsZero() {
		return nil, errors.New("invalid amount")
	}
	accountAddr := sdk.HexToAddress(from)
	tx := types.NewMsgFundCommunityPool(accountAddr, amount, gas, nonce, accountAddr)
	eth := cryptosuit.NewETHSignIdentity()
	signature, err := eth.Sign(tx.GetSignBytes(), privateKey)
	if err != nil {
		return nil, err
	}
	tx.SetSignature(signature)

	return tx.Bytes(), nil
}

func SignMsgSetWithdrawAddress(from, withdrawAddress sdk.AccAddress, gas, nonce uint64, priv string) ([]byte, error) {
	msg := types.NewMsgSetWithdrawAddress(from, withdrawAddress, from)
	txByte, err := app.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, priv, cdc)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}

func NewMsgSetWithdrawAddressMsg(from, withdrawAddress sdk.AccAddress) []byte {
	msg := types.NewMsgSetWithdrawAddress(from, withdrawAddress, from)
	return msg.Bytes()
}

func SignWithdrawDelegatorRewardTx(from, validatorAddress, delegatorAddress sdk.AccAddress, gas, nonce uint64, priv string) ([]byte, error) {
	msg := types.NewMsgWithdrawDelegatorReward(from, validatorAddress, delegatorAddress)
	txByte, err := app.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, priv, cdc)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}

func NewWithdrawDelegatorRewardMsg(from, validatorAddress, delegatorAddress sdk.AccAddress) []byte {
	msg := types.NewMsgWithdrawDelegatorReward(from, validatorAddress, delegatorAddress)
	return msg.Bytes()
}

func SignWithdrawValidatorCommissionTx(from, validatorAddress sdk.AccAddress, gas, nonce uint64, priv string) ([]byte, error) {
	msg := types.NewMsgWithdrawValidatorCommission(from, validatorAddress)
	txByte, err := app.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, priv, cdc)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}

func NewWithdrawValidatorCommissionMsg(from, validatorAddress sdk.AccAddress) []byte {
	msg := types.NewMsgWithdrawValidatorCommission(from, validatorAddress)
	return msg.Bytes()
}