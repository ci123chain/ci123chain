package cryptosuite

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding/hex"
	"encoding/pem"
	"golang.org/x/crypto/sha3"
)

const (
	EncryptVersion1 = "V1"
)

type KeySave struct {
	Key        interface{}
	EncryptKey string
}

//密码加盐
func GenFromPassword(password, passwordSalt []byte) []byte {
	sum := sha512.Sum512(append(password, passwordSalt...))
	return sum[:]
}

//使用一个对称秘钥+额外的key（例如用户密码，或是某个环境变量）来加密某个key
func EncryptKey(text []byte, syKey []byte, extraKey []byte) (string, error) {
	realKey := sha3.Sum256(append(syKey, extraKey...))
	var iv = realKey[:aes.BlockSize]
	encrypted := make([]byte, len(text))
	block, err := aes.NewCipher(realKey[:])
	if err != nil {
		return "", err
	}
	encrypter := cipher.NewCFBEncrypter(block, iv)
	encrypter.XORKeyStream(encrypted, text)
	return EncryptVersion1 + hex.EncodeToString(encrypted), nil
}

func DecryptKey(encrypted string, key []byte, extraKey []byte) (res []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	version := encrypted[:2]
	if version == EncryptVersion1 {
		realKey := sha3.Sum256(append(key, extraKey...))
		return decrypt(encrypted[2:], realKey[:])
	} else {
		return decrypt(encrypted, key)
	}
}

func decrypt(encrypted string, key []byte) ([]byte, error) {
	var err error
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	src, err := hex.DecodeString(encrypted)
	if err != nil {
		return nil, err
	}
	var iv = key[:aes.BlockSize]
	decrypted := make([]byte, len(src))
	var block cipher.Block
	block, err = aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	decrypter := cipher.NewCFBDecrypter(block, iv)
	decrypter.XORKeyStream(decrypted, src)
	return decrypted, nil
}

func pemToHex(input []byte) string {
	p, _ := pem.Decode(input)
	return hex.EncodeToString(p.Bytes)
}
