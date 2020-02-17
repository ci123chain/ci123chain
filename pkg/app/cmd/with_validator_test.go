package cmd

import (
	"fmt"
	"github.com/tanhuiya/ci123chain/pkg/validator"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
)


var cdc = amino.NewCodec()

func init() {
	cryptoAmino.RegisterAmino(cdc)
}

// {"type":"tendermint/PrivKeySecp256k1","value":"oQLmM5pM5wL78a6LJntQY8tPGQPpp050udIA5YZMkCc="}
func TestValidator(t *testing.T) {
	// 私钥编码后的数据
	privBz := "oQLmM5pM5wL78a6LJntQY8tPGQPpp050udIA5YZMkCc="

	privStr := fmt.Sprintf(`{"type":"%s","value":"%s"}`, secp256k1.PrivKeyAminoName, privBz)
	var privateKey secp256k1.PrivKeySecp256k1
	// 生成的对象privateKey 直接使用
	err := cdc.UnmarshalJSON([]byte(privStr), &privateKey)
	if err != nil {
		panic(err)
	}

	pv := validator.GenFilePV("", "", privateKey)
	jsonBytes, _ := cdc.MarshalJSONIndent(pv.Key, "", "  ")

	fmt.Printf("%s", jsonBytes)
}
