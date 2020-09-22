package init

import (
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/sdk/validator"
	"github.com/tendermint/tendermint/crypto"
	"testing"
	"time"
)

var pubKey crypto.PubKey
var privKey string
var address string
func TestNewInitChainFiles(t *testing.T) {
	privKey, pubKey, address, _ = validator.NewValidatorKey()

	var cInfo = ChainInfo{
		ChainID:     "ci0",
		GenesisTime: time.Now(),
	}
	var vInfo = ValidatorInfo{
		PubKey:  pubKey,
		Name:    "validator1",
	}
	var sInfo = StakingInfo{
		Address:           sdk.HexToAddress("0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c"),
		PubKey:            pubKey,
		Tokens:            10000000,
		CommissionInfo:    CommissionInfo{
			Rate:          1,
			MaxRate:       40,
			MaxChangeRate: 5,
		},
		UpdateTime:        time.Now(),
	}
	var supInfo = SupplyInfo{
		Amount: 200000000000,
	}
	var accInfo = AccountInfo{
		Address: sdk.HexToAddress("0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c"),
		Amount: 1000000000000000,
	}
	persistentPeers := address + "@127.0.0.1:80"

	initFiles, err := NewInitChainFiles(cInfo, vInfo, sInfo, supInfo, accInfo, privKey, persistentPeers)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println(initFiles)
}