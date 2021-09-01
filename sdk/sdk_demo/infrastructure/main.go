package main

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
)

var (
	from = "0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c"
	nonce uint64 = 1
	gas uint64 = 2000
	priv = "2b452434ac4f7cf9c5d61d62f23834f34e851fb6efdb8d4a8c6e214a8bc93d70"
	key = "key"
	content = "value"
)

func main() {
	acc := types.HexToAddress(from)
	tx, err := SignStoreContent(acc, gas, nonce, priv, key, content)
	if err != nil {
		fmt.Println(fmt.Sprintf("签名失败，原因是: %v", err.Error()))
		return
	}
	fmt.Println(fmt.Sprintf("签名后的交易: %v", tx))
}
