package cosmos_gravity

import (
	"fmt"
	"github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func TestSendToEth(t *testing.T) {
	cosmosRpc := "http://127.0.0.1:1317"

	EthAddress := "0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c"
	EthKey := "2b452434ac4f7cf9c5d61d62f23834f34e851fb6efdb8d4a8c6e214a8bc93d70"
	contact := NewContact(cosmosRpc)

	privKey, _ := crypto.HexToECDSA(EthKey)
	to := common.HexToAddress(EthAddress)
	txRes, err := SendToEth(privKey, &to, types.NewCoin("gravity0x5c702Fbbcfb8EF5cc70c4E4341AA437ef9D55281", types.NewInt(1000000)), types.NewCoin("gravity0x5c702Fbbcfb8EF5cc70c4E4341AA437ef9D55281", types.NewInt(10000)), contact)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(txRes)
}

func TestRequestBatch(t *testing.T) {
	cosmosRpc := "http://127.0.0.1:1317"

	EthKey := "a8a54b2d8197bc0b19bb8a084031be71835580a01e70a45a13babd16c9bc1563"
	contact := NewContact(cosmosRpc)

	privKey, _ := crypto.HexToECDSA(EthKey)
	txRes, err := SendRequestBatch(privKey, "gravity0x5c702Fbbcfb8EF5cc70c4E4341AA437ef9D55281", contact)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(txRes)
}