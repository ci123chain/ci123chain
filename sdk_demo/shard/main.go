package main

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/sdk/broadcast"
	shardsdk "github.com/ci123chain/ci123chain/sdk/shard"
)

var (
	online, async, isIBC bool
	from = "0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c"
	gas = "20000"
	nonce = ""
	t ="ADD"
	name = "ciChain-1"
	height = "800"
	priv = "2b452434ac4f7cf9c5d61d62f23834f34e851fb6efdb8d4a8c6e214a8bc93d70"
	proxy = "lb"
	requestURL = "http://ciChain:3030/tx/broadcast"
	//var requestURL = "http://ciChain:3030/tx/broadcast_async"

	onlienReqUrl = "http://ciChain:3030/tx/addShard"
	offlineGas = uint64(20000)
	offlineNonce = uint64(2)
	offlineHeight = int64(900)
)

func main() {
	online = false
	//online = true
	async = false
	//async == true
	isIBC = false
	//
	if online == true {
		fmt.Println("---------------添加分片在线签名交易----------------------")
		fmt.Println("---交易结果：---")
		shardsdk.HttpAddShardTx(from, gas, nonce, t, name, height, priv, proxy, onlienReqUrl)
	}else {
		fmt.Println("---------------添加分片离线签名交易----------------------")
		tx, err := signAddShardTxDemo()
		if err != nil {
			fmt.Println("签名失败，参数错误")
			fmt.Println(err)
			return
		}
		fmt.Print("---签名的交易：---")
		fmt.Print(tx)
		if async == false {
			fmt.Println("---同步广播交易，等待交易结果：---")
			b, _, err := sdk.SendTransaction(tx, async, isIBC, requestURL)
			if err != nil {
				fmt.Println("交易失败")
				fmt.Println(err)
				return
			}
			fmt.Println(string(b))
		}else {
			fmt.Println("---异步广播交易，无返回结果---")
			_, _, _ =sdk.SendTransaction(tx, async, isIBC, requestURL)
		}
	}
}
