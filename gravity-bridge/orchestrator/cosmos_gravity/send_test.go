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
	txRes, err := SendToEth(privKey, &to, types.NewChainCoin(types.NewInt(100000)), types.NewChainCoin(types.NewInt(100000)), contact)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(txRes)
}