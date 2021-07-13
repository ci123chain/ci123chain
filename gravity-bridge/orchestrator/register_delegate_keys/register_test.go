package register_delegate_keys

import (
	"fmt"
	"github.com/ci123chain/ci123chain/gravity-bridge/orchestrator/cosmos_gravity"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"testing"
)

func TestRegisterDelegateKeys(t *testing.T) {
	cosmosRpc := "http://127.0.0.1:1317"

	CosmosAddress := "0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c"
	EthAddress := "0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c"
	EthKey := "2b452434ac4f7cf9c5d61d62f23834f34e851fb6efdb8d4a8c6e214a8bc93d70"
	contact := cosmos_gravity.NewContact(cosmosRpc)

	privKey, _ := crypto.HexToECDSA(EthKey)
	txRes, err := UpdateGravityDelegateAddresses(contact, common.HexToAddress(EthAddress), common.HexToAddress(CosmosAddress), privKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(txRes)
}