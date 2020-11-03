package init

import (
	"encoding/json"
	"fmt"
	sdk "github.com/ci123chain/ci123chain/pkg/abci/types"
	"github.com/ci123chain/ci123chain/pkg/staking/types"
	"github.com/ci123chain/ci123chain/sdk/account"
	"github.com/ci123chain/ci123chain/sdk/validator"
	"strings"
	"testing"
	"time"
)

func TestNewInitChainFiles(t *testing.T) {
	privKey1, vpubKey1, address1, _ := validator.NewValidatorKey() //node address and privKey/pubKey
	privKey2, vpubKey2, address2, _ := validator.NewValidatorKey() //node address and privKey/pubKey
	acc1 := account.NewAccount() //account address and privKey
	acc2 := account.NewAccount()

	var cInfo = ChainInfo{
		ChainID:     "ci0",
		GenesisTime: time.Now(),
	}
	var vInfo = []ValidatorInfo{
		{
			PubKey:  vpubKey1,
			Name:    "validator1",
		},
		{
			PubKey:  vpubKey2,
			Name:    "validator2",
		},
	}
	var sInfo = []StakingInfo{
		{
			Address:           sdk.HexToAddress(acc1.Address),
			PubKey:            vpubKey1,
			Tokens:            "10000000",
			CommissionInfo:    CommissionInfo{
				Rate:          1,
				MaxRate:       40,
				MaxChangeRate: 5,
			},
			UpdateTime:        time.Now(),
			MinSelfDelegation: "1000000",
			Description: types.Description{
				Moniker:         "moniker1",
				Identity:        "",
				Website:         "",
				SecurityContact: "",
				Details:         "",
			},
		},
		{
			Address:           sdk.HexToAddress(acc2.Address),
			PubKey:            vpubKey2,
			Tokens:            "10000000",
			CommissionInfo:    CommissionInfo{
				Rate:          1,
				MaxRate:       40,
				MaxChangeRate: 5,
			},
			UpdateTime:        time.Now(),
			MinSelfDelegation: "1000000",
			Description: types.Description{
				Moniker:         "moniker1",
				Identity:        "",
				Website:         "",
				SecurityContact: "",
				Details:         "",
			},
		},
	}
	var supInfo = SupplyInfo{
		Amount: "200000000000",
	}
	var accInfo = []AccountInfo{
		{
			Address: sdk.HexToAddress(acc1.Address),
			Amount: "1000000000000000",
		},
		{
			Address: sdk.HexToAddress(acc2.Address),
			Amount: "1000000000000000",
		},
	}
	persistentPeers := strings.ToLower(address1) + "@127.0.0.1:26656" + "," + strings.ToLower(address2) + "@127.0.0.1:36656"

	//生成genesisFile
	genesisFiles, err := NewGenesisFiles(cInfo, vInfo, sInfo, supInfo, accInfo)
	if err != nil{
		fmt.Println(err)
	}

	//生成的nodeKey是privKey1的
	initFiles1, err := NewInitFiles(privKey1, persistentPeers)
	if err != nil{
		fmt.Println(err)
	}
	//生成的nodeKey是privKey2的
	initFiles2, err := NewInitFiles(privKey2, persistentPeers)
	if err != nil{
		fmt.Println(err)
	}
	b, _ := json.Marshal(genesisFiles)
	fmt.Println(string(b))
	c, _ := json.Marshal(initFiles1)
	fmt.Println(string(c))
	d, _ := json.Marshal(initFiles2)
	fmt.Println(string(d))
}