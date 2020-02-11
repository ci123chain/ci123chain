package main

import (
	"encoding/json"
	"fmt"
	sdk "github.com/tanhuiya/ci123chain/sdk/account"
)

type Response struct {
	Ret uint32 `json:"ret"`
	Data interface{} `json:"data"`
	Message string   `json:"message"`
}

func main() {
	var res Response
	var online bool
	online = false
	//online = true

	if online == true {
		fmt.Println("---------------在线生成新的账户----------------------")
		b, err := sdk.NewAccountOnLine()
		if err != nil {
			fmt.Println("---生成失败---")
			fmt.Println(err)
			return
		}
		err = json.Unmarshal(b, &res)
		if err != nil {
			fmt.Println("---解析失败---")
			fmt.Println(err)
			return
		}
		fmt.Println("---生成的账户信息：---")
		fmt.Println(res.Data)
	}else {
		fmt.Println("---------------离线生成新的账户----------------------")
		address, privateKey, err := sdk.NewAccountOffLine()
		if err != nil {
			fmt.Println("---生成失败---")
			fmt.Println(err)
			return
		}
		fmt.Println("---生成的账户信息：---")
		fmt.Println("address:", address)
		fmt.Println("privateKey:", privateKey)
	}
}
