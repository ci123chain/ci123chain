package main

import (
	"encoding/hex"
	sdk "github.com/ci123chain/ci123chain/sdk/ibc"
)


func SignIBC(from, to string, amount uint64, priv string, denom string) ([]byte, error) {
	privateKey, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}
	txByte, err := sdk.SignIBCTransferMsg(from, to, amount, privateKey, denom)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}

func SignIBCApplyTx(from string, uniqueID, observerID []byte, priv string) ([]byte, error) {
	privateKey, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}
	txByte, err := sdk.SignApplyIBCMsg(from, uniqueID, observerID, privateKey)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}

func SignIBCBankSendTx(from string, raw []byte, priv string) ([]byte, error) {
	privateKey, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}
	txByte, err := sdk.SignIBCBankSendMsg(from, raw, privateKey)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}

func SignIBCReceiptTx(from string, raw []byte, priv string) ([]byte, error) {
	privateKey, err := hex.DecodeString(priv)
	if err != nil {
		return nil, err
	}
	txByte, err := sdk.SignIBCReceiptMsg(from, raw, privateKey)
	if err != nil {
		return nil, err
	}
	return txByte, nil
}


