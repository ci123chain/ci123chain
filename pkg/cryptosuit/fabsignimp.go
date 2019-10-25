package cryptosuit

import (
	"CI123Chain/pkg/client/helper"
	"bytes"
	"github.com/pkg/errors"
	"github.com/tanhuiya/fabric-crypto/cryptosuite"
	"github.com/tanhuiya/fabric-crypto/cryptoutil"
)

type fabsignimp struct {

}

func (fab fabsignimp) Sign(msg []byte, priv []byte) ([]byte, error) {
	pkey, err := cryptoutil.GetPrivKeyFromKey(priv, cryptosuite.GetDefault())
	signature, err := cryptoutil.SignMsg(msg, pkey, cryptosuite.GetDefault())
	if err != nil {
		return nil, err
	}
	return signature, nil
}


func (fab fabsignimp)Verifier(msg []byte, signature []byte, pub []byte, address []byte) (bool, error) {
	// 公钥验证
	valid, err := cryptoutil.VerifyFromPubString(msg, pub, signature)

	if !valid || err != nil {
		return false, errors.Wrap(err, "valid signature failed")
	}
	pubaddress, err := cryptoutil.EcdsaPubToAddress(pub)
	if err != nil {
		return false, err
	}
	// 匹配地址
	addrByte, err := helper.StrToAddress(pubaddress)
	if err != nil {
		return false, err
	}
	if !bytes.Equal(addrByte[:], address) {
		return false, errors.Wrap(err, "address not match")
	}

	return true, nil
}

func (fab fabsignimp) GetPubKey(privKey []byte) ([]byte, error) {
	pkey, err := cryptoutil.GetPrivKeyFromKey(privKey, cryptosuite.GetDefault())
	if err != nil {
		return nil, err
	}
	pubKey, err := pkey.PublicKey()
	if err != nil {
		return nil, err
	}
	bytes, err := pubKey.Bytes()
	if err != nil {
		return nil, err
	}
	return bytes, nil
}