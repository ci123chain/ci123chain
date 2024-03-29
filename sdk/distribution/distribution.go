package distribution

import (
	"errors"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	types2 "github.com/ci123chain/ci123chain/pkg/app/types"
	"github.com/ci123chain/ci123chain/pkg/distribution/types"
)

var cdc = types2.GetCodec()

//todo
func SignFundCommunityPoolTx(from string, amount sdk.Coin, gas, nonce uint64, priv string) ([]byte, error) {
	//Amount := sdk.NewCoin(sdk.NewInt(amount))
	if amount.IsNegative() || amount.IsZero() {
		return nil, errors.New("invalid amount")
	}
	accountAddr := sdk.HexToAddress(from)
	msg := types.NewMsgFundCommunityPool(accountAddr, amount, gas, nonce, accountAddr)
	txByte, err := types2.SignCommonTx(accountAddr, nonce, gas, []sdk.Msg{msg}, priv, cdc)
	if err != nil {
		return nil, err
	}

	return txByte, nil
}

func SignMsgSetWithdrawAddress(from, withdrawAddress sdk.AccAddress, gas, nonce uint64, priv string) ([]byte, error) {
	msg := types.NewMsgSetWithdrawAddress(from, withdrawAddress, from)
	txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, priv, cdc)
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
	txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, priv, cdc)
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
	txByte, err := types2.SignCommonTx(from, nonce, gas, []sdk.Msg{msg}, priv, cdc)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}

func NewWithdrawValidatorCommissionMsg(from, validatorAddress sdk.AccAddress) []byte {
	msg := types.NewMsgWithdrawValidatorCommission(from, validatorAddress)
	return msg.Bytes()
}