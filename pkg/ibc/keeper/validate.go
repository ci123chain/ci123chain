package keeper

import (
	"crypto/ecdsa"
	"encoding/json"
	sdkerrors "github.com/ci123chain/ci123chain/pkg/abci/types/errors"
	"github.com/ci123chain/ci123chain/pkg/cryptosuit"
	"github.com/ci123chain/ci123chain/pkg/ibc/types"
	"github.com/tanhuiya/fabric-crypto/cryptoutil"
)

// 验证 apply 消息
func ValidateRawIBCMessage(tx types.IBCMsgBankSend) (*types.IBCInfo, error) {
	var signObj types.ApplyReceipt

	// 反序列化
	err := json.Unmarshal(tx.RawMessage, &signObj)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	if len(signObj.Signature) < 1 || len(signObj.IBCMsgBytes) < 1 {
		return nil, sdkerrors.Wrap(sdkerrors.ErrNoSignatures, "signature or ibcMsgBytes len less than 1")
	}

	sid := cryptosuit.NewFabSignIdentity()
	privKey, _ := cryptoutil.DecodePriv([]byte(Priv))

	pubKey := privKey.Public().(*ecdsa.PublicKey)
	pubketBz := cryptoutil.MarshalPubkey(pubKey)
	valid, err := sid.Verifier(signObj.GetSignBytes(), signObj.Signature, pubketBz, nil)
	if !valid  {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInternal, "signature verify faield")
	}

	var ibcMsg types.IBCInfo
	err = json.Unmarshal(signObj.IBCMsgBytes, &ibcMsg)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}
	return &ibcMsg, nil
}


// 验证 回执 消息
func ValidateRawReceiptMessage(tx types.IBCReceiveReceiptMsg) (*types.BankReceipt, error) {
	var receiveObj types.BankReceipt
	// 反序列化
	err := json.Unmarshal(tx.RawMessage, &receiveObj)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	sid := cryptosuit.NewFabSignIdentity()
	pubBz, err := getPublicKey()
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrParams, "get pubkey failed")
	}

	valid, err := sid.Verifier(receiveObj.GetSignBytes(), receiveObj.Signature, pubBz, nil)
	if !valid || err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrParams, "signature verify failed")
	}
	return &receiveObj, nil
}
