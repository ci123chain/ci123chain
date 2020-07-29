package main

import (
	"fmt"
	common "github.com/ci123chain/ci123chain/sdk/broadcast"
	"os"
)

var (
	from = "0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c"
	validatorAddress = "0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c"
	delegatorAddress = "0x2561735fA825F004AFe2cb8Dfe87BFeF539B8Fc5"
	withdrawAddress = "0x2561735fA825F004AFe2cb8Dfe87BFeF539B8Fc5"
	gas uint64 = 20000
	nonce uint64 = 0
	amount int64 = 100
	priv = "2b452434ac4f7cf9c5d61d62f23834f34e851fb6efdb8d4a8c6e214a8bc93d70"
	requestURL = "http://ciChain:3030/tx/broadcast"
)

func main() {
	var sel int
	var async int
	fmt.Println("please input select value, you input one of [0, 1, 2, 3]")
	_, _ = fmt.Scanln(&sel)
	fmt.Println("please select sync or async")
	_, _ = fmt.Scanln(&async)
	if async != 0 && async != 1 {
		fmt.Println(fmt.Sprintf("err async value: %d", async))
		os.Exit(1)
	}
	var tx string
	var err error
	switch sel {
	case 0:
		//从个人账户直接转账到CommunityPool.
		tx, err = SignCommunityPoolTx(from, amount, gas, nonce, priv)
	case 1:
		//提取validator佣金
		tx, err = SignWithdrawCommissionTx(from, validatorAddress, gas, nonce, priv)
	case 2:
		//提取delegator奖金
		tx, err = SignWithdrawRewardsTx(from, validatorAddress, delegatorAddress, gas, nonce, priv)
	case 3:
		//更改提取奖金或佣金的账户地址（提取的奖金会转入这个更改后的账户地址）
		tx, err = SignSetWithdrawAddressTx(from, withdrawAddress, gas, nonce, priv)
	default:
		fmt.Println(fmt.Sprintf("err select value: %d",sel))
		os.Exit(1)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Print("---签名的交易：---")
	fmt.Println(tx)
	if async == 1 {
		//
		fmt.Println("---异步广播交易，无返回结果---")
		_ , _, _ = common.SendTransaction(tx, true, false, requestURL)
	}else if async == 0 {
		//
		fmt.Println("---同步广播交易，等待交易结果：---")
		b, _, err := common.SendTransaction(tx, false, false, requestURL)
		if err != nil {
			fmt.Println("交易失败")
			fmt.Println(err)
			return
		}
		fmt.Print("---交易结果：---")
		fmt.Println(string(b))
	}
}
