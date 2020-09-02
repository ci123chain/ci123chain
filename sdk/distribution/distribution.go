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


func SignSetWithdrawAddressTx(from, withdrawAddress string, gas, nonce uint64, priv string) ([]byte, error) {

	privateKey, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}
	accountAddr := sdk.HexToAddress(from)
	withdrawAddr := sdk.HexToAddress(withdrawAddress)

	tx := types.NewSetWithdrawAddressTx(accountAddr, withdrawAddr, accountAddr, gas, nonce)

	eth := cryptosuit.NewETHSignIdentity()
	signature, err := eth.Sign(tx.GetSignBytes(), privateKey)
	if err != nil {
		return nil, err
	}
	tx.SetSignature(signature)

	return tx.Bytes(), nil
}

func SignWithdrawDelegatorRewardTx(from, validatorAddress, delegatorAddress string, gas, nonce uint64, priv string) ([]byte, error) {
	privateKey, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}
	accountAddr := sdk.HexToAddress(from)
	valAddr := sdk.HexToAddress(validatorAddress)
	delAddr := sdk.HexToAddress(delegatorAddress)

	tx := types.NewWithdrawDelegatorRewardTx(accountAddr, valAddr, delAddr, gas, nonce)
	/*sid := cryptosuit.NewFabSignIdentity()
	pub, err  := sid.GetPubKey(privateKey)
	if err != nil {
		return nil, err
	}

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

func SignWithdrawValidatorCommissionTx(from, validatorAddress string, gas, nonce uint64, priv string) ([]byte, error) {
	privateKey, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}
	accountAddr := sdk.HexToAddress(from)
	valAddr := sdk.HexToAddress(validatorAddress)

	tx := types.NewWithdrawValidatorCommissionTx(accountAddr, valAddr, gas, nonce)
	/*sid := cryptosuit.NewFabSignIdentity()
	pub, err  := sid.GetPubKey(privateKey)
	if err != nil {
		return nil, err
	}

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