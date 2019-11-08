package sdk

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tanhuiya/fabric-crypto/cryptoutil"
	"testing"
)

func TestMortgage(t *testing.T)  {
	// 将pem 格式私钥转化为 十六机制 字符串
	priKey, err := cryptoutil.DecodePriv([]byte(testPrivKey))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)

	signdata, err := SignMortgage("0x204bCC42559Faf6DFE1485208F7951aaD800B313",
		"0xD1a14962627fAc768Fe885Eeb9FF072706B54c19", 10, 1, "uniqueID#########1231234###", privByte)

	assert.NoError(t, err)
	fmt.Println(hex.EncodeToString(signdata))
}