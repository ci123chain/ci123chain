package keeper

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"github.com/tanhuiya/ci123chain/pkg/cryptosuit"
	"github.com/tanhuiya/ci123chain/pkg/ibc/types"
	"github.com/tanhuiya/fabric-crypto/cryptoutil"
)

// 验证 apply 消息
func ValidateRawIBCMessage(tx types.IBCMsgBankSend) (*types.IBCInfo, error) {
	var signObj types.ApplyReceipt
	// 反序列化
	err := json.Unmarshal(tx.RawMessage, &signObj)
	if err != nil {
		return nil, err
	}

	sid := cryptosuit.NewFabSignIdentity()
	privKey, _ := cryptoutil.DecodePriv([]byte(Priv))

	pubKey := privKey.Public().(*ecdsa.PublicKey)
	pubketBz := cryptoutil.MarshalPubkey(pubKey)
	valid, err := sid.Verifier(signObj.GetSignBytes(), signObj.Signature, pubketBz, nil)
	if !valid || err != nil {
		return nil, errors.New("pkg invalid signature; " + err.Error())
	}

	var ibcMsg types.IBCInfo
	err = json.Unmarshal(signObj.IBCMsgBytes, &ibcMsg)
	if err != nil {
		return nil, err
	}
	return &ibcMsg, nil
}


// 验证 回执 消息
func ValidateRawReceiptMessage(tx types.IBCReceiveReceiptMsg) (*types.BankReceipt, error) {
	var receiveObj types.BankReceipt
	// 反序列化
	err := json.Unmarshal(tx.RawMessage, &receiveObj)
	if err != nil {
		return nil, err
	}

	sid := cryptosuit.NewFabSignIdentity()
	pubBz, err := getPublicKey()
	if err != nil {
		return nil, err
	}

	valid, err := sid.Verifier(receiveObj.GetSignBytes(), receiveObj.Signature, pubBz, nil)
	if !valid || err != nil {
		return nil, errors.New("pkg invalid signature; " + err.Error())
	}
	return &receiveObj, nil
}
