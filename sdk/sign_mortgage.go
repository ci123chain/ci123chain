package sdk

import (
	"github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/client/helper"
	"github.com/tanhuiya/ci123chain/pkg/cryptosuit"
	"github.com/tanhuiya/ci123chain/pkg/mortgage"
	types2 "github.com/tanhuiya/ci123chain/pkg/mortgage/types"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
)

// Mortgage
func SignMortgage(from, to string, amount, gas uint64, uniqueID string, priv []byte) ([]byte, error) {
	tx, err := buildMortgageTx(from, to, amount, gas, uniqueID)
	if err != nil {
		return nil, err
	}
	sid := cryptosuit.NewFabSignIdentity()
	pub, err  := sid.GetPubKey(priv)

	tx.SetPubKey(pub)
	signbyte := tx.GetSignBytes()
	signature, err := sid.Sign(signbyte, priv)
	tx.SetSignature(signature)
	return tx.Bytes(), nil
}

func buildMortgageTx (from, to string, amount, gas uint64, uniqueID string) (*types2.MsgMortgage, error) {
	fromAddr, err := helper.StrToAddress(from)
	if err != nil {
		return nil, err
	}
	toAddr, err := helper.StrToAddress(to)
	if err != nil {
		return nil, err
	}
	nonce, err := transaction.GetNonceByAddress(fromAddr)
	mort := mortgage.NewMortgageMsg(fromAddr, toAddr, gas, nonce, types.Coin(amount), []byte(uniqueID))
	return mort, nil
}

// MortgageDone

func SignMortgageDone(from string, gas uint64, uniqueID string, priv []byte) ([]byte, error) {
	tx, err := buildMortgageDoneTx(from, gas, uniqueID)
	if err != nil {
		return nil, err
	}
	sid := cryptosuit.NewFabSignIdentity()
	pub, err  := sid.GetPubKey(priv)

	tx.SetPubKey(pub)
	signbyte := tx.GetSignBytes()
	signature, err := sid.Sign(signbyte, priv)
	tx.SetSignature(signature)
	return tx.Bytes(), nil
}



func buildMortgageDoneTx (from string, gas uint64, uniqueID string) (*types2.MsgMortgageDone, error) {
	fromAddr, err := helper.StrToAddress(from)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	nonce, err := transaction.GetNonceByAddress(fromAddr)
	mort := mortgage.NewMsgMortgageDone(fromAddr, gas, nonce, []byte(uniqueID))
	return mort, nil
}


// MortgageCancel

func SignMortgageCancel(from string, gas uint64, uniqueID string, priv []byte) ([]byte, error) {
	tx, err := buildMortgageCancelTx(from, gas, uniqueID)
	if err != nil {
		return nil, err
	}
	sid := cryptosuit.NewFabSignIdentity()
	pub, err  := sid.GetPubKey(priv)

	tx.SetPubKey(pub)
	signbyte := tx.GetSignBytes()
	signature, err := sid.Sign(signbyte, priv)
	tx.SetSignature(signature)
	return tx.Bytes(), nil
}

func buildMortgageCancelTx (from string, gas uint64, uniqueID string) (*types2.MsgMortgageCancel, error) {
	fromAddr, err := helper.StrToAddress(from)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	nonce, err := transaction.GetNonceByAddress(fromAddr)
	mort := mortgage.NewMsgMortgageCancel(fromAddr, gas, nonce, []byte(uniqueID))
	return mort, nil
}