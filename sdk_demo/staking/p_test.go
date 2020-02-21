package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	sdk "github.com/tanhuiya/ci123chain/pkg/abci/types"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"io"
	"testing"
)

var cdc = amino.NewCodec()

func init() {
	cryptoAmino.RegisterAmino(cdc)
}
func Test1(t *testing.T) {

	privSecret := "nodenodenodenodenodenodenodenode"
	sbz := []byte(privSecret)
	var b32 [32]byte
	io.ReadFull(bytes.NewReader(sbz), b32[:])
	privateKey := secp256k1.PrivKeySecp256k1(b32)

	pk := privateKey.PubKey()
	fmt.Println(pk.Address())
	b, err := cdc.MarshalJSONIndent(pk, "", "")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	byte := pk.Bytes()
	fmt.Println(byte)
	addr := sdk.ToAccAddress(byte).String()
	fmt.Println(addr)
	hb := sdk.HexToAddress(addr)
	fmt.Println(hb)

}

type rt struct {
	Addr   sdk.AccAddress  `json:"addr"`
}

func Test6(t *testing.T) {
	var p = "AtuUBaKV6Vmchuq/ZqGH7sD2y5bUmJnwjOW9A3b3VnER"
	pubStr := fmt.Sprintf(`{"type":"%s","value":"%s"}`, secp256k1.PubKeyAminoName, p)
	var pk secp256k1.PubKeySecp256k1
	cdc.UnmarshalJSON([]byte(pubStr), &pk)
	pp := crypto.PubKey(pk)
	var ty = rt{Addr:sdk.ToAccAddress(pp.Address())}
	bz, err := json.Marshal(ty)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bz))
	//fmt.Println(pk.Address())
}

func Test2(t *testing.T) {
	var ad = "0BBA3A8E46DDC7328F078B9244306EC097A68E17"
	var ty = "0x0bba3a8E46dDC7328F078B9244306eC097A68E17"
	addr := sdk.HexToAddress(ad)
	fmt.Println(addr.String())
	aty := sdk.HexToAddress(ty)
	fmt.Println(aty.String())
}

func Test3(t *testing.T) {
	delegateTx, err := SignDelegateTx(from, amount, gas, nonce, pri, validatorAddress, delegatorAddress)
	if err != nil {
		fmt.Println("签名失败，参数错误")
		fmt.Println(err)
		return
	}
	fmt.Println(delegateTx)
}

func Test4(t *testing.T) {
	delegateTx, err := SignDelegateTx(from, amount, gas, nonce, pri, validatorAddress, delegatorAddress)
	if err != nil {
		fmt.Println("签名失败，参数错误")
		fmt.Println(err)
		return
	}
	fmt.Println(delegateTx)
}