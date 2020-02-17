package cmd

import (
	"bytes"
	"fmt"
	"github.com/tanhuiya/ci123chain/pkg/validator"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"io"
	"testing"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/encoding/amino"
)


var cdc = amino.NewCodec()

func init() {
	cryptoAmino.RegisterAmino(cdc)
}

// {"type":"tendermint/PrivKeySecp256k1","value":"oQLmM5pM5wL78a6LJntQY8tPGQPpp050udIA5YZMkCc="}
// 通过 json 文件里编码后的私钥解码
func TestGenValidator(t *testing.T) {
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

// 通过secret 构造私钥
func TestGenValidator2(t *testing.T) {
	// 私钥编码后的数据
	privSecret := "nodenodenodenodenodenodenodenode"
	sbz := []byte(privSecret)
	var b32 [32]byte
	io.ReadFull(bytes.NewReader(sbz), b32[:])
	privateKey := secp256k1.PrivKeySecp256k1(b32)

	pv := validator.GenFilePV("", "", privateKey)
	jsonBytes, _ := cdc.MarshalJSONIndent(pv.Key, "", "  ")
	fmt.Printf("%s", jsonBytes)
}

