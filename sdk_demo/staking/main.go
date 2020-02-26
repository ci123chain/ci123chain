package main

import (
	"fmt"
	sdk "github.com/tanhuiya/ci123chain/sdk/broadcast"
	sSDK "github.com/tanhuiya/ci123chain/sdk/staking"
)

var (
	async bool
	isIBC bool

	create bool
	delegate bool
	redelegate bool
	undelegate bool
	online bool
	requestURL = "http://ciChain:3030/tx/broadcast"
	// requestURL = "http://ciChain:3030/tx/broadcast_async"

	from = "0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c"
	amount uint64 = 20
	gas uint64 = 20000
	nonce uint64 = 1
	pri = "2b452434ac4f7cf9c5d61d62f23834f34e851fb6efdb8d4a8c6e214a8bc93d70"
	proxy = "lb"
	minSelfDelegation int64 = 2
	validatorAddress = "0xB6727FCbC60A03A6689AEE6E5fBC83a7FDc9beBf"
	validatorSrcAddress = ""
	validatorDstAddress = ""
	delegatorAddress = "0xB6727FCbC60A03A6689AEE6E5fBC83a7FDc9beBf"
	rate int64 = 1.000000000000000000
	maxRate int64 = 1.000000000000000000
	maxChangeRate int64 = 1.000000000000000000
	moniker = "first"
	identity = "identity"
	website = "website"
	securityContact = "security"
	details = "details"
	publicKey = "7b0a2274797065223a202274656e6465726d696e742f5075624b6579536563703235366b31222c0a2276616c7565223a20224138485a776b6442307544497150594536783530394d304a654c56585a2b70613172343279706e4630316e36220a7d"
	//"validatorAddress": "0xB6727FCbC60A03A6689AEE6E5fBC83a7FDc9beBf"
	//"address": "B6727FCBC60A03A6689AEE6E5FBC83A7FDC9BEBF"
	//{
	//"type": "tendermint/PubKeySecp256k1",
	//"value": "A8HZwkdB0uDIqPYE6x509M0JeLVXZ+pa1r42ypnF01n6"
	//}

	onlineGas = "20000"
	onlineNonce = "2"
	onlineAmount = "2"
	delegateURL = "http://ciChain:3030/staking/delegate"
	redelegateURL = "http://ciChain:3030/staking/redelegate"
	undelegateURL = "http://ciChain:3030/staking/undelegate"
)


func main() {
	async = false
	isIBC = false
	delegate = false
	redelegate = false
	undelegate = false
	create = true
	online = false

	if create {
		fmt.Println("---------------添加验证者离线签名交易----------------------")
		tx, err := SignCreateValidatorTx(from, amount, gas, nonce, pri, minSelfDelegation, validatorAddress, delegatorAddress, rate, maxRate,
			maxChangeRate, moniker, identity, website, securityContact, details, publicKey)
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
			fmt.Print("---交易结果：---")
			fmt.Println(string(b))
		}else {
			fmt.Println("---异步广播交易，无返回结果---")
			_, _, _ =sdk.SendTransaction(tx, async, isIBC, requestURL)
		}
		fmt.Println("---------------添加验证者离线签名交易完成----------------------")
	}

	if delegate {
		if online {
			fmt.Println("---------------抵押在线签名交易----------------------")
			sSDK.HttpDelegateTx(from, onlineGas, onlineNonce, onlineAmount, pri, validatorAddress, delegatorAddress,proxy, delegateURL)
			fmt.Println("---------------抵押在线签名交易完成----------------------")
		}else {
			fmt.Println("---------------抵押离线签名交易----------------------")
			delegateTx, err := SignDelegateTx(from, amount, gas, nonce, pri, validatorAddress, delegatorAddress)
			if err != nil {
				fmt.Println("签名失败，参数错误")
				fmt.Println(err)
				return
			}
			fmt.Print("---签名的交易：---")
			fmt.Print(delegateTx)
			if async == false {
				fmt.Println("---同步广播交易，等待交易结果：---")
				b, _, err := sdk.SendTransaction(delegateTx, async, isIBC, requestURL)
				if err != nil {
					fmt.Println("交易失败")
					fmt.Println(err)
					return
				}
				fmt.Print("---交易结果：---")
				fmt.Println(string(b))
			}else {
				fmt.Println("---异步广播交易，无返回结果---")
				_, _, _ =sdk.SendTransaction(delegateTx, async, isIBC, requestURL)
			}
			fmt.Println("---------------抵押离线签名交易完成----------------------")
		}
	}

	if redelegate {
		if online {
			fmt.Println("---------------重新抵押在线签名交易----------------------")
			sSDK.HttpRedelegateTx(from, onlineGas, onlineNonce, onlineAmount, pri, validatorSrcAddress, validatorDstAddress, delegatorAddress,proxy, redelegateURL)
			fmt.Println("---------------重新抵押在线签名交易完成----------------------")
		}else {
			fmt.Println("---------------重新抵押离线签名交易----------------------")
			redelegateTx, err := SignRelegateTx(from, amount, gas, nonce, pri, validatorSrcAddress, validatorDstAddress, delegatorAddress)
			if err != nil {
				fmt.Println("签名失败，参数错误")
				fmt.Println(err)
				return
			}
			fmt.Print("---签名的交易：---")
			fmt.Print(redelegateTx)
			if async == false {
				fmt.Println("---同步广播交易，等待交易结果：---")
				b, _, err := sdk.SendTransaction(redelegateTx, async, isIBC, requestURL)
				if err != nil {
					fmt.Println("交易失败")
					fmt.Println(err)
					return
				}
				fmt.Print("---交易结果：---")
				fmt.Println(string(b))
			}else {
				fmt.Println("---异步广播交易，无返回结果---")
				_, _, _ =sdk.SendTransaction(redelegateTx, async, isIBC, requestURL)
			}
			fmt.Println("---------------重新抵押离线签名交易完成----------------------")
		}
	}

	if undelegate {
		if online {
			fmt.Println("---------------解除抵押在线签名交易----------------------")
			sSDK.HttpUndelegateTx(from, onlineGas, onlineNonce, onlineAmount, pri, validatorAddress, delegatorAddress,proxy, undelegateURL)
			fmt.Println("---------------解除抵押在线签名交易完成----------------------")
		}else {
			fmt.Println("---------------解除抵押离线签名交易----------------------")
			undelegateTx, err := SignUndelegate(from, amount, gas, nonce, pri, validatorAddress, delegatorAddress)
			if err != nil {
				fmt.Println("签名失败，参数错误")
				fmt.Println(err)
				return
			}
			fmt.Print("---签名的交易：---")
			fmt.Print(undelegateTx)
			if async == false {
				fmt.Println("---同步广播交易，等待交易结果：---")
				b, _, err := sdk.SendTransaction(undelegateTx, async, isIBC, requestURL)
				if err != nil {
					fmt.Println("交易失败")
					fmt.Println(err)
					return
				}
				fmt.Print("---交易结果：---")
				fmt.Println(string(b))
			}else {
				fmt.Println("---异步广播交易，无返回结果---")
				_, _, _ =sdk.SendTransaction(undelegateTx, async, isIBC, requestURL)
			}
			fmt.Println("---------------解除抵押离线签名交易完成----------------------")
		}
	}

}
