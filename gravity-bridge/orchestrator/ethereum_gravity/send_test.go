package ethereum_gravity

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/gravity_utils/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/umbracle/go-web3"
	"github.com/umbracle/go-web3/jsonrpc"
	"math/big"
	"testing"
	"time"
)

func TestSendToCosmos(t *testing.T) {
	client, _ := jsonrpc.NewClient("http://127.0.0.1:8545")

	Erc20Address := "0x5c702Fbbcfb8EF5cc70c4E4341AA437ef9D55281"
	EthKey := "a8a54b2d8197bc0b19bb8a084031be71835580a01e70a45a13babd16c9bc1563"
	GravityContract := "0x2b7dEe2CF60484325716A1c6A193519c8c3b19F3"

 	timeout := 60 * time.Second
	CosmosDestination := "0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c"

	privKey, _ := crypto.HexToECDSA(EthKey)

	txRes, err := SendToCosmos(Erc20Address, GravityContract, big.NewInt(100000000), common.HexToAddress(CosmosDestination), privKey, timeout, client)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(txRes)
}

type valsetWrapper struct {
	Type string `json:"type"`
	Value types.ValSetMember
}

func TestDeploy(t *testing.T) {
	var a valsetWrapper
	bz, err := base64.StdEncoding.DecodeString("ewogICJ0eXBlIjogImdyYXZpdHkvVmFsc2V0IiwKICAidmFsdWUiOiB7CiAgICAibm9uY2UiOiAiMzYyNiIsCiAgICAibWVtYmVycyI6IFsKICAgICAgewogICAgICAgICJwb3dlciI6ICI0Mjk0OTY3Mjk1IiwKICAgICAgICAiZXRoZXJldW1fYWRkcmVzcyI6ICIweDNGNDNFNzVBYWJhMmMyZkQ2RTIyN0MxMEM2RTdEQzEyNUE5M0RFM2MiCiAgICAgIH0KICAgIF0sCiAgICAiaGVpZ2h0IjogIjM2MjYiCiAgfQp9")
	fmt.Println(bz)
	err = json.Unmarshal(bz, &a)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(a)
}

func TestTxVerify(t *testing.T) {
	client, _ := jsonrpc.NewClient("http://127.0.0.1:8545")
	hash := web3.HexToHash("0x9ef6819edb07e0aa67267fe130da8a9fd2a733673b4ba907530d1117587f4639")
	receipt, err := client.Eth().GetTransactionReceipt(hash)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(receipt)
}