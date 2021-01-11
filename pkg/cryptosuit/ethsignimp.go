package cryptosuit

import (
	"bytes"
	"github.com/ethereum/go-ethereum/crypto"
)

type ethsignimp struct {

}

// 签名
func (eth ethsignimp)Sign(msg []byte, identity []byte) ([]byte, error) {
	key := crypto.ToECDSAUnsafe(identity)
	signature, err := crypto.Sign(msg, key)
	return signature, err
}

func (eth ethsignimp) Verifier(msg []byte, signature []byte, _ []byte, address []byte) (bool, error) {
	rawPub, err := crypto.Ecrecover(msg, signature)
	if err != nil {
		return false, err
	}
	pubKey, err := crypto.UnmarshalPubkey(rawPub)
	if err != nil {
		return false, err
	}
	signer := crypto.PubkeyToAddress(*pubKey)
	if !bytes.Equal(signer[:], address) {
		return false, err
	}
	return true, nil
}


func (fab ethsignimp) GetPubKey(privKey []byte) ([]byte, error) {
	return nil, nil
}