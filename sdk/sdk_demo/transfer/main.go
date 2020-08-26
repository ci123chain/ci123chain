package main

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/sdk/broadcast"
	transfersdk "github.com/ci123chain/ci123chain/sdk/transfer"
)


var (
	isIBC, isFabric, online, async bool
	from = "0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c"
	to = "0x505A74675dc9C71eF3CB5DF309256952917E801e"
	amount = "2"
	gas = "20000"
	nonce = ""
	priv = "2b452434ac4f7cf9c5d61d62f23834f34e851fb6efdb8d4a8c6e214a8bc93d70"
	proxy = "lb"
	requestURL = "http://ciChain:3030/tx/broadcast"
	// requestURL = "http://ciChain:3030/tx/broadcast_async"
	offlineGas = uint64(20000)
	offlineNonce = uint64(2)
	offlineAmount = uint64(2)
	onlineReqUrl = "http://ciChain:3030/tx/transfers"
)

func main() {
	online = false
	//online = true
	async = false
	//async = true
	//
	isIBC = false
	if online == true {
		fmt.Println("---------------普通转账在线签名交易----------------------")

		fmt.Println("---交易结果：---")
		transfersdk.HttpTransferTx(from, to, gas, nonce, amount, priv, proxy, onlineReqUrl)
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
			b,_, err := sdk.SendTransaction(tx, async, isIBC, requestURL)
			if err != nil {
				fmt.Println("交易失败")
				fmt.Println(err)
				return
			}
			fmt.Println(string(b))
		}else {
			fmt.Println("---异步广播交易，无返回结果---")
			_,_, _ = sdk.SendTransaction(tx, async, isIBC, requestURL)
		}
	}
}
