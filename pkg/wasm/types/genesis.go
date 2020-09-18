package types

const (
	communityContractAddress = "0xfffffffffffffffffffffffffffffffffffffff1"
	isOfficial = "true"
	aclContractAddress = "0xfffffffffffffffffffffffffffffffffffffff2"
	votingContractAddress = "0xfffffffffffffffffffffffffffffffffffffff3"
	supportRequiredPct = "600000000000000000"  //60*10^16 :60%
	minAcceptQuorumPct = "500000000000000000"  //50*10^16 :50%
	openTime = "1290600000"  //2 weeks
	InitMethod = "init"
	InvokeMethod = "invoke"
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

func DefaultGenesisState() GenesisState{

	contracts := []Contract{
		{
			Index:   0,
			Code:    communityCode,
			Method:  InitMethod,
			Params:  []string{isOfficial},
			Address: communityContractAddress,
		},
		{
			Index:   1,
			Code:    aclCode,
			Method:  InitMethod,
			Params:  []string{communityContractAddress},
			Address: aclContractAddress,
		},
		{
			Index:   2,
			Code:    votingCode,
			Method:  InitMethod,
			Params:  []string{aclContractAddress, communityContractAddress, supportRequiredPct, minAcceptQuorumPct, openTime, isOfficial},
			Address: votingContractAddress,
		},
		{
			Index:   3,
			Code:    communityCode,
			Method:  InvokeMethod,
			Params:  []string{SetContractMethod, aclContractAddress, votingContractAddress},
			Address: communityContractAddress,
		},

	}

	return GenesisState{
		Contracts:contracts,
		Invoker:invoker,
		Name: "OfficialContract",
		Version: "v0.0.1",
		Author: "Official",
		Email: "ci123chain@corp-ci.com",
		Describe: "OfficialContract",
	}
}