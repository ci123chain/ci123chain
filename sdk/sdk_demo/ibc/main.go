package main

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/sdk/broadcast"
)
var (
	async bool
	isIBC bool
	from = "0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c"
	to = "0x505A74675dc9C71eF3CB5DF309256952917E801e"
	amount = uint64(1)
	gas = uint64(20000)
	nonce = uint64(2)
	priv = "2b452434ac4f7cf9c5d61d62f23834f34e851fb6efdb8d4a8c6e214a8bc93d70"
	requestURL = "http://ciChain:3030/tx/broadcast"
	// requestURL = "http://ciChain:3030/tx/broadcast_async"
)

func main() {
	async = false
	isIBC = true
	//
	fmt.Println("---------------跨链转账离线签名交易----------------------")
	fmt.Println("---发送跨链消息---")

	txByte, err := SignIBC(from, to, amount, gas,nonce, priv)
	if err != nil {
		fmt.Println("---签名失败，参数错误---")
		fmt.Println(err)
		return
	}
	tx := hex.EncodeToString(txByte)
	fmt.Println(tx)
	_, registRet, _ := sdk.SendTransaction(tx, async, isIBC, requestURL)
	if registRet.Data == "" {
		fmt.Println("---发送跨链消息失败---")
		return
	}
	fmt.Println("发送跨链消息完成：UniqueID = " + registRet.Data)

	fmt.Println("---申请处理该跨链消息---")
	uid := []byte(registRet.Data)
	observerID := []byte("1234567812345679")
	signdata, err := SignIBCApplyTx(from, uid, observerID, gas, nonce, priv)
	if err != nil {
		fmt.Println("---签名失败，参数错误---")
		fmt.Println(err)
		return
	}
	applyTx := hex.EncodeToString(signdata)
	fmt.Println(applyTx)
	_, applyRet, _ := sdk.SendTransaction(applyTx, async, isIBC, requestURL)
	if len(applyRet.RawLog) < 1 {
		fmt.Println("---申请处理该跨链消息成功---")
	}else {
		fmt.Println("---申请处理该跨链消息失败---")
		return
	}

	fmt.Println("---向对方转账---")
	pkg := applyRet.Data
	signdata, err = SignIBCBankSendTx(from, []byte(pkg), gas, nonce, priv)
	if err != nil {
		fmt.Println("---签名失败，参数错误---")
		fmt.Println(err)
		return
	}
	receiptTx := hex.EncodeToString(signdata)
	fmt.Println(receiptTx)
	_, receiptRet, _ := sdk.SendTransaction(receiptTx, async, isIBC, requestURL)
	if len(receiptRet.RawLog) < 1 {
		fmt.Println("---向对方转账成功---")
	}else {
		fmt.Println("---向对方转账失败---")
		return
	}

	fmt.Println("---发送回执---")
	receivepkg := receiptRet.Data
	signdata, err = SignIBCReceiptTx(from, []byte(receivepkg), gas, nonce, priv)
	if err != nil {
		fmt.Println("---签名失败，参数错误---")
		fmt.Println(err)
		return
	}
	receiveTx := hex.EncodeToString(signdata)
	fmt.Println(receiveTx)
	_, ret, _ := sdk.SendTransaction(receiveTx, async, isIBC, requestURL)
	if len(ret.RawLog) < 1 {
		fmt.Println("---发送回执成功---")
	}else {
		fmt.Println("---发送回执失败---")
		return
	}

}