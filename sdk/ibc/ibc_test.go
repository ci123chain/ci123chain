package ibc

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tanhuiya/fabric-crypto/cryptoutil"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

var testPrivKey = `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgp4qKKB0WCEfx7XiB
5Ul+GpjM1P5rqc6RhjD5OkTgl5OhRANCAATyFT0voXX7cA4PPtNstWleaTpwjvbS
J3+tMGTG67f+TdCfDxWYMpQYxLlE8VkbEzKWDwCYvDZRMKCQfv2ErNvb
-----END PRIVATE KEY-----`

// 生成 跨链 交易
func TestIBCMsg(t *testing.T)  {
	// 将pem 格式私钥转化为 十六机制 字符串
	priKey, err := cryptoutil.DecodePriv([]byte(testPrivKey))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)

	signdata, err := SignIBCTransferMsg("0x204bCC42559Faf6DFE1485208F7951aaD800B313",
		"0xD1a14962627fAc768Fe885Eeb9FF072706B54c19", 10, 1, privByte)

	assert.NoError(t, err)
	httpPost(hex.EncodeToString(signdata))
}


// 生成 apply 签名交易
const UniqueID = "8A37BA6B8013ABE59F278EFE33D5B188"
const ObserverID = "1234567812345678"
func TestApplyIBCMsg(t *testing.T)  {
	// 将pem 格式私钥转化为 十六机制 字符串
	priKey, err := cryptoutil.DecodePriv([]byte(testPrivKey))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)

	uid := []byte(UniqueID)
	signdata, err := SignApplyIBCMsg("0x204bCC42559Faf6DFE1485208F7951aaD800B313", uid, []byte(ObserverID), 1, privByte)

	assert.NoError(t, err)
	//fmt.Println(hex.EncodeToString(signdata))
	httpPost(hex.EncodeToString(signdata))
}


// 生成 bank 到 account 扣钱交易
const pkg =
`{"signature":"MEQCIF8Lp+L1p3p0WV5tf68QdN7Mf04mLMLIx9j0FSOefhE1AiBEcPxKCIu7y0B09R5IUHmy5tXeiX5rhPWTXXFyd2Tljw==","ibc_msg_bytes":"eyJ1bmlxdWVfaWQiOiI4QTM3QkE2QjgwMTNBQkU1OUYyNzhFRkUzM0Q1QjE4OCIsIm9ic2VydmVyX2lkIjoiMTIzNDU2NzgxMjM0NTY3OCIsImJhbmtfYWRkcmVzcyI6IjB4NTA1QTc0Njc1ZGM5QzcxZUYzQ0I1REYzMDkyNTY5NTI5MTdFODAxZSIsImFwcGx5X3RpbWUiOiIyMDE5LTExLTEzVDE2OjU4OjMyLjAwNzkyKzA4OjAwIiwic3RhdGUiOiJwcm9jZXNzaW5nIiwiZnJvbV9hZGRyZXNzIjoiMHgyMDRiQ0M0MjU1OUZhZjZERkUxNDg1MjA4Rjc5NTFhYUQ4MDBCMzEzIiwidG9fYWRkcmVzcyI6IjB4RDFhMTQ5NjI2MjdmQWM3NjhGZTg4NUVlYjlGRjA3MjcwNkI1NGMxOSIsImFtb3VudCI6MTB9"}`

func TestBankSendMsg(t *testing.T)  {
	// 将pem 格式私钥转化为 十六机制 字符串
	priKey, err := cryptoutil.DecodePriv([]byte(testPrivKey))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)

	signdata, err := SignIBCBankSendMsg("0x204bCC42559Faf6DFE1485208F7951aaD800B313", []byte(pkg), 1, privByte)

	assert.NoError(t, err)
	httpPost(hex.EncodeToString(signdata))
}

const pkgReceipt =
`{"unique_id":"8A37BA6B8013ABE59F278EFE33D5B188","observer_id":"1234567812345678","signature":"MEQCIHNGnWa/xk4n+WOERiXphkytHN+iOIfQiwJTozixLBXnAiBPyZbtpUHwbp4ATUg90Tmye/iNZ9sc7Q3jO0RsZ1Jfag=="}`

func TestReceiptMsg(t *testing.T)  {
	// 将pem 格式私钥转化为 十六机制 字符串
	priKey, err := cryptoutil.DecodePriv([]byte(testPrivKey))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)

	signdata, err := SignIBCReceiptMsg("0x204bCC42559Faf6DFE1485208F7951aaD800B313", []byte(pkgReceipt), 1, privByte)
	assert.NoError(t, err)
	httpPost(hex.EncodeToString(signdata))
}

const FromAddr  = "0x204bCC42559Faf6DFE1485208F7951aaD800B313"
const ToAddr  = "0xD1a14962627fAc768Fe885Eeb9FF072706B54c19"
func TestAll(t *testing.T)  {
	// 将pem 格式私钥转化为 十六机制 字符串
	priKey, err := cryptoutil.DecodePriv([]byte(testPrivKey))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)

	fmt.Println("---发送跨链消息")
	signdata, err := SignIBCTransferMsg(FromAddr,
		ToAddr, 10, 1, privByte)
	registRet := httpPost(hex.EncodeToString(signdata))
	fmt.Println("发送跨链消息完成：UniqueID = " + registRet.Data)
	fmt.Println()

	fmt.Println("---申请处理该跨链消息")
	uid := []byte(registRet.Data)
	signdata, err = SignApplyIBCMsg(FromAddr, uid, []byte(ObserverID), 1, privByte)
	applyRet := httpPost(hex.EncodeToString(signdata))
	assert.True(t, len(applyRet.RawLog) < 1)
	fmt.Println("申请处理该跨链消息结束")
	fmt.Println()

	fmt.Println("---第二个申请处理该跨链消息")
	ObserverID2 := "12313213213213124321"
	signdata, err = SignApplyIBCMsg(FromAddr, uid, []byte(ObserverID2), 1, privByte)
	applyRetErr := httpPost(hex.EncodeToString(signdata))
	assert.True(t, len(applyRetErr.RawLog) > 0)
	fmt.Println("申请处理该跨链消息结束失败")
	fmt.Println()

	// bank转账，该交易参数应该是observer 从 fabric 获得，此处是模拟
	fmt.Println("---向对方转账")
	pkg := applyRet.Data
	signdata, err = SignIBCBankSendMsg(FromAddr, []byte(pkg), 1, privByte)
	receiptRet := httpPost(hex.EncodeToString(signdata))
	assert.True(t, len(receiptRet.RawLog) < 1)
	fmt.Println("向对方转账成功")
	fmt.Println()

	//发送回执
	fmt.Println("---发送回执")
	receivepkg := receiptRet.Data
	signdata, err = SignIBCReceiptMsg(FromAddr, []byte(receivepkg), 1, privByte)
	ret := httpPost(hex.EncodeToString(signdata))
	assert.True(t, len(ret.RawLog) < 1)
	fmt.Println("发送回执成功")
}



func httpPost(param string) retData {
	resp, err := http.PostForm("http://localhost:1317/tx/broadcast",
		url.Values{"data": {param}})
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	var ret retData
	err = json.Unmarshal(body, &ret)
	if err != nil {
		fmt.Println(string(body))
	}
	if len(ret.Data) > 0 {
		fmt.Println(ret.Data)
	}
	if len(ret.RawLog) > 0 {
		fmt.Println(ret.RawLog)
	}
	return ret
}

type retData struct {
	Data string `json:"data"`
	RawLog  string `json:"raw_log"`
}