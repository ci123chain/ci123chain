package transaction

import (
	"github.com/tanhuiya/ci123chain/pkg/client/helper"
	"github.com/tanhuiya/ci123chain/pkg/cryptosuit"
	"fmt"
	"github.com/magiconair/properties/assert"
	"github.com/tanhuiya/fabric-crypto/cryptosuite"
	"github.com/tanhuiya/fabric-crypto/cryptoutil"
	"testing"
)

var testPrivKey = `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgp4qKKB0WCEfx7XiB
5Ul+GpjM1P5rqc6RhjD5OkTgl5OhRANCAATyFT0voXX7cA4PPtNstWleaTpwjvbS
J3+tMGTG67f+TdCfDxWYMpQYxLlE8VkbEzKWDwCYvDZRMKCQfv2ErNvb
-----END PRIVATE KEY-----`

func TestFabSign(t *testing.T)  {
	sid := cryptosuit.GetSignIdentity(cryptosuit.FabSignType)

	privKey, _  := cryptoutil.GetPrivKeyFromKey([]byte(testPrivKey), cryptosuite.GetDefault())
	pubKey, _ := privKey.PublicKey()
	address, _ := cryptoutil.PublicKeyToAddress(pubKey)

	froms, _ := helper.ParseAddrs(address)
	nonce, err := GetNonceByAddress(froms[0])
	tx := NewTransferTx(froms[0], froms[0], 1, nonce, 10, true)

	pubByte, err := sid.GetPubKey([]byte(testPrivKey))

	if err != nil {
		panic(err)
	}
	tx.SetPubKey(pubByte)
	sig, err := sid.Sign(tx.GetSignBytes(), []byte(testPrivKey))
	if err != nil {
		panic(err)
	}
	fmt.Println(sig)

	valid, err := sid.Verifier(tx.GetSignBytes(), sig, pubByte, []byte(address))
	assert.Equal(t, true, valid)
}