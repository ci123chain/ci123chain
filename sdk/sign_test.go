package sdk

import (
	"gitlab.oneitfarm.com/blockchain/ci123chain/pkg/client/helper"
	"gitlab.oneitfarm.com/blockchain/ci123chain/pkg/cryptosuit"
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tanhuiya/fabric-crypto/cryptoutil"
	"testing"
)

var testPrivKey = `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgp4qKKB0WCEfx7XiB
5Ul+GpjM1P5rqc6RhjD5OkTgl5OhRANCAATyFT0voXX7cA4PPtNstWleaTpwjvbS
J3+tMGTG67f+TdCfDxWYMpQYxLlE8VkbEzKWDwCYvDZRMKCQfv2ErNvb
-----END PRIVATE KEY-----`

var testCert = `-----BEGIN CERTIFICATE-----
MIICGTCCAcCgAwIBAgIRALR/1GXtEud5GQL2CZykkOkwCgYIKoZIzj0EAwIwczEL
MAkGA1UEBhMCVVMxEzARBgNVBAgTCkNhbGlmb3JuaWExFjAUBgNVBAcTDVNhbiBG
cmFuY2lzY28xGTAXBgNVBAoTEG9yZzEuZXhhbXBsZS5jb20xHDAaBgNVBAMTE2Nh
Lm9yZzEuZXhhbXBsZS5jb20wHhcNMTcwNzI4MTQyNzIwWhcNMjcwNzI2MTQyNzIw
WjBbMQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMN
U2FuIEZyYW5jaXNjbzEfMB0GA1UEAwwWVXNlcjFAb3JnMS5leGFtcGxlLmNvbTBZ
MBMGByqGSM49AgEGCCqGSM49AwEHA0IABPIVPS+hdftwDg8+02y1aV5pOnCO9tIn
f60wZMbrt/5N0J8PFZgylBjEuUTxWRsTMpYPAJi8NlEwoJB+/YSs29ujTTBLMA4G
A1UdDwEB/wQEAwIHgDAMBgNVHRMBAf8EAjAAMCsGA1UdIwQkMCKAIIeR0TY+iVFf
mvoEKwaToscEu43ZXSj5fTVJornjxDUtMAoGCCqGSM49BAMCA0cAMEQCID+dZ7H5
AiaiI2BjxnL3/TetJ8iFJYZyWvK//an13WV/AiARBJd/pI5A7KZgQxJhXmmR8bie
XdsmTcdRvJ3TS/6HCA==
-----END CERTIFICATE-----`

func TestSignTx(t *testing.T) {
	// 将pem 格式私钥转化为 十六机制 字符串
	priKey, err := cryptoutil.DecodePriv([]byte(testPrivKey))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)

	signedData, err := SignTx("0x204bCC42559Faf6DFE1485208F7951aaD800B313",
		"0x204bCC42559Faf6DFE1485208F7951aaD800B313",
		1,
		1,
		privByte,
		true)

	assert.NoError(t, err)
	fmt.Println(hex.EncodeToString(signedData))
}


func TestVerifier(t *testing.T)  {
	tx, _ := buildTransferTx("0x204bCC42559Faf6DFE1485208F7951aaD800B313",
		"0x204bCC42559Faf6DFE1485208F7951aaD800B313",
		1,
		1,
		true,
		)
	sid := cryptosuit.NewFabSignIdentity()

	priKey, err := cryptoutil.DecodePriv([]byte(testPrivKey))
	assert.NoError(t, err)

	signature, err := sid.Sign(tx.GetSignBytes(), cryptoutil.MarshalPrivateKey(priKey))
	assert.NoError(t, err)
	addrbyte, _ := helper.StrToAddress("0x204bCC42559Faf6DFE1485208F7951aaD800B313")

	pubkey, _ := cryptoutil.DecodePub([]byte(testCert))
	assert.Equal(t, "04f2153d2fa175fb700e0f3ed36cb5695e693a708ef6d2277fad3064c6ebb7fe4dd09f0f1598329418c4b944f1591b1332960f0098bc365130a0907efd84acdbdb", hex.EncodeToString(cryptoutil.MarshalPubkey(pubkey)))
	signature, _ = hex.DecodeString(hex.EncodeToString(signature))

	valid, err := Verifier(tx.GetSignBytes(), signature , cryptoutil.MarshalPubkey(pubkey), addrbyte[:])
	assert.NoError(t, err)
	assert.Equal(t, true, valid)
}
