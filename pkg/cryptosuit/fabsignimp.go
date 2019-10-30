package cryptosuit

import (
	"gitlab.oneitfarm.com/blockchain/ci123chain/pkg/client/helper"
	"bytes"
	"crypto/ecdsa"
	"github.com/pkg/errors"
	"github.com/tanhuiya/fabric-crypto/cryptoutil"
)

type fabsignimp struct {

}

func (fab fabsignimp) Sign(msg []byte, priv []byte) ([]byte, error) {
	privKey, err := cryptoutil.UnMarshalPrivateKey(priv)
	if err != nil {
		return nil, errors.Wrap(err, "Decode privatekey error")
	}
	signature, err := cryptoutil.SignMsg2(msg, privKey)
	if err != nil {
		return nil, errors.Wrap(err, "sign action error")
	}
	return signature, nil
}


func (fab fabsignimp)Verifier(msg []byte, signature []byte, pub []byte, address []byte) (bool, error) {
	pubKey, err := cryptoutil.UnMarshalPubKey(pub)
	if err != nil {
		return false, err
	}
	// 公钥验证
	valid, err := cryptoutil.Verifier2(msg, signature, pubKey)

	if !valid || err != nil {
		return false, errors.Wrap(err, "valid signature failed")
	}

	if len(address) > 0{
		pubaddress, err := cryptoutil.PublicKeyToAddress(pubKey)
		if err != nil {
			return false, err
		}
		// 匹配地址
		addrByte, err := helper.StrToAddress(pubaddress)
		if err != nil {
			return false, err
		}
		if !bytes.Equal(addrByte[:], address) {
			return false, errors.Wrap(err, "address and privateKey not match")
		}
	}

	return true, nil
}

func (fab fabsignimp) GetPubKey(privKey []byte) ([]byte, error) {
	priv, err := cryptoutil.UnMarshalPrivateKey(privKey)
	if err != nil {
		return nil, err
	}
	pub := priv.Public().(*ecdsa.PublicKey)
	pubByte := cryptoutil.MarshalPubkey(pub)
	return pubByte, nil
}