package types

import (
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	//invoker = "0x9BA7dc2269895DF1004Ec75D8326644295508069"
	communityContractAddress = "0xfffffffffffffffffffffffffffffffffffffff1"
	isOfficial = "true"
	aclContractAddress = "0xfffffffffffffffffffffffffffffffffffffff2"
	votingContractAddress = "0xfffffffffffffffffffffffffffffffffffffff3"
	supportRequiredPct = "6000000000"
	minAcceptQuorumPct = "5000000000"
	openTime = "1290600000"
	InitMethod = "init"
	//InvokeMethod = "invoke"
	SetContractMethod = "initial_contract"
	invoker = communityContractAddress
)

type Contract struct {
	Index    int        `json:"index"`
	Code     string     `json:"code"`
	Method   string     `json:"method"`
	Params   []string   `json:"params"`
	Address  string     `json:"address"`
}


type GenesisState struct {
	Contracts []Contract `json:"contracts"`
	Invoker  string      `json:"invoker"` //invoker.
	Name     string      `json:"name"`
	Version  string      `json:"version"`
	Author   string      `json:"author"`
	Email    string      `json:"email"`
	Describe string      `json:"describe"`
}

func DefaultGenesisState(_ []tmtypes.GenesisValidator) GenesisState{

	contracts := []Contract{
		{
			Index: 0,
			Code: community_code,
			Method: InitMethod,
			Params: []string{isOfficial},
			Address: communityContractAddress,
		},
		{
			Index: 1,
			Code: acl_code,
			Method: InitMethod,
			Params: []string{communityContractAddress},
			Address: aclContractAddress,
		},
		{
			Index: 2,
			Code: voting_code,
			Method: InitMethod,
			Params: []string{aclContractAddress, communityContractAddress, supportRequiredPct, minAcceptQuorumPct, openTime, isOfficial},
			Address: votingContractAddress,
		},
		{
			Index: 3,
			Code: community_code,
			Method: SetContractMethod,
			Params: []string{aclContractAddress, votingContractAddress},
			Address: communityContractAddress,
		},

	}

	return GenesisState{
		Contracts:contracts,
		Invoker:invoker,
	}
}