package main

import (
	"fmt"
	sdk "github.com/tanhuiya/ci123chain/sdk/broadcast"
)

const (
	async = false
	isIBC = false

	create = true
	delegate = false
	redelegate = false
	undelegate = false
)

var (
	from = "0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c"
	amount uint64 = 2
	gas uint64 = 20000
	nonce uint64 = 1
	pri = "2b452434ac4f7cf9c5d61d62f23834f34e851fb6efdb8d4a8c6e214a8bc93d70"
	minSelfDelegation int64 = 2
	validatorAddress = "0xaFD809610Bf8D4f26caCEAbEfF6ad1144d0e6A1D"
	validatorSrcAddress = ""
	validatorDstAddress = ""
	delegatorAddress = "0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c"
	rate int64 = 2
	maxRate int64 = 3
	maxChangeRate int64 = 4
	moniker = ""
	identity = ""
	website = ""
	securityContact = ""
	details = ""
	pubKeyTp = ""
	pubKeyVal = ""
)


func main() {

	if create {
		fmt.Println("---------------添加验证者离线签名交易----------------------")
		tx, err := SignCreateValidatorTx(from, amount, gas, nonce, pri, minSelfDelegation, validatorAddress, delegatorAddress, rate, maxRate,
			maxChangeRate, moniker, identity, website, securityContact, details, pubKeyTp, pubKeyVal)
		if err != nil {
			fmt.Println("签名失败，参数错误")
			fmt.Println(err)
			return
		}
		fmt.Print("---签名的交易：---")
		fmt.Print(tx)
		if async == false {
			fmt.Println("---同步广播交易，等待交易结果：---")
			b, _, err := sdk.SendTransaction(tx, async, isIBC)
			if err != nil {
				fmt.Println("交易失败")
				fmt.Println(err)
				return
			}
			fmt.Print("---交易结果：---")
			fmt.Println(string(b))
		}else {
			fmt.Println("---异步广播交易，无返回结果---")
			_, _, _ =sdk.SendTransaction(tx, async, isIBC)
		}
		fmt.Println("---------------添加验证者离线签名交易完成----------------------")
	}

	if delegate {
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
			b, _, err := sdk.SendTransaction(delegateTx, async, isIBC)
			if err != nil {
				fmt.Println("交易失败")
				fmt.Println(err)
				return
			}
			fmt.Print("---交易结果：---")
			fmt.Println(string(b))
		}else {
			fmt.Println("---异步广播交易，无返回结果---")
			_, _, _ =sdk.SendTransaction(delegateTx, async, isIBC)
		}
		fmt.Println("---------------抵押离线签名交易完成----------------------")
	}

	if redelegate {
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
			b, _, err := sdk.SendTransaction(redelegateTx, async, isIBC)
			if err != nil {
				fmt.Println("交易失败")
				fmt.Println(err)
				return
			}
			fmt.Print("---交易结果：---")
			fmt.Println(string(b))
		}else {
			fmt.Println("---异步广播交易，无返回结果---")
			_, _, _ =sdk.SendTransaction(redelegateTx, async, isIBC)
		}
		fmt.Println("---------------重新抵押离线签名交易完成----------------------")
	}

	if undelegate {
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
			b, _, err := sdk.SendTransaction(undelegateTx, async, isIBC)
			if err != nil {
				fmt.Println("交易失败")
				fmt.Println(err)
				return
			}
			fmt.Print("---交易结果：---")
			fmt.Println(string(b))
		}else {
			fmt.Println("---异步广播交易，无返回结果---")
			_, _, _ =sdk.SendTransaction(undelegateTx, async, isIBC)
		}
		fmt.Println("---------------解除抵押离线签名交易完成----------------------")
	}

}
