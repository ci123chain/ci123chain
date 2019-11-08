package sdk

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tanhuiya/fabric-crypto/cryptoutil"
	"testing"
)

var UniqueID = "uniqueID#########12312345###1234"

func TestMortgage(t *testing.T)  {
	// 将pem 格式私钥转化为 十六机制 字符串
	priKey, err := cryptoutil.DecodePriv([]byte(testPrivKey))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)

	signdata, err := SignMortgage("0x204bCC42559Faf6DFE1485208F7951aaD800B313",
		"0xD1a14962627fAc768Fe885Eeb9FF072706B54c19", 10, 1, UniqueID, privByte)

	assert.NoError(t, err)
	fmt.Println(hex.EncodeToString(signdata))
}


func TestMortgageCancel(t *testing.T)  {
	// 将pem 格式私钥转化为 十六机制 字符串
	priKey, err := cryptoutil.DecodePriv([]byte(testPrivKey))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)

	signdata, err := SignMortgageCancel("0x204bCC42559Faf6DFE1485208F7951aaD800B313", 1, UniqueID, privByte)
	assert.NoError(t, err)
	fmt.Println(hex.EncodeToString(signdata))
}


func TestMortgageDone(t *testing.T)  {
	// 将pem 格式私钥转化为 十六机制 字符串
	priKey, err := cryptoutil.DecodePriv([]byte(testPrivKey))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)

	signdata, err := SignMortgageDone("0x204bCC42559Faf6DFE1485208F7951aaD800B313", 1, UniqueID, privByte)
	assert.NoError(t, err)
	fmt.Println(hex.EncodeToString(signdata))
}


// 使用错误的 from和 priv 签名，返回广播消息会失败
func TestMortgageDoneFake(t *testing.T) {
	// 将pem 格式私钥转化为 十六机制 字符串
	priKey, err := cryptoutil.DecodePriv([]byte(testPriv2))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)

	signdata, err := SignMortgageDone("0x55B9dB2c4d8336FA4a2005EFAC447c8f4568c7c7", 1, UniqueID, privByte)
	assert.NoError(t, err)
	fmt.Println(hex.EncodeToString(signdata))
}