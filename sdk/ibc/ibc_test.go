package ibc

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tanhuiya/fabric-crypto/cryptoutil"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"testing"
)

var testPrivKey = `-----BEGIN PRIVATE KEY-----
MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgp4qKKB0WCEfx7XiB
5Ul+GpjM1P5rqc6RhjD5OkTgl5OhRANCAATyFT0voXX7cA4PPtNstWleaTpwjvbS
J3+tMGTG67f+TdCfDxWYMpQYxLlE8VkbEzKWDwCYvDZRMKCQfv2ErNvb
-----END PRIVATE KEY-----`

var ip = "127.0.0.1"
var port = "1317"

// 生成 跨链 交易
func TestIBCMsg(t *testing.T)  {
	// 将pem 格式私钥转化为 十六机制 字符串

	//获取nonce，地址， 端口
	//nonce := httpQuery(ip, port, FromAddr)

	priKey, err := cryptoutil.DecodePriv([]byte(testPrivKey))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)


	signdata, err := SignIBCTransferMsg("0x204bCC42559Faf6DFE1485208F7951aaD800B313",
		"0xD1a14962627fAc768Fe885Eeb9FF072706B54c19", 6, privByte)
	fmt.Println(hex.EncodeToString(signdata))

	//assert.NoError(t, err)
	httpPost(hex.EncodeToString(signdata))
}


// 生成 apply 签名交易
const UniqueID = "61849E3829B6B42616BC2736FA44CBBE"
const ObserverID = "1234567812345679"
func TestApplyIBCMsg(t *testing.T)  {
	//nonce := httpQuery(ip, port, FromAddr)

	// 将pem 格式私钥转化为 十六机制 字符串
	priKey, err := cryptoutil.DecodePriv([]byte(testPrivKey))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)

	uid := []byte(UniqueID)
	signdata, err := SignApplyIBCMsg("0x204bCC42559Faf6DFE1485208F7951aaD800B313", uid, []byte(ObserverID), privByte)

	assert.NoError(t, err)
	//fmt.Println(hex.EncodeToString(signdata))
	httpPost(hex.EncodeToString(signdata))
}


// 生成 bank 到 account 扣钱交易
const pkg =
`{"signature":"MEQCIF8Lp+L1p3p0WV5tf68QdN7Mf04mLMLIx9j0FSOefhE1AiBEcPxKCIu7y0B09R5IUHmy5tXeiX5rhPWTXXFyd2Tljw==","ibc_msg_bytes":"eyJ1bmlxdWVfaWQiOiI4QTM3QkE2QjgwMTNBQkU1OUYyNzhFRkUzM0Q1QjE4OCIsIm9ic2VydmVyX2lkIjoiMTIzNDU2NzgxMjM0NTY3OCIsImJhbmtfYWRkcmVzcyI6IjB4NTA1QTc0Njc1ZGM5QzcxZUYzQ0I1REYzMDkyNTY5NTI5MTdFODAxZSIsImFwcGx5X3RpbWUiOiIyMDE5LTExLTEzVDE2OjU4OjMyLjAwNzkyKzA4OjAwIiwic3RhdGUiOiJwcm9jZXNzaW5nIiwiZnJvbV9hZGRyZXNzIjoiMHgyMDRiQ0M0MjU1OUZhZjZERkUxNDg1MjA4Rjc5NTFhYUQ4MDBCMzEzIiwidG9fYWRkcmVzcyI6IjB4RDFhMTQ5NjI2MjdmQWM3NjhGZTg4NUVlYjlGRjA3MjcwNkI1NGMxOSIsImFtb3VudCI6MTB9"}`

func TestBankSendMsg(t *testing.T)  {
	//nonce := httpQuery(ip, port, FromAddr)

	// 将pem 格式私钥转化为 十六机制 字符串
	priKey, err := cryptoutil.DecodePriv([]byte(testPrivKey))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)
	pub := priKey.Public().(*ecdsa.PublicKey)
	addr, _  := cryptoutil.PublicKeyToAddress(pub)

	signdata, err := SignIBCBankSendMsg(addr, []byte(pkg),  privByte)


	assert.NoError(t, err)
	httpPost(hex.EncodeToString(signdata))
}

const pkgReceipt =
`{"unique_id":"8A37BA6B8013ABE59F278EFE33D5B188","observer_id":"1234567812345678","signature":"MEQCIHNGnWa/xk4n+WOERiXphkytHN+iOIfQiwJTozixLBXnAiBPyZbtpUHwbp4ATUg90Tmye/iNZ9sc7Q3jO0RsZ1Jfag=="}`

func TestReceiptMsg(t *testing.T)  {
	//nonce := httpQuery(ip, port, FromAddr)

	// 将pem 格式私钥转化为 十六机制 字符串
	priKey, err := cryptoutil.DecodePriv([]byte(testPrivKey))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)

	signdata, err := SignIBCReceiptMsg("0x204bCC42559Faf6DFE1485208F7951aaD800B313", []byte(pkgReceipt),  privByte)
	assert.NoError(t, err)
	httpPost(hex.EncodeToString(signdata))
}

const FromAddr  = "0x204bCC42559Faf6DFE1485208F7951aaD800B313"
const ToAddr  = "0xD1a14962627fAc768Fe885Eeb9FF072706B54c19"


type ciRes struct{
	Ret 	uint32 	`json:"ret"`
	Data 	string	`json:"data"`
	Message	string	`json:"message"`
}

func TestAll(t *testing.T)  {

	nonce := httpQuery(ip, port, FromAddr)
	// 将pem 格式私钥转化为 十六机制 字符串
	priKey, err := cryptoutil.DecodePriv([]byte(testPrivKey))
	assert.NoError(t, err)
	privByte := cryptoutil.MarshalPrivateKey(priKey)

	fmt.Println("---发送跨链消息")
	nonce = httpQuery(ip, port, FromAddr)
	signdata, err := SignIBCTransferMsg(FromAddr,
		ToAddr, 20, privByte)
	registRet := httpPost(hex.EncodeToString(signdata))
	fmt.Println(nonce)

	fmt.Println("发送跨链消息完成：UniqueID = " + registRet.Data)
	fmt.Println()

	fmt.Println("---申请处理该跨链消息")
	nonce = httpQuery(ip, port, FromAddr)
	uid := []byte(registRet.Data)
	signdata, err = SignApplyIBCMsg(FromAddr, uid, []byte(ObserverID), privByte)
	applyRet := httpPost(hex.EncodeToString(signdata))
	assert.True(t, len(applyRet.RawLog) < 1)
	fmt.Println("申请处理该跨链消息结束")
	fmt.Println()


	fmt.Println("---第二个申请处理该跨链消息")
	nonce = httpQuery(ip, port, FromAddr)
	ObserverID2 := "12313213213213124321"
	signdata, err = SignApplyIBCMsg(FromAddr, uid, []byte(ObserverID2), privByte)

	applyRetErr := httpPost(hex.EncodeToString(signdata))
	assert.True(t, len(applyRetErr.RawLog) > 0)
	fmt.Println("申请处理该跨链消息结束失败")
	fmt.Println()


	// bank转账，该交易参数应该是observer 从 fabric 获得，此处是模拟
	fmt.Println("---向对方转账")
	nonce = httpQuery(ip, port, FromAddr)
	pkg := applyRet.Data
	signdata, err = SignIBCBankSendMsg(FromAddr, []byte(pkg), privByte)
	fmt.Println(hex.EncodeToString(signdata))
	receiptRet := httpPost(hex.EncodeToString(signdata))
	assert.True(t, len(receiptRet.RawLog) < 1)
	fmt.Println("向对方转账成功")
	fmt.Println()

	//发送回执
	fmt.Println("---发送回执")
	nonce = httpQuery(ip, port, FromAddr)
	receivepkg := receiptRet.Data
	signdata, err = SignIBCReceiptMsg(FromAddr, []byte(receivepkg), privByte)
	ret := httpPost(hex.EncodeToString(signdata))
	assert.True(t, len(ret.RawLog) < 1)
	fmt.Println("发送回执成功")
}

func httpQuery(ip, port, param string) uint64 {
	url := "http://" + ip + ":" + port + "/ibctx/nonce/" + param
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var ret string
	err = json.Unmarshal(body, &ret)
	reps, err := strconv.ParseInt(ret, 10, 64)
	nonce := uint64(reps)

	return nonce
}
func httpPost(param string) retData {
	resp, err := http.PostForm("http://127.0.0.1:1317/tx/broadcast",
		url.Values{"data": {param}})
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(body)
	var ret retData
	var ty ciRes
	err = json.Unmarshal(body, &ty)
	fmt.Println(ty.Data)
	d := []byte(ty.Data)
	err = json.Unmarshal(d, &ret)

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

type BalanceData struct {
	Balance uint64 `json:"balance"`
}

func Test90(t *testing.T) {
	var ret = BalanceData{Balance:10000000000}
	bt, err := json.Marshal(ret)
	if err != nil {
		//
		panic(err)
	}
	var op = ciRes{
		Ret:     0,
		Data:    string(bt),
		Message: "",
	}
	nb, err := json.Marshal(op)
	if err != nil {
		panic(err)
	}
	fmt.Println(nb)
	fmt.Println(string(nb))

	var ty ciRes
	var res BalanceData
	err1 := json.Unmarshal(nb, &ty)
	if err1 != nil {
		panic(err1)
	}

	d := []byte(ty.Data)
	err2 := json.Unmarshal(d, &res)
	if err2 != nil {
		panic(err2)
	}
	fmt.Println(res.Balance)
}