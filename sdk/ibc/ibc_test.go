package ibc

import (
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

func TestIBCMsg(t *testing.T)  {
	// 将pem 格式私钥转化为 十六机制 字符串
	priKey, err := cryptoutil.DecodePriv([]byte(testPrivKey))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)

	signdata, err := SignIBCTransferMsg("0x204bCC42559Faf6DFE1485208F7951aaD800B313",
		"0xD1a14962627fAc768Fe885Eeb9FF072706B54c19", 10, 1, privByte)

	assert.NoError(t, err)
	fmt.Println(hex.EncodeToString(signdata))
}

