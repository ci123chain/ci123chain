package main

import (
	"fmt"
	sdk "github.com/tanhuiya/ci123chain/sdk/broadcast"
	transfersdk "github.com/tanhuiya/ci123chain/sdk/transfer"
)

const (
	isIBC = false
)

func main() {
	//
	var online, async bool
	online = false
	//online = true
	async = false
	//async = true
	//
	if online == true {
		fmt.Println("---------------普通转账在线签名交易----------------------")
		from := "0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c"
		to := "0x505A74675dc9C71eF3CB5DF309256952917E801e"
		amount := "2"
		gas := "20000"
		nonce := "2"
		priv := "2b452434ac4f7cf9c5d61d62f23834f34e851fb6efdb8d4a8c6e214a8bc93d70"
		proxy := "lb"
		fmt.Println("---交易结果：---")
		transfersdk.HttpTransferTx(from, to, gas, nonce, amount, priv, proxy)
	}else {
		fmt.Println("---------------普通转账离线签名交易----------------------")
		tx, err := SignTransferTxDemo()
		if err != nil {
			fmt.Println("签名失败，参数错误")
			fmt.Println(err)
			return
		}
		fmt.Println("---签名的交易：---")
		fmt.Println(tx)
		if async == true {
			fmt.Println("---同步广播交易，等待交易结果：---")
			b,_, err := sdk.SendTransaction(tx, async, isIBC)
			if err != nil {
				fmt.Println("交易失败")
				fmt.Println(err)
				return
			}
			fmt.Println(string(b))
		}else {
			fmt.Println("---异步广播交易，无返回结果---")
			_,_, _ = sdk.SendTransaction(tx, async, isIBC)
		}
	}
}
