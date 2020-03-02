package main

import (
	"encoding/hex"
	sdk "github.com/tanhuiya/ci123chain/sdk/ibc"
)


func SignIBC(from, to string, amount, gas, nonce uint64, priv string) ([]byte, error) {
	privateKey, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}
	txByte, err := sdk.SignIBCTransferMsg(from, to, amount, gas, nonce, privateKey)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}

func SignIBCApplyTx(from string, uniqueID, observerID []byte, gas, nonce uint64, priv string) ([]byte, error) {
	privateKey, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}
	txByte, err := sdk.SignApplyIBCMsg(from, uniqueID, observerID, gas, nonce, privateKey)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}

func SignIBCBankSendTx(from string, raw []byte, gas, nonce uint64, priv string) ([]byte, error) {
	privateKey, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}
	txByte, err := sdk.SignIBCBankSendMsg(from, raw, gas, nonce, privateKey)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}

func SignIBCReceiptTx(from string, raw []byte, gas, nonce uint64, priv string) ([]byte, error) {
	privateKey, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}
	txByte, err := sdk.SignIBCReceiptMsg(from, raw, gas, nonce, privateKey)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}


