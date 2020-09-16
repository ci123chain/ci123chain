package distribution

import (
	"encoding/hex"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"github.com/ci123chain/ci123chain/pkg/distribution/types"
)

func SignFundCommunityPoolTx(from string, amount int64, gas, nonce uint64, priv string) ([]byte, error) {
	//
	privateKey, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}
	Amount := sdk.NewCoin(sdk.NewInt(amount))
	accountAddr := sdk.HexToAddress(from)
	tx := types.NewMsgFundCommunityPool(accountAddr, Amount, gas, nonce, accountAddr)

	/*sid := cryptosuit.NewFabSignIdentity()
	pub, err  := sid.GetPubKey(privateKey)

	tx.SetPubKey(pub)
	signbyte := tx.GetSignBytes()
	signature, err := sid.Sign(signbyte, privateKey)*/
	eth := cryptosuit.NewETHSignIdentity()
	signature, err := eth.Sign(tx.GetSignBytes(), privateKey)
	if err != nil {
		return nil, err
	}
	tx.SetSignature(signature)

	return tx.Bytes(), nil
}


func SignMsgSetWithdrawAddress(from, withdrawAddress sdk.AccAddress, priv string) (sdk.Msg, error) {

	privateKey, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}

	msg := types.NewMsgSetWithdrawAddress(from, withdrawAddress, from)

	eth := cryptosuit.NewETHSignIdentity()
	signature, err := eth.Sign(msg.GetSignBytes(), privateKey)
	if err != nil {
		return nil, err
	}
	msg.SetSignature(signature)

	return msg, nil
}

func SignWithdrawDelegatorRewardTx(from, validatorAddress, delegatorAddress sdk.AccAddress, priv string) (sdk.Msg, error) {
	privateKey, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}

	msg := types.NewMsgWithdrawDelegatorReward(from, validatorAddress, delegatorAddress)
	/*sid := cryptosuit.NewFabSignIdentity()
	pub, err  := sid.GetPubKey(privateKey)
	if err != nil {
		return nil, err
	}

	tx.SetPubKey(pub)
	signbyte := tx.GetSignBytes()
	signature, err := sid.Sign(signbyte, privateKey)*/
	eth := cryptosuit.NewETHSignIdentity()
	signature, err := eth.Sign(msg.GetSignBytes(), privateKey)
	if err != nil {
		return nil, err
	}
	msg.SetSignature(signature)

	return msg, nil
}

func SignWithdrawValidatorCommissionTx(from, validatorAddress sdk.AccAddress, priv string) (sdk.Msg, error) {
	privateKey, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}

	msg := types.NewMsgWithdrawValidatorCommission(from, validatorAddress)
	/*sid := cryptosuit.NewFabSignIdentity()
	pub, err  := sid.GetPubKey(privateKey)
	if err != nil {
		return nil, err
	}

	tx.SetPubKey(pub)
	signbyte := tx.GetSignBytes()
	signature, err := sid.Sign(signbyte, privateKey)*/
	eth := cryptosuit.NewETHSignIdentity()
	signature, err := eth.Sign(msg.GetSignBytes(), privateKey)
	if err != nil {
		return nil, err
	}
	msg.SetSignature(signature)

	return msg, nil
}