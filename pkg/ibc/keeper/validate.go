package keeper

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tanhuiya/ci123chain/pkg/cryptosuit"
	"github.com/tanhuiya/ci123chain/pkg/ibc/types"
	"github.com/tanhuiya/fabric-crypto/cryptoutil"
	"github.com/tanhuiya/ci123chain/pkg/transaction"
)

// 验证 apply 消息
func ValidateRawIBCMessage(tx types.IBCMsgBankSend) (*types.IBCInfo, sdk.Error) {
	var signObj types.ApplyReceipt

	// 反序列化
	err := json.Unmarshal(tx.RawMessage, &signObj)
	if err != nil {
		return nil, types.ErrFailedUnmarshal(types.DefaultCodespace, err.Error())
	}
	if len(signObj.Signature) < 1 || len(signObj.IBCMsgBytes) < 1 {
		return nil, types.ErrBadBankSignature(types.DefaultCodespace, errors.New("signature or ibcMsgBytes len less than 1"))
	}

	sid := cryptosuit.NewFabSignIdentity()
	privKey, _ := cryptoutil.DecodePriv([]byte(Priv))

	pubKey := privKey.Public().(*ecdsa.PublicKey)
	pubketBz := cryptoutil.MarshalPubkey(pubKey)
	valid, err := sid.Verifier(signObj.GetSignBytes(), signObj.Signature, pubketBz, nil)
	if !valid  {
		return nil, types.ErrBadBankSignature(types.DefaultCodespace, err)
	}

	var ibcMsg types.IBCInfo
	err = json.Unmarshal(signObj.IBCMsgBytes, &ibcMsg)
	//fmt.Println(string(ibcMsg.UniqueID))
	if err != nil {
		return nil, types.ErrFailedUnmarshal(types.DefaultCodespace, err.Error())
	}
	return &ibcMsg, nil
}


// 验证 回执 消息
func ValidateRawReceiptMessage(tx types.IBCReceiveReceiptMsg) (*types.BankReceipt, sdk.Error) {
	var receiveObj types.BankReceipt
	// 反序列化
	err := json.Unmarshal(tx.RawMessage, &receiveObj)
	if err != nil {
		return nil, types.ErrFailedUnmarshal(types.DefaultCodespace, err.Error())
	}

	sid := cryptosuit.NewFabSignIdentity()
	pubBz, err := getPublicKey()
	if err != nil {
		return nil, transaction.ErrBadPubkey(types.DefaultCodespace, err)
	}

	valid, err := sid.Verifier(receiveObj.GetSignBytes(), receiveObj.Signature, pubBz, nil)
	if !valid || err != nil {
		return nil, types.ErrBadReceiptSignature(types.DefaultCodespace, err)
	}
	return &receiveObj, nil
}
