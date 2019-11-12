package ibc

import (
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tanhuiya/fabric-crypto/cryptoutil"
	"testing"
)

const UniqueID = "698D871C159F25B51C41AC09C5552FE3"
const ObserverID = "1234567812345678"

func TestApplyIBCMsg(t *testing.T)  {
	// 将pem 格式私钥转化为 十六机制 字符串
	priKey, err := cryptoutil.DecodePriv([]byte(testPrivKey))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)

	uid := []byte(UniqueID)
	signdata, err := SignApplyIBCMsg("0x204bCC42559Faf6DFE1485208F7951aaD800B313", uid, []byte(ObserverID), 1, privByte)

	assert.NoError(t, err)
	fmt.Println(hex.EncodeToString(signdata))
}


