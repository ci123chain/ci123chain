package cryptosuite

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
	"hash"
)

type EccK1 struct {
	curve        elliptic.Curve
	sigAlgorithm x509.SignatureAlgorithm
	hashFunction func() hash.Hash
}

func NewEccK1() *EccK1 {
	return &EccK1{
		curve:        crypto.S256(),
		sigAlgorithm: x509.ECDSAWithSHA256,
		hashFunction: sha3.New256,
	}
}

func (c EccK1) GenerateKey() (priKeyByte []byte, pubKeyByte []byte, err error) {
	key, err := ecdsa.GenerateKey(c.curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	priKeyByte = crypto.FromECDSA(key)
	pubKeyByte = crypto.FromECDSAPub(&key.PublicKey)
	return priKeyByte, pubKeyByte, nil
}

func (c EccK1) Sign(priKey []byte, digest []byte) (signature []byte, err error) {
	priKeyIns, err := crypto.ToECDSA(priKey)
	if err != nil {
		return nil, err
	}
	return crypto.Sign(digest[:], priKeyIns)
}

func (c EccK1) Verify(from []byte, digest []byte, sig []byte) (bool, error) {
	if sig[64] == 27 || sig[64] == 28 {
		sig[64] -= 27
	}

	key, err := crypto.SigToPub(digest, sig)
	if err != nil {
		return false, err
	}
	addr := crypto.PubkeyToAddress(*key)
	if from != nil {
		if !bytes.Equal(from, addr.Bytes()) {
			return false, errors.New("public key not matching signature")
		}
	}

	parsedKey := crypto.FromECDSAPub(key)
	flag := crypto.VerifySignature(parsedKey, digest[:], sig[:len(sig)-1])
	return flag, nil
}

func Hash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

func (c EccK1) GetPriKeyPem(key interface{}) ([]byte, error) {
	if eccKey, ok := key.(*ecdsa.PrivateKey); !ok {
		return nil, ErrInvalidKeyType
	} else {
		keyBytes := crypto.FromECDSA(eccKey)
		block := pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: keyBytes,
		}
		pemBytes := pem.EncodeToMemory(&block)
		return pemBytes, nil
	}
}

func (c EccK1) GetPubKeyPem(key interface{}) ([]byte, error) {
	if eccKey, ok := key.(*ecdsa.PublicKey); !ok {
		return nil, ErrInvalidKeyType
	} else {
		keyBytes := crypto.FromECDSAPub(eccKey)
		block := pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: keyBytes,
		}
		pemBytes := pem.EncodeToMemory(&block)
		return pemBytes, nil
	}
}

func (c EccK1) ParsePemPriKey(key []byte) (interface{}, error) {
	//p, _ := pem.Decode(key)
	//if p == nil || len(p.Bytes) == 0 {
	//	return nil, errors.New("key parse error")
	//}
	//panic("implement me")
	p, _ := pem.Decode(key)
	if p == nil || len(p.Bytes) == 0 {
		return nil, errors.New("key parse error")
	}
	ecdsaPrivateKey, err := crypto.ToECDSA(p.Bytes)
	if err != nil {
		return nil, err
	}
	return ecdsaPrivateKey, nil
}

func (c EccK1) ParsePemPubKey(key []byte) (interface{}, error) {
	return crypto.UnmarshalPubkey(key)
}
